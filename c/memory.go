package c

/*
#include <string.h>
#include <stdlib.h>
 */
import "C"
import "unsafe"

func Malloc(size int) (unsafe.Pointer) {
	return C.malloc(C.size_t(size))
}

func Free(data unsafe.Pointer) {
	C.free(data)
}

func Memmove(dest, src unsafe.Pointer, length int) {
	C.memmove(dest, src, C.size_t(length))
}

func Memcpy(dest unsafe.Pointer, src []byte, length int) {
	srcData := C.CBytes(src)
	C.memcpy(dest, srcData, C.size_t(length))
}