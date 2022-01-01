package mem

import (
	"errors"
	"fmt"
	"sync"
)


const (
	m4K int = 4096
	m16K int = 16384
	m64K int = 655535
	m256K int = 262144
	m1M int = 1048576
	m4M int = 4194304
	m8M int = 8388608
)

//内存管理池类型
type Pool map[int] *Buf

const (
	//总内存池最大限制 单位是Kb 所以目前限制是 5GB
	EXTRA_MEM_LIMIT int = 5 * 1024 * 1024
)

//单例对象
var bufPoolInstance *BufPool
var once sync.Once


//Buf内存池，为单例模式
type BufPool struct {
	//所有buffer的一个map集合句柄
	Pool Pool
	PoolLock sync.RWMutex

	//总buffer池的内存大小 单位为KB
	TotalMem uint64
}

func MemPool() *BufPool{
	once.Do(func() {
		bufPoolInstance = new(BufPool)
		bufPoolInstance.Pool = make(map[int]*Buf)
		bufPoolInstance.TotalMem = 0
		bufPoolInstance.initPool()
	})

	return bufPoolInstance
}

func (bp *BufPool) makeBufList(cap int, num int) {
	bp.Pool[cap] = NewBuf(cap)

	var prev *Buf
	prev = bp.Pool[cap]
	for i := 1; i < num; i ++ {
		prev.Next = NewBuf(cap)
		prev = prev.Next
	}
	bp.TotalMem += (uint64(cap)/1024) * uint64(num)
}

/*
	初始化内存池 主要是预先开辟一定量的空间
	这里BufPool是一个hash，每个key都是不同空间容量
	对应的value是一个Buf集合的链表

BufPool -->   [m4K] -- Buf-Buf-Buf-Buf...(BufList)
              [m16K] -- Buf-Buf-Buf-Buf...(BufList)
              [m64K] -- Buf-Buf-Buf-Buf...(BufList)
              [m256K] -- Buf-Buf-Buf-Buf...(BufList)
              [m1M] -- Buf-Buf-Buf-Buf...(BufList)
              [m4M] -- Buf-Buf-Buf-Buf...(BufList)
              [m8M] -- Buf-Buf-Buf-Buf...(BufList)
 */
func (bp *BufPool) initPool() {
	//----> 开辟4K buf 内存池
	// 4K的Buf 预先开辟5000个，约20MB供开发者使用
	bp.makeBufList(m4K, 5000)

	//----> 开辟16K buf 内存池
	//16K的Buf 预先开辟1000个，约16MB供开发者使用
	bp.makeBufList(m16K, 1000)

	//----> 开辟64K buf 内存池
	//64K的Buf 预先开辟500个，约32MB供开发者使用
	bp.makeBufList(m64K, 500)

	//----> 开辟256K buf 内存池
	//256K的Buf 预先开辟200个，约50MB供开发者使用
	bp.makeBufList(m256K, 200)

	//----> 开辟1M buf 内存池
	//1M的Buf 预先开辟50个，约50MB供开发者使用
	bp.makeBufList(m1M, 50)

	//----> 开辟4M buf 内存池
	//4M的Buf 预先开辟20个，约80MB供开发者使用
	bp.makeBufList(m4M, 20)

	//----> 开辟8M buf 内存池
	//8M的io_buf 预先开辟10个，约80MB供开发者使用
	bp.makeBufList(m8M, 10)
}


//当前内存池的总内存容量（KB单位）
func (bp *BufPool) MemSize() uint64 {
	return bp.TotalMem
}

/*
	开辟一个Buf
	1 如果上层需要N个字节的大小的空间，找到与N最接近的buf hash组，取出，
	2 如果该组已经没有节点使用，可以额外申请
	3 总申请长度不能够超过最大的限制大小 EXTRA_MEM_LIMIT
	4 如果有该节点需要的内存块，直接取出，并且将该内存块从pool摘除
*/
func (bp *BufPool) Alloc(N int) (*Buf, error) {
	//1 找到N最接近哪hash 组
	var index int
	if N <= m4K {
		index = m4K
	} else if (N <= m16K) {
		index = m16K
	} else if (N <= m64K) {
		index = m64K
	} else if (N <= m256K) {
		index = m256K
	} else if (N <= m1M) {
		index = m1M
	} else if (N <= m4M) {
		index = m4M
	} else if (N <= m8M) {
		index = m8M
	} else {
		return nil, errors.New("Alloc size Too Large!");
	}

	//2 如果该组已经没有，需要额外申请，那么需要加锁保护
	bp.PoolLock.Lock()
	if bp.Pool[index] == nil {
		if (bp.TotalMem + uint64(index/1024)) >= uint64(EXTRA_MEM_LIMIT) {
			errStr := fmt.Sprintf("already use too many memory!\n")
			return nil, errors.New(errStr)
		}

		newBuf := NewBuf(index)
		bp.TotalMem += uint64(index/1024)
		bp.PoolLock.Unlock()
		fmt.Printf("Alloc Mem Size: %d KB\n", newBuf.Capacity/1024)
		return newBuf, nil
	}

	//3 如果有该组有Buf内存存在，那么得到一个Buf并返回，并且从pool中摘除该内存块
	targetBuf := bp.Pool[index]
	bp.Pool[index] = targetBuf.Next
	bp.TotalMem -= uint64(index/1024)
	bp.PoolLock.Unlock()
	targetBuf.Next = nil
	fmt.Printf("Alloc Mem Size: %d KB\n", targetBuf.Capacity/1024)
	return targetBuf, nil
}


//当Alloc之后，当前Buf被使用完，需要重置这个Buf,需要将该buf放回pool中
func (bp *BufPool) Revert(buf *Buf) error {
	//每个buf的容量都是固定的 在hash的key中取值
	index := buf.Capacity
	//重置io_buf中的内置位置指针
	buf.Clear()

	bp.PoolLock.Lock()
	//找到对应的hash组 buf首届点地址
	if _, ok := bp.Pool[index]; !ok {
		errStr := fmt.Sprintf("Index %d not in BufPoll!\n", index)
		return errors.New(errStr)
	}

	//将buffer插回链表头部
	buf.Next = bp.Pool[index]
	bp.Pool[index] = buf
	bp.TotalMem += uint64(index/1024)
	bp.PoolLock.Unlock()
	fmt.Printf("Revert Mem Size: %d KB\n",index/1024)

	return nil
}