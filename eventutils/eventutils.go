// I should rename winutils to something else but what this package does
// is interface some windows stuff that we need for the:
// - Broadcast Channel
// - Valid Data Event windows thing
package mmaputils

import (
	"io"
	"time"
)

type Reader interface {
	io.Reader
	io.ReaderAt
	io.ReadCloser
}

type EventUtils struct {
	Utils *utils
}

func Init() (*EventUtils, error) {
	u, err := newUtils()
	if err != nil {
		return nil, err
	}
	return &EventUtils{u}, nil
}

func (u *EventUtils) Close() {
	if u.Utils == nil {
		return
	}

	u.Utils.Close()
}

// OpenEvent will open the named windows event
func (u *EventUtils) OpenEvent(name string) error {
	return u.Utils.OpenEvent(name)
}

// OpenBroadcastChannel will open the broadcast channel
func (u *EventUtils) OpenBroadcastChannel(name string) error {
	return u.Utils.OpenBroadcastChannel(name)
}

// CheckValidDataEvent checks if our windows even is telling us we are good to go
func (u *EventUtils) CheckValidDataEvent(timeout time.Duration) bool {
	return u.Utils.CheckValidDataEvent(timeout)
}

// SignalEvent triggers a frame pulse (used by mock servers and tests)
func (u *EventUtils) SignalEvent() error {
	return u.Utils.SignalEvent()
}

// SignalEvent triggers a frame pulse (used by mock servers and tests)
func SignalEvent(name string) {
	signalEvent(name)
}

// CleanupEvent cleans up event resources (used by mock servers and tests)
func CleanupEvent(name string) {
	cleanupEvent(name)
}
