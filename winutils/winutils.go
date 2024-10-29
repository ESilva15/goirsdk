// I should rename winutils to something else but what this package does
// is interface some windows stuff that we need for the:
// - Broadcast Channel
// - Valid Data Event windows thing
package winutils

import (
	"ibtReader/sharedMem"
	"io"
	"time"
)

type Reader interface {
	io.Reader
	io.ReaderAt
	io.ReadCloser
}

type IRacingWinUtils struct {
	Utils *utils
}

func Init() (*IRacingWinUtils, error) {
	u, err := newUtils()
	if err != nil {
		return nil, err
	}
	return &IRacingWinUtils{u}, nil
}

func (u *IRacingWinUtils) Close() {
	u.Utils.Close()
}

// OpenMemMap returns a Reader interface that can be used to read the data
// No need to encapsulate it
func OpenMemMap(path string, size uint32) (Reader, error) {
	file, err := sharedMem.Open(path, size)
	if err != nil {
		return nil, err
	}

	return file, nil
}

// OpenWinEvent will open the named windows event
func (u *IRacingWinUtils) OpenWinEvent(name string) error {
	return u.Utils.OpenEvent(name)
}

// OpenWinEvent will open the broadcast channel
func (u *IRacingWinUtils) OpenBroadcastChannel(name string) error {
	return u.Utils.OpenBroadcastChannel(name)
}

// CheckValidDataEvent checks if our windows even is telling us we are good to go
func (u *IRacingWinUtils) CheckValidDataEvent(timeout time.Duration) bool {
	return u.Utils.CheckValidDataEvent(timeout)
}
