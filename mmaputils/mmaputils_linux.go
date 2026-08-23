//go:build (linux && cgo) || (darwin && cgo)

package mmaputils

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/ESilva15/goirsdk/sharedMem"
)

type utils struct {
	socketPath string
	listener   *net.UnixConn
}

func newUtils() (*utils, error) {
	return &utils{}, nil
}

func (u *utils) Close() {
	if u.listener != nil {
		u.listener.Close()
		os.Remove(u.socketPath)
	}
}

// OpenMemMap returns a Reader interface that can be used to read the data
// No need to encapsulate it
func OpenMemMap(name string, size uint32) (Reader, error) {
	file, err := sharedMem.Open(name, size)
	if err != nil {
		return nil, fmt.Errorf("failed to open `%s` with err: %+v", name, err)
	}

	return file, nil
}

// OpenEvent creates or connects to a Unix socket for event signaling on Linux
func (u *utils) OpenEvent(eventName string) error {
	u.socketPath = fmt.Sprintf("/tmp/iracing_%s.sock", eventName)
	_ = os.Remove(u.socketPath)

	addr, err := net.ResolveUnixAddr("unixgram", u.socketPath)
	if err != nil {
		return err
	}

	l, err := net.ListenUnixgram("unixgram", addr)
	if err != nil {
		return err
	}
	u.listener = l
	return nil
}

func (u *utils) OpenBroadcastChannel(name string) error {
	// No-op or log stub on Linux
	return nil
}

// CheckValidDataEvent waits for a pulse byte sent over the socket
func (u *utils) CheckValidDataEvent(timeout time.Duration) bool {
	if u.listener == nil {
		return false
	}

	_ = u.listener.SetReadDeadline(time.Now().Add(timeout))
	buf := make([]byte, 1)
	_, err := u.listener.Read(buf)
	return err == nil
}

func (u *utils) SendBroadcastMessage(id, p1, p2 uintptr) error {
	return nil
}

// package winutils
//
// import (
// 	"errors"
// 	"sync"
// 	"time"
// )
//
// const (
// 	WAIT_OBJECT_0 = 0
// 	WAIT_TIMEOUT  = 258
// )
//
// var (
// 	once             sync.Once
// 	ErrUnsupportedOS = errors.New("not found")
// )
//
// type utils struct {
// }
//
// // INITIALIZATION
// func newUtils() (*utils, error) {
// 	return nil, ErrUnsupportedOS
// }
//
// func (u *utils) Close() {
// }
//
// // openEvent opens a windows.Handle for a given event
// func (u *utils) OpenEvent(eventName string) error {
// 	return ErrUnsupportedOS
// }
//
// // OpenBroadcastChannel opens up a broadcast channel to send commands to iracing
// func (u *utils) OpenBroadcastChannel(name string) error {
// 	return ErrUnsupportedOS
// }
//
// // INITIALIZATION
//
// // openEvent waits for a good response for some given time
// func (u *utils) CheckValidDataEvent(timeout time.Duration) bool {
// 	return false
// }
//
// // SendBroadcastMessage sends a message trough the broadcast channel
// func (u *utils) SendBroadcastMessage(id, p1, p2 uintptr) error {
//   return ErrUnsupportedOS
// }
