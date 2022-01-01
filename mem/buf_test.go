package mem_test

import (
	"zmem/mem"
	"fmt"
	"testing"
)


func TestBufPoolCopy(t *testing.T) {
	pool := mem.MemPool()

	buffer, err := pool.Alloc(1)
	if err != nil {
		fmt.Println("pool Alloc Error ", err)
		return
	}

	buffer.SetBytes([]byte("Aceld12345"))
	fmt.Printf("Buffer GetBytes = %+v\n", string(buffer.GetBytes()))

	buffer2, err := pool.Alloc(1)
	if err != nil {
		fmt.Println("pool Alloc Error ", err)
		return
	}
	buffer2.Copy(buffer)
	fmt.Printf("Buffer2 GetBytes = %+v\n", string(buffer2.GetBytes()))

	buffer2.Pop(5)
	fmt.Printf("Buffer2 GetBytes = %+v\n", string(buffer2.GetBytes()))
}


func TestBufPoolSetGet(t *testing.T) {
	pool := mem.MemPool()

	buffer, err := pool.Alloc(1)
	if err != nil {
		fmt.Println("pool Alloc Error ", err)
		return
	}

	buffer.SetBytes([]byte("Aceld12345"))
	fmt.Printf("GetBytes = %+v, ToString = %s\n", buffer.GetBytes(), string(buffer.GetBytes()))
	buffer.Pop(4)
	fmt.Printf("GetBytes = %+v, ToString = %s\n", buffer.GetBytes(), string(buffer.GetBytes()))
}

func TestBufPoolAdjust(t *testing.T) {
	pool := mem.MemPool()

	buffer, err := pool.Alloc(4096)
	if err != nil {
		fmt.Println("pool Alloc Error ", err)
		return
	}

	buffer.SetBytes([]byte("Aceld12345"))
	fmt.Printf("GetBytes = %+v, Head = %d, Length = %d\n", buffer.GetBytes(), buffer.Head(), buffer.Length())
	buffer.Pop(4)
	fmt.Printf("GetBytes = %+v, Head = %d, Length = %d\n", buffer.GetBytes(), buffer.Head(), buffer.Length())
	buffer.Adjust()
	fmt.Printf("GetBytes = %+v, Head = %d, Length = %d\n", buffer.GetBytes(), buffer.Head(), buffer.Length())
}





