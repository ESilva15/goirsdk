//go:build windows && cgo
// +build windows,cgo

package sharedMem

import (
	"io"
	"log"

	"golang.org/x/sys/windows"
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

// open shared memory. return shmi object.
func open(name string, size uint32) (*shmi, error) {
	return create(name, size)
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
