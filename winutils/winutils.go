package winutils

import (
	"fmt"
	"sync"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	WAIT_OBJECT_0 = 0
	WAIT_TIMEOUT  = 258
)

var (
	once      sync.Once
	user32DLL *windows.LazyDLL = nil
)

// openEvent opens a windows.Handle for a given event
func OpenEvent(eventName string) (windows.Handle, error) {
	name, err := windows.UTF16PtrFromString(eventName)
	if err != nil {
		return 0, err
	}

	handle, err := windows.OpenEvent(windows.SYNCHRONIZE, false, name)
	if err != nil {
		return 0, fmt.Errorf("error opening event %s: %w", eventName, err)
	}

	return handle, nil
}

// closeEvent closes a given windows.Handle
func CloseEvent(h windows.Handle) {
	windows.CloseHandle(h)
}

// openEvent waits for a good response for some given time
func CheckValidDataEvent(handle windows.Handle, timeout time.Duration) bool {
	t0 := time.Now().UnixNano()
	timeoutInt := uint32(timeout / time.Millisecond)

	result, err := windows.WaitForSingleObject(handle, timeoutInt)
	if err != nil {
		remaining := timeoutInt - uint32((time.Now().UnixNano()-t0)/1000000)
		if remaining > 0 {
			time.Sleep(time.Duration(remaining) * time.Millisecond)
		}
		return false
	}

	// Check the result of the wait
	if result == WAIT_OBJECT_0 {
		return true
	} else if result == WAIT_TIMEOUT {
		return false
	}

	return false
}

// loadUser32DLL loads the user32.dll which is used to create some processes
func loadUser32DLL() {
	user32DLL = windows.NewLazyDLL("user32.dll")
}

// OpenBroadcastChannel opens up a broadcast channel to send commands to iracing
func OpenBroadcastChannel(name string) (uintptr, error) {
	if user32DLL == nil {
		once.Do(loadUser32DLL)
	}
	registerWindowsMessageW := user32DLL.NewProc("RegisterWindowMessageW")

	msgPtr, err := windows.UTF16PtrFromString(name)
	if err != nil {
		return 0, err
	}

	ret, _, err := registerWindowsMessageW.Call(uintptr(unsafe.Pointer(msgPtr)))
	if ret == 0 {
		return 0, err
	}

	return ret, nil
}

// SendBroadcastMessage sends a message trough the broadcast channel
func SendBroadcastMessage(id, p1, p2 uintptr) error {
	if user32DLL == nil {
		once.Do(loadUser32DLL)
	}

	sendMsg := user32DLL.NewProc("SendNotifyMessageW")
	ret, _, err := sendMsg.Call(0xffff, id, p1, p2)

	if ret == 1 {
		return nil
	} else {
		return err
	}
}
