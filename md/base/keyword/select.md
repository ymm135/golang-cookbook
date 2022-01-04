# Select实现  
[参考文章](https://draveness.me/golang/docs/part2-foundation/ch05-keyword/golang-select/)  
`select` 是操作系统中的系统调用，我们经常会使用 `select`、`poll` 和 `epoll` 等函数构建
I/O 多路复用模型提升程序的性能。Go 语言的 `select` 与操作系统中的 `select` 比较相似。

## 测试demo
[code](code/base/keyword/select/base-select.go)  
```go
package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan string)

	go func() {
		ch <- "hello channel"
		ch <- "!"
		ch <- "quit"
	}()

	//需要延迟一段时间，要不然接收不到数据
	time.Sleep(time.Second)

	for { //
		select {
		case str := <-ch:
			fmt.Println(str)
			if str == "quit" {
				goto end
			}
		default:
			fmt.Println("default")
		}
	}
end:
}
```

输出结果:
```shell
hello channel
!
default
default
default
default
default
default
default
default
default
default
default
default
default
default
default
default
default
quit
```

## 数据结构 
`chan`的数据结构在`go/src/runtime/chan.go`
```go
type hchan struct {
	qcount   uint           // total data in the queue
	dataqsiz uint           // size of the circular queue
	buf      unsafe.Pointer // points to an array of dataqsiz elements
	elemsize uint16
	closed   uint32
	elemtype *_type // element type
	sendx    uint   // send index
	recvx    uint   // receive index
	recvq    waitq  // list of recv waiters
	sendq    waitq  // list of send waiters

	// lock protects all fields in hchan, as well as several
	// fields in sudogs blocked on this channel.
	//
	// Do not change another G's status while holding this lock
	// (in particular, do not ready a G), as this can deadlock
	// with stack shrinking.
	lock mutex
}
```

## 实现原理 
`select` 语句在编译期间会被转换成 `OSELECT` 节点。每个 `OSELECT` 节点都会持有一组 `OCASE` 节点，
每一个 `OCASE` 既包含执行条件也包含满足条件后执行的代码。如果 `OCASE` 的执行条件是空，那就意味着这是一个 `default` 节点。  

具体实现代码在`cmd/compile/internal/gc/select.go`的`func walkselectcases(cases *Nodes) []*Node` 方法  
```go
ncas := cases.Len()
	sellineno := lineno

	// optimization: zero-case select
	if ncas == 0 {
		return []*Node{mkcall("block", nil, nil)}
	}

	// optimization: one-case select: single op.
	if ncas == 1 {
        cas := cases.First()
        setlineno(cas)
        l := cas.Ninit.Slice()
        if cas.Left != nil { // not default:
        n := cas.Left
        l = append(l, n.Ninit.Slice()...)
        n.Ninit.Set(nil)
        switch n.Op {
            default:
            Fatalf("select %v", n.Op)
            
            case OSEND:
            // already ok
            
            case OSELRECV, OSELRECV2:
		}
		l = append(l, cas.Nbody.Slice()...)
        l = append(l, nod(OBREAK, nil, nil))
        return l
	}

    // optimization: two-case select but one is default: single non-blocking op.
    if ncas == 2 && dflt != nil {
        ...
		
		case OSEND:
            // if selectnbsend(c, v) { body } else { default body }
            ch := n.Left
            r.Left = mkcall1(chanfn("selectnbsend", 2, ch.Type), types.Types[TBOOL], &r.Ninit, ch, n.Right)

        case OSELRECV:
            // if selectnbrecv(&v, c) { body } else { default body }
            ch := n.Right.Left
            elem := n.Left
            if elem == nil {
            elem = nodnil()
            }
            r.Left = mkcall1(chanfn("selectnbrecv", 2, ch.Type), types.Types[TBOOL], &r.Ninit, elem, ch)
			...
    }
``` 

在编译阶段会根据具体情况做不同的优化，我们这里关注`channel`的收发情况。`go/src/runtime/chan.go`源码中写的很明白  
```go
// compiler implements
//
//	select {
//	case c <- v:
//		... foo
//	default:
//		... bar
//	}
//
// as
//
//	if selectnbsend(c, v) {
//		... foo
//	} else {
//		... bar
//	}
//
func selectnbsend(c *hchan, elem unsafe.Pointer) (selected bool) {
return chansend(c, elem, false, getcallerpc())
}

// compiler implements
//
//	select {
//	case v = <-c:
//		... foo
//	default:
//		... bar
//	}
//
// as
//
//	if selectnbrecv(&v, c) {
//		... foo
//	} else {
//		... bar
//	}
//
func selectnbrecv(elem unsafe.Pointer, c *hchan) (selected bool) {
selected, _ = chanrecv(c, elem, false)
return
}

// compiler implements
//
//	select {
//	case v, ok = <-c:
//		... foo
//	default:
//		... bar
//	}
//
// as
//
//	if c != nil && selectnbrecv2(&v, &ok, c) {
//		... foo
//	} else {
//		... bar
//	}
//
func selectnbrecv2(elem unsafe.Pointer, received *bool, c *hchan) (selected bool) {
// TODO(khr): just return 2 values from this function, now that it is in Go.
selected, *received = chanrecv(c, elem, false)
return
}
``` 

从代码中可以看出，`select`中的`channel`相关操作，最终在编译阶段都会转换为`go/src/runtime/chan.go`的发送和接收操作。
- selectnbsend 
- selectnbrecv 、 selectnbrecv2 

`select`相当于`channel`收发操作的一层封装。  








