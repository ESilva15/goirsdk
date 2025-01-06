//go:build (linux && cgo) || (darwin && cgo)
// +build linux,cgo darwin,cgo

package winutils

import (
	"errors"
	"sync"
	"time"
)

const (
	WAIT_OBJECT_0 = 0
	WAIT_TIMEOUT  = 258
)

var (
	once             sync.Once
	ErrUnsupportedOS = errors.New("not found")
)

type utils struct {
}

// INITIALIZATION
func newUtils() (*utils, error) {
	return nil, ErrUnsupportedOS
}

func (u *utils) Close() {
}

// openEvent opens a windows.Handle for a given event
func (u *utils) OpenEvent(eventName string) error {
	return ErrUnsupportedOS
}

// OpenBroadcastChannel opens up a broadcast channel to send commands to iracing
func (u *utils) OpenBroadcastChannel(name string) error {
	return ErrUnsupportedOS
}

// INITIALIZATION

// openEvent waits for a good response for some given time
func (u *utils) CheckValidDataEvent(timeout time.Duration) bool {
	return false
}

// SendBroadcastMessage sends a message trough the broadcast channel
func (u *utils) SendBroadcastMessage(id, p1, p2 uintptr) error {
  return ErrUnsupportedOS
}
