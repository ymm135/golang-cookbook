- # defer  

目录:  
- [执行顺序](#执行顺序)
- [预计算参数](#预计算参数)
- [数据结构及实现原理](#数据结构及实现原理)


## 执行顺序  
先声明后执行，相当于堆栈，`先进后出`  
```go
func main() {
	defer println("A")
	defer println("B")
	defer println("C")
}
```  

输出结果:  
```shell
C
B
A
```  

## 预计算参数  
[参考文章](https://draveness.me/golang/docs/part2-foundation/ch05-keyword/golang-defer/)  

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	testParam1()
	testParam2()
}

func testParam1() {
	startedAt := time.Now()
	defer fmt.Println("testParam1=", time.Since(startedAt))

	time.Sleep(time.Second)
}

func testParam2() {
	startedAt := time.Now()
	defer func() { fmt.Println("testParam1=", time.Since(startedAt)) }()

	time.Sleep(time.Second)
}
```

输出结果是:
```
testParam1= 185ns
testParam1= 1.004880495s
```

从结果可以看出`defer`后会立刻拷贝函数中引用的外部参数，所以 `time.Since(startedAt)` 的结果不是在 `main` 函数退出之前计算的  
`defer`关键字后增加匿名函数,虽然调用 defer 关键字时也使用`值传递`，但是因为拷贝的是函数指针，所以 `time.Since(startedAt)` 
会在 `main` 函数返回前调用并打印出符合预期的结果。  

> 匿名函数也是函数，函数执行是只要找到函数指针指向的`代码地址`(引用传递)即可，不会像`值类型`那样固定。  


## 数据结构及实现原理  

```go
// A _defer holds an entry on the list of deferred calls.
// If you add a field here, add code to clear it in freedefer and deferProcStack
// This struct must match the code in cmd/compile/internal/gc/reflect.go:deferstruct
// and cmd/compile/internal/gc/ssa.go:(*state).call.
// Some defers will be allocated on the stack and some on the heap.
// All defers are logically part of the stack, so write barriers to
// initialize them are not required. All defers must be manually scanned,
// and for heap defers, marked.
type _defer struct {
siz     int32 // includes both arguments and results
started bool
heap    bool
// openDefer indicates that this _defer is for a frame with open-coded
// defers. We have only one defer record for the entire frame (which may
// currently have 0, 1, or more defers active).
openDefer bool
sp        uintptr  // sp at time of defer
pc        uintptr  // pc at time of defer
fn        *funcval // can be nil for open-coded defers
_panic    *_panic  // panic that is running defer
link      *_defer

// If openDefer is true, the fields below record values about the stack
// frame and associated function that has the open-coded defer(s). sp
// above will be the sp for the frame, and pc will be address of the
// deferreturn call in the function.
fd   unsafe.Pointer // funcdata for the function associated with the frame
varp uintptr        // value of varp for the stack frame
// framepc is the current pc associated with the stack frame. Together,
// with sp above (which is the sp associated with the stack frame),
// framepc/sp can be used as pc/sp pair to continue a stack trace via
// gentraceback().
framepc uintptr
}
```  

defer数据结构
```shell
┌──────────────────────────┐       ┌──────────────────────────┐       ┌──────────────────────────┐       ┌──────────────────────────┐                                                                                                                                                                                                             
│        goroutine         │       │          defer           │       │          defer           │       │          defer           │
├─────────────┬────────────┤       ├─────────────┬────────────┤       ├─────────────┬────────────┤       ├─────────────┬────────────┤
│    ...      │_defer->link│       │    ...      │    link    │       │    ...      │    link    │       │    ...      │    link    │
└─────────────┴─────┬──────┘       └─────────────┴─────┬──────┘       └─────────────┴──────┬─────┘       └─────────────┴────────────┘  
                    │                            ↑     │                            ↑      │                           ↑
                    └────────────────────────────┘     └────────────────────────────┘      └───────────────────────────┘
```

后调用的 defer 函数会先执行：
- 后调用的 defer 函数会被追加到 Goroutine _defer 链表的最前面；
- 运行 runtime._defer 时是从前到后依次执行；



