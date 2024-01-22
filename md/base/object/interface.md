- # 接口    

目录:  

- [实例](#实例)
- [数据结构](#数据结构)


## 实例
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

也可以从汇编实现看出结构体赋值给指针`iman = &man`
```go
    22:	func main() {
    23:		var test_interface interface{}
    24:		man := Man{name: "xiaoming", age: 18}
    25:		var iman IMan
    26:
=>  27:		iman = &man
    28:		iman.walk()
    29:		man.walk()
```
汇编实现如下:
```shell
	base-interface.go:24	0x10cdd7a	488d05df370100			lea rax, ptr [rip+0x137df]
	base-interface.go:24	0x10cdd81	48890424			mov qword ptr [rsp], rax
	base-interface.go:24	0x10cdd85	e8d6fff3ff			call $runtime.newobject
	base-interface.go:24	0x10cdd8a	488b442408			mov rax, qword ptr [rsp+0x8]
	base-interface.go:24	0x10cdd8f	4889442470			mov qword ptr [rsp+0x70], rax            #[rsp+0x70] 就是变量man的地址
	base-interface.go:24	0x10cdd94	0f57c0				xorps xmm0, xmm0
	base-interface.go:24	0x10cdd97	0f118424b8000000		movups xmmword ptr [rsp+0xb8], xmm0
	base-interface.go:24	0x10cdd9f	48c78424c800000000000000	mov qword ptr [rsp+0xc8], 0x0
	base-interface.go:24	0x10cddab	488d05385d0200			lea rax, ptr [rip+0x25d38]
	base-interface.go:24	0x10cddb2	48898424b8000000		mov qword ptr [rsp+0xb8], rax
	base-interface.go:24	0x10cddba	48c78424c000000008000000	mov qword ptr [rsp+0xc0], 0x8
	base-interface.go:24	0x10cddc6	48c78424c800000012000000	mov qword ptr [rsp+0xc8], 0x12
	base-interface.go:24	0x10cddd2	488b7c2470			mov rdi, qword ptr [rsp+0x70]
	base-interface.go:24	0x10cddd7	48c7470808000000		mov qword ptr [rdi+0x8], 0x8
	base-interface.go:24	0x10cdddf	48c7471012000000		mov qword ptr [rdi+0x10], 0x12
	base-interface.go:24	0x10cdde7	833d02760d0000			cmp dword ptr [runtime.writeBarrier], 0x0
	base-interface.go:24	0x10cddee	7405				jz 0x10cddf5
	base-interface.go:24	0x10cddf0	e946030000			jmp 0x10ce13b
	base-interface.go:24	0x10cddf5	488907				mov qword ptr [rdi], rax
	base-interface.go:24	0x10cddf8	eb00				jmp 0x10cddfa                   
	base-interface.go:25	0x10cddfa	0f57c0				xorps xmm0, xmm0
	base-interface.go:25	0x10cddfd	0f11842488000000		movups xmmword ptr [rsp+0x88], xmm0		#[rsp+0x88]是iman变量的地址
=>	base-interface.go:27	0x10cde05*	488b442470			mov rax, qword ptr [rsp+0x70]			#[rsp+0x70]存储的就是man变量的地址
	base-interface.go:27	0x10cde0a	4889442448			mov qword ptr [rsp+0x48], rax
	base-interface.go:27	0x10cde0f	488d0d823c0300			lea rcx, ptr [rip+0x33c82]
	base-interface.go:27	0x10cde16	48898c2488000000		mov qword ptr [rsp+0x88], rcx            #把`tab  *itab`赋值给iman变量
	base-interface.go:27	0x10cde1e	4889842490000000		mov qword ptr [rsp+0x90], rax            #就是把man变量的地址赋给iman.data变量，也就是指向结构体的实现  
	base-interface.go:28	0x10cde26	488b842488000000		mov rax, qword ptr [rsp+0x88]
```

从汇编实现可以看出接口赋值需要填充两部分，一部分是`tab  *itab`，另一部分是结构体实现`data unsafe.Pointer`部分。  














