# 接口    
[code](../../../code/base/object/interface/base-interface.go)  
```go
package main

import (
	"fmt"
	"unsafe"
)

type IMan interface {
	walk() int
}

type Man struct {
	name string
	age  int
}

func (man *Man) walk() int {
	fmt.Println("walk name:", man.name, ",age:", man.age)
	return man.age
}

func main() {
	var test_interface interface{}
	man := Man{name: "xiaoming", age: 18}
	var iman IMan

	iman = &man
	iman.walk()
	man.walk()

	fmt.Printf("iman  %T 占中的字节数是 %d \n", iman, unsafe.Sizeof(iman))
	fmt.Println(man, iman, test_interface)
}
```

打印输出
```
walk name: xiaoming ,age: 18
walk name: xiaoming ,age: 18
iman  *main.Man 占中的字节数是 16 
{xiaoming 18} &{xiaoming 18} <nil>
```

## 数据结构  
`runtime/runtime2.go`定义`iface`
```go
type iface struct {
	tab  *itab
	data unsafe.Pointer
}

// layout of Itab known to compilers
// allocated in non-garbage-collected memory
// Needs to be in sync with
// ../cmd/compile/internal/gc/reflect.go:/^func.dumptabs.
type itab struct {
inter *interfacetype
_type *_type
hash  uint32 // copy of _type.hash. Used for type switches.
_     [4]byte
fun   [1]uintptr // variable sized. fun[0]==0 means _type does not implement inter.
}
```  

`interfacetype`定义在文件`runtime/type.go`  
```go
type interfacetype struct {
	typ     _type
	pkgpath name
	mhdr    []imethod   // method handler 
```

`接口`数据结构`data unsafe.Pointer`中会存储实现者的`指针`。  
通过`dlv`查看内存结构可以证明这一点:  
```shell
(dlv) print &man
(*main.Man)(0xc00000c030)
(dlv) print &iman
(*main.IMan)(0xc000065ee0)
(dlv) x -fmt hex -count 32 -size 1 0xc000065ee0
0xc000065ee0:   0x78   0xf6   0x0f   0x01   0x00   0x00   0x00   0x00   
0xc000065ee8:   0x30   0xc0   0x00   0x00   0xc0   0x00   0x00   0x00   #(*main.Man)(0xc00000c030)
0xc000065ef0:   0x08   0x01   0x40   0x01   0x00   0x00   0x00   0x00   
0xc000065ef8:   0xb8   0x02   0x00   0x00   0xc0   0x00   0x00   0x00
``` 
内存图`iman`的`data`指向实现的结构体`man`  
```
┌──────────────────────────┐                                                                                                                                                                                                          
│     interface(iman)      │
├─────────────┬────────────┤
│    tab      │    data    │ 
└──────┬──────┴─────┬──────┘
       │            │         ┌─────────────┐         ┌─────────────┐
       ↓            └────────→│ 0xc00000c030├────────→│ struct(man) │
┌─────────────┐               └─────────────┘         └─────────────┘
│  0x10ff678  │               
└─────────────┘
```  












