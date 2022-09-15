- # 网络轮询器(NetPoller) 
网络轮询器不仅用于监控网络 `I/O`，还能用于监控文件的 `I/O`，
它利用了操作系统提供的 `I/O` 多路复用模型来提升 `I/O` 设备的利用率以及程序的性能。  

- [类型](#类型)
- [文件](#文件)
- [网络轮询器的实现](#网络轮询器的实现)
- [Unix Socket](#unix-socket)
## 类型
操作系统中包含阻塞 `I/O`、非阻塞 `I/O`、信号驱动 `I/O` 与异步 `I/O` 以及 `I/O` 多路复用五种 `I/O` 模型。我们在本节中会介绍上述五种模型中的三种：

- 阻塞 `I/O` 模型；
- 非阻塞 `I/O` 模型；
- `I/O` 多路复用模型；

在 Unix 和类 Unix 操作系统中，文件描述符（File descriptor，FD）
是用于访问文件或者其他 `I/O` 资源的抽象句柄，例如：管道或者网络套接字1。而不同的 `I/O` 模型会使用不同的方式操作文件描述符。


## 文件  
文件的定义:
```go
// os/file.go
// File represents an open file descriptor.
type File struct {
	*file // os specific
}

//os/file_unix.go
// file is the real representation of *File.
// The extra level of indirection ensures that no clients of os
// can overwrite this data, which could cause the finalizer
// to close the wrong file descriptor.
type file struct {
	pfd         poll.FD
	name        string
	dirinfo     *dirInfo // nil unless directory being read
	nonblock    bool     // whether we set nonblocking mode
	stdoutOrErr bool     // whether this is stdout or stderr
	appendMode  bool     // whether file is opened for appending
}
```

重要的是`poll.FD`, 系统文件及网络都是基于它实现的
```go
// FD is a file descriptor. The net and os packages use this type as a
// field of a larger type representing a network connection or OS file.
type FD struct {
	// 锁定 sysfd 并序列化对 Read 和 Write 方法的访问。
	fdmu fdMutex

	// 系统文件描述符。 在关闭之前不可变。
	Sysfd int

	// I/O 轮询器(poller)
	pd pollDesc

	// 写入缓存
	iovecs *[]syscall.Iovec

	// 当文件关闭时发出信号。
	csema uint32

	// 如果此文件已设置为阻塞模式，则非零。
	isBlocking uint32

	// 这是否是流式描述符，而不是像 UDP 套接字这样的基于数据包的描述符。 不可变。
	IsStream bool

	// 零字节读取是否表示 EOF。 对于基于消息的套接字连接，这是错误的。
	ZeroReadIsEOF bool

	// 这是否是文件而不是网络套接字。
	isFile bool
}
```
网络轮询器的关键在于`pd pollDesc`.源码位置`internal/poll/fd_poll_runtime.go`
```go
type pollDesc struct {
	runtimeCtx uintptr
}
```
这个`runtimeCtx`指向的就是`FD`  
```go
func (pd *pollDesc) init(fd *FD) error {
	serverInit.Do(runtime_pollServerInit)
	ctx, errno := runtime_pollOpen(uintptr(fd.Sysfd))
	if errno != 0 {
		if ctx != 0 {
			runtime_pollUnblock(ctx)
			runtime_pollClose(ctx)
		}
		return errnoErr(syscall.Errno(errno))
	}
	pd.runtimeCtx = ctx
	return nil
}
```


## 网络轮询器的实现 
`pollDesc`拥有的方法有: 

```go
func (pd *pollDesc) init(fd *FD) error
func (pd *pollDesc) close()
func (pd *pollDesc) evict()
func (pd *pollDesc) prepare(mode int, isFile bool) error
func (pd *pollDesc) prepareRead(isFile bool) error
func (pd *pollDesc) prepareWrite(isFile bool) error
func (pd *pollDesc) wait(mode int, isFile bool) error
func (pd *pollDesc) waitRead(isFile bool) error
func (pd *pollDesc) waitWrite(isFile bool) error
func (pd *pollDesc) waitCanceled(mode int)
func (pd *pollDesc) pollable()
```

最终调用
```go
func runtime_pollServerInit()
func runtime_pollOpen(fd uintptr) (uintptr, int)
func runtime_pollClose(ctx uintptr)
func runtime_pollWait(ctx uintptr, mode int) int
func runtime_pollWaitCanceled(ctx uintptr, mode int) int
func runtime_pollReset(ctx uintptr, mode int) int
func runtime_pollSetDeadline(ctx uintptr, d int64, mode int)
func runtime_pollUnblock(ctx uintptr)
func runtime_isPollServerDescriptor(fd uintptr) bool
```

具体的实现在`runtime/netpoll.go`，先看下`pollDesc`结构体定义:  
```go
// Network poller descriptor.
//
// No heap pointers.
//
//go:notinheap
type pollDesc struct {
	link *pollDesc // in pollcache, protected by pollcache.lock

	// The lock protects pollOpen, pollSetDeadline, pollUnblock and deadlineimpl operations.
	// This fully covers seq, rt and wt variables. fd is constant throughout the PollDesc lifetime.
	// pollReset, pollWait, pollWaitCanceled and runtime·netpollready (IO readiness notification)
	// proceed w/o taking the lock. So closing, everr, rg, rd, wg and wd are manipulated
	// in a lock-free way by all operations.
	// NOTE(dvyukov): the following code uses uintptr to store *g (rg/wg),
	// that will blow up when GC starts moving objects.
	lock    mutex // protects the following fields
	fd      uintptr
	closing bool
	everr   bool      // marks event scanning error happened
	user    uint32    // user settable cookie
	rseq    uintptr   // protects from stale read timers
	rg      uintptr   // pdReady, pdWait, G waiting for read or nil
	rt      timer     // read deadline timer (set if rt.f != nil)
	rd      int64     // read deadline
	wseq    uintptr   // protects from stale write timers
	wg      uintptr   // pdReady, pdWait, G waiting for write or nil
	wt      timer     // write deadline timer
	wd      int64     // write deadline
	self    *pollDesc // storage for indirect interface. See (*pollDesc).makeArg.
}
```

目前所有重要的数据结构都已出炉，接下来就要看看具体实现  
还是先看看源码中的注释`runtime/netpoll.go`  
```
// Integrated network poller (platform-independent part).
// A particular implementation (epoll/kqueue/port/AIX/Windows)
// must define the following functions:
//
// func netpollinit()
//     Initialize the poller. Only called once.
//
// func netpollopen(fd uintptr, pd *pollDesc) int32
//     Arm edge-triggered notifications for fd. The pd argument is to pass
//     back to netpollready when fd is ready. Return an errno value.
//
// func netpoll(delta int64) gList
//     Poll the network. If delta < 0, block indefinitely. If delta == 0,
//     poll without blocking. If delta > 0, block for up to delta nanoseconds.
//     Return a list of goroutines built by calling netpollready.
//
// func netpollBreak()
//     Wake up the network poller, assumed to be blocked in netpoll.
//
// func netpollIsPollDescriptor(fd uintptr) bool
//     Reports whether fd is a file descriptor used by the poller.
```

google翻译
```
网络轮询器与。特定实现（epoll/kqueue/port/AIX/Windows）必须定义以下函数：

// func netpollinit()
初始化轮询器。 只调用一次。

// func netpollopen(fd uintptr, pd *pollDesc) int32
为 fd 设置边缘触发通知。 pd 参数是在 fd 准备好时传回 netpollready。 返回一个 errno 值。

// func netpoll(delta int64) gList
轮询网络。 如果 delta < 0，则无限期阻塞。 如果 delta == 0，则轮询而不阻塞。 如果 delta > 0，最多阻塞 delta 纳秒。返回通过调用 netpollready 构建的 goroutine 列表。

// 函数 netpollBreak()
唤醒网络轮询器，假设在 netpoll 中被阻塞。

// func netpollIsPollDescriptor(fd uintptr) bool
报告 fd 是否是轮询器使用的文件描述符。 
```

## Unix Socket


client
```go
//客户端
func (this *UnixSocket) ClientSendContext(context string) {
	addr, err := net.ResolveUnixAddr("unixgram", this.filename)
	if err != nil {
		panic("Cannot resolve unix addr: " + err.Error())
	}
	//拔号
	c, err := net.DialUnix("unixgram", nil, addr)
	if err != nil {
		panic("DialUnix failed.")
	}
	//写出
	_, err = c.Write([]byte(context))
	if err != nil {
		panic("Writes failed.")
	}
}
Footer

```

server
```go
func (this *UnixSocket) createServer() {
	fmt.Println("socket监听执行========================================")
	os.Remove(this.filename)
	addr, err := net.ResolveUnixAddr("unixgram", this.filename)
	if err != nil {
		panic("Cannot resolve unix addr: " + err.Error())
	}
	c, err := net.ListenUnixgram("unixgram", addr)
	defer c.Close()
	if err != nil {
		panic("Cannot listen to unix domain socket: " + err.Error())
	}
	os.Chmod(this.filename, 0666)
	for {
		data := make([]byte, 4096)
		nr, _, err := c.ReadFrom(data)
		if err != nil {
			fmt.Printf("conn.ReadFrom error: %s\n", err)
			return
		}
		go this.HandleServerConn(c, string(data[0:nr]))
	}

}
```

`net/unixsock.go` 源码  
```go
// ResolveUnixAddr returns an address of Unix domain socket end point.
//
// The network must be a Unix network name.
//
// See func Dial for a description of the network and address
// parameters.
func ResolveUnixAddr(network, address string) (*UnixAddr, error) {
	switch network {
	case "unix", "unixgram", "unixpacket":
		return &UnixAddr{Name: address, Net: network}, nil
	default:
		return nil, UnknownNetworkError(network)
	}
}

// ListenUnixgram acts like ListenPacket for Unix networks.
//
// The network must be "unixgram".
func ListenUnixgram(network string, laddr *UnixAddr) (*UnixConn, error) {
	switch network {
	case "unixgram":
	default:
		return nil, &OpError{Op: "listen", Net: network, Source: nil, Addr: laddr.opAddr(), Err: UnknownNetworkError(network)}
	}
	if laddr == nil {
		return nil, &OpError{Op: "listen", Net: network, Source: nil, Addr: nil, Err: errMissingAddress}
	}
	sl := &sysListener{network: network, address: laddr.String()}
	c, err := sl.listenUnixgram(context.Background(), laddr)
	if err != nil {
		return nil, &OpError{Op: "listen", Net: network, Source: nil, Addr: laddr.opAddr(), Err: err}
	}
	return c, nil
}

func (sl *sysListener) listenUnixgram(ctx context.Context, laddr *UnixAddr) (*UnixConn, error) {
	fd, err := unixSocket(ctx, sl.network, laddr, nil, "listen", sl.ListenConfig.Control)
	if err != nil {
		return nil, err
	}
	return newUnixConn(fd), nil
}

func newUnixConn(fd *netFD) *UnixConn { return &UnixConn{conn{fd}} }
```

















