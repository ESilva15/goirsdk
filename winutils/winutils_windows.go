//go:build windows && cgo
// +build windows,cgo

package winutils

import (
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
	once sync.Once
)

type utils struct {
	user32DLL     *windows.LazyDLL
	wEvent        *windows.Handle
	wBroadcastChn uintptr
}

// INITIALIZATION
func newUtils() (*utils, error) {
	return &utils{
		user32DLL: openUser32DLL(),
	}, nil
}

func (u *utils) Close() {
	closeEvent(u.wEvent)
	// Do we need to unload the user32DLL ???
	// Do we need to close the broadcast channel ???
}

// openEvent opens a windows.Handle for a given event
func (u *utils) OpenEvent(eventName string) error {
	name, err := windows.UTF16PtrFromString(eventName)
	if err != nil {
		return err
	}

	event, err := windows.OpenEvent(windows.SYNCHRONIZE, false, name)
	if err != nil {
		return err
	}
	u.wEvent = &event

	return nil
}

// loadUser32DLL loads the user32.dll which is used to create some processes
func openUser32DLL() *windows.LazyDLL {
	return windows.NewLazyDLL("user32.dll")
}

// OpenBroadcastChannel opens up a broadcast channel to send commands to iracing
func (u *utils) OpenBroadcastChannel(name string) error {
	registerWindowsMessageW := u.user32DLL.NewProc("RegisterWindowMessageW")

	msgPtr, err := windows.UTF16PtrFromString(name)
	if err != nil {
		return err
	}

	ret, _, err := registerWindowsMessageW.Call(uintptr(unsafe.Pointer(msgPtr)))
	if ret == 0 {
		return err
	}
	u.wBroadcastChn = ret

	return nil
}

// INITIALIZATION

// closeEvent closes a given windows.Handle
func closeEvent(h *windows.Handle) {
	windows.CloseHandle(*h)
}

// openEvent waits for a good response for some given time
func (u *utils) CheckValidDataEvent(timeout time.Duration) bool {
	t0 := time.Now().UnixNano()
	timeoutInt := uint32(timeout / time.Millisecond)

	result, err := windows.WaitForSingleObject(*u.wEvent, timeoutInt)
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

// SendBroadcastMessage sends a message trough the broadcast channel
func (u *utils) SendBroadcastMessage(id, p1, p2 uintptr) error {
	sendMsg := u.user32DLL.NewProc("SendNotifyMessageW")
	ret, _, err := sendMsg.Call(0xffff, id, p1, p2)

	if ret == 1 {
		return nil
	} else {
		return err
	}
}
