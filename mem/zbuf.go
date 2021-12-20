package mem

import (
	"fmt"
	"net"
)

//应用层的buffer数据
type ZBuf struct {
	b *Buf
}

//清空当前的ZBuf
func (zb *ZBuf) Clear() {
	if zb.b != nil {
		//将Buf重新放回到buf_pool中
		MemPool().Revert(zb.b)
		zb.b = nil
	}
}

//弹出已使用的有效长度
func (zb *ZBuf) Pop(len int) {
	if zb.b == nil || len > zb.b.Length() {
		return
	}

	zb.b.Pop(len)

	//当此时Buf的可用长度已经为0时,将Buf重新放回BufPool中
	if zb.b.Length() == 0 {
		MemPool().Revert(zb.b)
		zb.b = nil
	}
}

//获取Buf中的数据
func (zb *ZBuf) Data() []byte {
	if zb.b == nil {
		return nil
	}
	return zb.b.GetBytes()
}

//重置缓冲区
func (zb *ZBuf) Adjust() {
	if zb.b != nil {
		zb.b.Adjust()
	}
}

//读取数据到Buf中
func (zb *ZBuf) Read(src []byte) (err error){
	if zb.b == nil {
		zb.b, err = MemPool().Alloc(len(src))
		if err != nil {
			fmt.Println("pool Alloc Error ", err)
		}
	} else {
		if zb.b.Head() != 0 {
			return nil
		}
		if zb.b.Capacity - zb.b.Length() < len(src) {
			//不够存，重新从内存池申请
			newBuf, err := MemPool().Alloc(len(src)+zb.b.Length())
			if err != nil {
				return nil
			}
			//将之前的Buf拷贝到新申请的Buf中去
			newBuf.Copy(zb.b)
			//将之前的Buf回收到内存池中
			MemPool().Revert(zb.b)
			//新申请的Buf成为当前的ZBuf
			zb.b = newBuf
		}
	}

	//将内容写进ZBuf缓冲中
	zb.b.SetBytes(src)

	return nil
}


//读取数据从Fd文件描述符中
func (zb *ZBuf) ReadFromFd(fd int) error {
	//...
	return nil
}

//将数据写入Fd文件描述符中
func (zb *ZBuf) WriteToFd(fd int) error {
	//...
	return nil
}

//读取数据从文件中
func (zb *ZBuf) ReadFromFile(path string) error {
	//...
	return nil
}

func (zb *ZBuf) WriteToFile(path string) error {
	//...
	return nil
}

//读取数据从网络连接中
func (zb *ZBuf) ReadFromConn(conn net.Conn) error {
	//...
	return nil
}

func (zb *ZBuf) WriteToConn(conn net.Conn) error {
	//...
	return nil
}





