package mem

import (
	"fmt"
	"testing"
)

func TestZBuf(t *testing.T) {
	buf := new(ZBuf)
	buf.Read([]byte("12345Aceld"))
	buf.Pop(5)
	fmt.Printf("data := %+v\n", string(buf.Data()))
	buf.Clear()
}
