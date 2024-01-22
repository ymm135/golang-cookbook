- # panic and recover  

目录:  
- [样例](#样例)
	- [panic](#panic)
	- [recover](#recover)
- [数据结构](#数据结构)


两个关键字的作用:  
- `panic` 能够改变程序的控制流，调用 `panic` 后会立刻停止执行当前函数的剩余代码，并在当前 Goroutine 中递归执行调用方的 defer；
- `recover` 可以中止 `panic` 造成的程序崩溃(`recover` 只有在发生 `panic` 之后调用才会生效)。它是一个只能在 `defer` 中发挥作用的函数，在其他作用域中调用不会发挥作用；

## 样例
### panic 
```go
func main() {
	defer fmt.Println("in main")
	defer func() {
		defer func() {
			fmt.Println("panic again and again")
			panic("panic again and again")
		}()
		fmt.Println("panic again")
		panic("panic again")
	}()

	panic("panic once")
}
``` 
`panic`嵌套的情况:  
打印输出
```
panic again
panic again and again
in main  
panic: panic once
        panic: panic again
        panic: panic again and again
```
从上述程序输出的结果，我们可以确定程序多次调用 `panic` 也不会影响 `defer` 函数的正常执行，
所以使用 `defer` 进行收尾工作一般来说都是安全的。

### recover 

```go
func main() {
	defer fmt.Println("in main")
	if err := recover(); err != nil {
		fmt.Println(err)
	}
	
	panic("unknown err")
}
```
`recover`没有在`panic`之后调用，不生效,输出
```
in main
panic: unknown err

goroutine 1 [running]:
main.main()
    /Users/xxx/work/mygithub/golang-cookbook/code/base/keyword/panic-recover/recover/base-recover2.go:11 +0x125
```

```go 
func main() {
	defer fmt.Println("in main")
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	panic("unknown err")
}
``` 
recover生效，打印输出是
```shell
unknown err
in main
```

## 数据结构 

`go/src/runtime/runtime2.go`文件中
```go
// A _panic holds information about an active panic.
//
// A _panic value must only ever live on the stack.
//
// The argp and link fields are stack pointers, but don't need special
// handling during stack growth: because they are pointer-typed and
// _panic values only live on the stack, regular stack pointer
// adjustment takes care of them.
type _panic struct {
	argp      unsafe.Pointer // pointer to arguments of deferred call run during panic; cannot move - known to liblink
	arg       interface{}    // argument to panic
	link      *_panic        // link to earlier panic 可以有多个
	pc        uintptr        // where to return to in runtime if this panic is bypassed
	sp        unsafe.Pointer // where to return to in runtime if this panic is bypassed
	recovered bool           // whether this panic is over
	aborted   bool           // the panic was aborted
	goexit    bool
}
```  


