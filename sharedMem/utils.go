package sharedMem

import (
	"unsafe"
)

func copySlice2Ptr(b []byte, p uintptr, off int64, size uint32) int {
	bb := unsafe.Slice((*byte)(*(*unsafe.Pointer)(unsafe.Pointer(&p))), int(size))
	return copy(bb[off:], b)
}

func copyPtr2Slice(p uintptr, b []byte, off int64, size uint32) int {
	bb := unsafe.Slice((*byte)(*(*unsafe.Pointer)(unsafe.Pointer(&p))), int(size))
	return copy(b, bb[off:size])
}
