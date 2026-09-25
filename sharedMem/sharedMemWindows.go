//go:build windows && cgo
// +build windows,cgo

package sharedMem

import (
	"io"
	"log"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	modkernel32          = windows.NewLazySystemDLL("kernel32.dll")
	procOpenFileMappingW = modkernel32.NewProc("OpenFileMappingW")
)

type shmi struct {
	h    windows.Handle
	v    uintptr
	size uint32
}

// create shared memory. return shmi object.
func create(name string, size uint32) (*shmi, error) {
	fnPtr, _ := windows.UTF16PtrFromString(name)

	flProtect := uint32(windows.PAGE_READONLY)

	h, errno := windows.CreateFileMapping(
		windows.InvalidHandle,
		nil,
		flProtect,
		0,
		size,
		fnPtr)
	if h == 0 {
		log.Fatal("could not open memmap file: ", errno)
	}

	addr, errno := windows.MapViewOfFile(h,
		windows.FILE_MAP_READ,
		0,
		0,
		uintptr(size))
	if addr == 0 {
		log.Printf("error in MapViewOfFile: %v", errno)
	}

	return &shmi{h, addr, size}, nil
}

func openFileMapping(desidredAccess uint32, inheritHandle bool, name *uint16) (windows.Handle, error) {
	var inherit uint32
	if inheritHandle {
		inherit = 1
	}

	r1, _, err := procOpenFileMappingW.Call(
		uintptr(desidredAccess),
		uintptr(inherit),
		uintptr(unsafe.Pointer(name)),
	)
	if r1 == 0 {
		return 0, err
	}

	return windows.Handle(r1), nil
}

// open shared memory. return shmi object.
func open(name string, size uint32) (*shmi, error) {
	fnPtr, err := windows.UTF16PtrFromString(name)
	if err != nil {
		return nil, err
	}

	memHandle, err := openFileMapping(windows.FILE_MAP_READ, false, fnPtr)
	if err != nil {
		return nil, err
	}

	memAddr, err := windows.MapViewOfFile(memHandle, windows.FILE_MAP_READ, 0, 0, 0)
	if err != nil {
		return nil, err
	}

	return &shmi{
		h:    memHandle,
		v:    memAddr,
		size: size,
	}, nil
}

func (o *shmi) close() error {
	if o.v != uintptr(0) {
		windows.UnmapViewOfFile(o.v)
		o.v = uintptr(0)
	}
	if o.h != windows.InvalidHandle {
		windows.CloseHandle(o.h)
		o.h = windows.InvalidHandle
	}
	return nil
}

// read shared memory. return read size.
func (o *shmi) readAt(p []byte, off int64) (n int, err error) {
	if off >= int64(o.size) {
		return 0, io.EOF
	}
	if max := int64(o.size) - off; int64(len(p)) > max {
		p = p[:max]
	}
	return copyPtr2Slice(o.v, p, off, o.size), nil
}

// write shared memory. return write size.
func (o *shmi) writeAt(p []byte, off int64) (n int, err error) {
	if off >= int64(o.size) {
		return 0, io.EOF
	}
	if max := int64(o.size) - off; int64(len(p)) > max {
		p = p[:max]
	}
	return copySlice2Ptr(p, o.v, off, o.size), nil
}
