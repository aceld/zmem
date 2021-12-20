package mem

import "C"
import (
	"ZMM/c"
	"fmt"
)
import "unsafe"

type Buf struct {
	//如果存在多个buffer，是采用链表的形式链接起来
	Next *Buf
	//当前buffer的缓存容量大小
	Capacity int
	//当前buffer有效数据长度
	length int
	//未处理数据的头部位置索引
	head int
	//当前io_buf所保存的数据地址
	data unsafe.Pointer
}


//构造，创建一个IoBuf对象
func NewIoBuf(size int) *Buf {
	return &Buf{
		Capacity: size,
		length: 0,
		head: 0,
		Next: nil,
		data : c.Malloc(size),
	}
}

//清空数据
func (b *Buf) Clear() {
	b.length = 0
	b.head = 0
}

//将已经处理过的数据，清空,将未处理的数据提前至数据首地址
func (b *Buf) Adjust() {
	if b.head != 0 {
		if (b.length != 0) {
			c.Memmove(b.data, unsafe.Pointer(uintptr(b.data) + uintptr(b.head)), b.length)
		}
		b.head = 0
	}
}

//将其他Buf对象数据考本到自己中
func (b *Buf) Copy(other *Buf) {
	c.Memcpy(b.data, other.GetBytes(), other.length)
	b.head = 0
	b.length = other.length
}

//处理长度为len的数据，移动head和修正length
func (b *Buf) Pop(len int) {
	if b.data == nil {
		fmt.Printf("pop data is nil")
		return
	}
	if len > b.length {
		fmt.Printf("pop len > length")
		return
	}
	b.length -= len
	b.head += len
}

//给一个Buf填充[]byte数据
func (b *Buf) SetBytes(src []byte) {
	c.Memcpy(unsafe.Pointer(uintptr(b.data)+uintptr(b.head)), src, len(src))
	b.length += len(src)
}

//获取一个Buf的数据，以[]byte形式展现
func (b *Buf) GetBytes() []byte {
	data := C.GoBytes(unsafe.Pointer(uintptr(b.data)+uintptr(b.head)), C.int(b.length))
	return data
}

func (b *Buf) Head() int {
	return b.head
}

func (b *Buf) Length() int {
	return b.length
}