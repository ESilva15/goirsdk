package mmaputils_test

import (
	"testing"
	"time"

	eventutils "github.com/ESilva15/goirsdk/eventutils"
)

const testEventName = "test_iRSDKDataValidEvent"

// Tests initialization and resource cleanup
func TestInitAndClose(t *testing.T) {
	sdkUtils, err := eventutils.Init()
	if err != nil {
		t.Fatalf("Init() failed: %v", err)
	}
	if sdkUtils == nil {
		t.Fatalf("Init() returned nil struct")
	}

	sdkUtils.Close()
}

// Tests timeout behavior when telemetry stalls or stops
func TestCheckValidDataEvent_Timeout(t *testing.T) {
	defer eventutils.CleanupEvent(testEventName)

	sdkUtils, err := eventutils.Init()
	if err != nil {
		t.Fatalf("Init() failed: %v", err)
	}
	defer sdkUtils.Close()

	if err := sdkUtils.OpenEvent(testEventName); err != nil {
		t.Fatalf("OpenEvent() failed: %v", err)
	}

	timeout := 40 * time.Millisecond
	start := time.Now()

	// Should return false because no producer signaled the event
	got := sdkUtils.CheckValidDataEvent(timeout)
	elapsed := time.Since(start)

	if got != false {
		t.Errorf("expected CheckValidDataEvent to return false on timeout, got true")
	}

	if elapsed < timeout {
		t.Errorf("expected timeout to wait at least %v, returned early after %v", timeout, elapsed)
	}
}

// Tests successful unblocking when a frame signal arrives
func TestCheckValidDataEvent_Signaled(t *testing.T) {
	defer eventutils.CleanupEvent(testEventName)

	sdkUtils, err := eventutils.Init()
	if err != nil {
		t.Fatalf("Init() failed: %v", err)
	}
	defer sdkUtils.Close()

	if err := sdkUtils.OpenEvent(testEventName); err != nil {
		t.Fatalf("OpenEvent() failed: %v", err)
	}

	// Simulate mock producer sending a frame signal after 10ms
	go func() {
		time.Sleep(10 * time.Millisecond)
		eventutils.SignalEvent(testEventName)
	}()

	start := time.Now()
	got := sdkUtils.CheckValidDataEvent(200 * time.Millisecond)
	elapsed := time.Since(start)

	if got != true {
		t.Errorf("expected CheckValidDataEvent to return true on signal, got false")
	}

	if elapsed >= 150*time.Millisecond {
		t.Errorf("CheckValidDataEvent took too long (%v), failed to wake up immediately", elapsed)
	}
}

// Tests multi-frame streaming loop (simulating steady 60Hz feed)
func TestCheckValidDataEvent_MultipleTicks(t *testing.T) {
	defer eventutils.CleanupEvent(testEventName)

	sdkUtils, err := eventutils.Init()
	if err != nil {
		t.Fatalf("Init() failed: %v", err)
	}
	defer sdkUtils.Close()

	if err := sdkUtils.OpenEvent(testEventName); err != nil {
		t.Fatalf("OpenEvent() failed: %v", err)
	}

	frameCount := 5
	go func() {
		for i := 0; i < frameCount; i++ {
			time.Sleep(10 * time.Millisecond)
			eventutils.SignalEvent(testEventName)
		}
	}()

	for i := 0; i < frameCount; i++ {
		ok := sdkUtils.CheckValidDataEvent(100 * time.Millisecond)
		if !ok {
			t.Fatalf("failed to receive pulse for tick %d", i)
		}
	}
}

func Test60FPS60SecondsSemaphores(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping 60-second endurance test in short mode")
	}

	const eventName = "test_60fps_event"
	const targetFPS = 60
	const durationSeconds = 60
	const totalFrames = targetFPS * durationSeconds // 3,600 frames
	const frameInterval = time.Second / targetFPS   // ~16.666ms

	// Timeout per frame set to 100ms to absorb normal OS thread scheduling jitter
	const frameWaitTimeout = 250 * time.Millisecond

	// Cleanup prior event state
	eventutils.CleanupEvent(eventName)
	defer eventutils.CleanupEvent(eventName)

	// Initialize Reader
	uReader, err := eventutils.Init()
	if err != nil {
		t.Fatalf("Failed to initialize reader: %v", err)
	}
	defer uReader.Close()

	if err := uReader.OpenEvent(eventName); err != nil {
		t.Fatalf("Failed to open event on reader: %v", err)
	}

	// Initialize Writer
	uWriter, err := eventutils.Init()
	if err != nil {
		t.Fatalf("Failed to initialize writer: %v", err)
	}
	defer uWriter.Close()

	if err := uWriter.OpenEvent(eventName); err != nil {
		t.Fatalf("Failed to open event on writer: %v", err)
	}

	stopWriter := make(chan struct{})
	writerDone := make(chan struct{})

	// Producer Goroutine: Emits pulse at 60 FPS
	go func() {
		defer close(writerDone)
		ticker := time.NewTicker(frameInterval)
		defer ticker.Stop()

		for {
			select {
			case <-stopWriter:
				return
			case <-ticker.C:
				if err := uWriter.SignalEvent(); err != nil {
					t.Errorf("SignalEvent failed on writer: %v", err)
					return
				}
			}
		}
	}()

	startTime := time.Now()
	receivedFrames := 0

	// Consumer Loop: Consumes 3,600 frames continuously
	for i := 1; i <= totalFrames; i++ {
		ok := uReader.CheckValidDataEvent(frameWaitTimeout)
		if !ok {
			close(stopWriter)
			<-writerDone
			t.Fatalf("FAILED at frame %d/%d (elapsed: %v). Semaphore timed out after %v.",
				i, totalFrames, time.Since(startTime), frameWaitTimeout)
		}
		receivedFrames++
	}

	elapsed := time.Since(startTime)
	close(stopWriter)
	<-writerDone

	actualFPS := float64(receivedFrames) / elapsed.Seconds()
	t.Logf("Passed: Processed %d/%d frames continuously in %v (Average FPS: %.2f)",
		receivedFrames, totalFrames, elapsed, actualFPS)
}
