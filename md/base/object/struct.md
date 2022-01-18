# 结构体  
[code](code/base/object/struct/base-struct.go)   
```go
package main

import "fmt"

type Man struct {
	name string
	age  int
}

func (man *Man) Walk() int {
	fmt.Println("man Walk")
	return 0
}

func main() {
	man := Man{name: "xiaoming", age: 18}

	var walkFun func() int
	walkFun = man.Walk
	walkFun()
}
```

## 数据结构

在go源码中暂时没有搜索到`struct`的定义，先使用`dlv`debug看看  

```shell
dlv debug main.go 

(dlv) print man
main.Man {
	name: "xiaoming",
	age: 18,}
(dlv) print &man
(*main.Man)(0xc0000b6018)
(dlv) print &man.name
(*string)(0xc0000b6018)     #从这看出man与man.name地址相同,结构体开始就是"string"
(dlv) print &man.age
(*int)(0xc0000b6028)        #偏移16个字节
```

查看结构体内存视图:  
```shell
(dlv) x -fmt hex -count 32 -size 1 0xc0000b6018    # 查看man结构体
0xc0000b6018:   0x8f   0x14   0x0f   0x01   0x00   0x00   0x00   0x00   
0xc0000b6020:   0x08   0x00   0x00   0x00   0x00   0x00   0x00   0x00   
0xc0000b6028:   0x12   0x00   0x00   0x00   0x00   0x00   0x00   0x00   
0xc0000b6030:   0x00   0x00   0x00   0x00   0x00   0x00   0x00   0x00

# 在查看name指向的字符串地址(0x010f148f) 
# xiaoming 对应16进制 78 69 61 6f 6d 69 6e 67
(dlv) x -fmt hex -count 32 -size 1 0x010f148f
0x10f148f:   0x78   0x69   0x61   0x6f   0x6d   0x69   0x6e   0x67   
0x10f1497:   0x20   0x28   0x66   0x6f   0x72   0x63   0x65   0x64   
0x10f149f:   0x29   0x20   0x2d   0x3e   0x20   0x6e   0x6f   0x64   
0x10f14a7:   0x65   0x3d   0x20   0x62   0x6c   0x6f   0x63   0x6b

```  

那么man结构体数据结构为(图形为特殊符号):
```shell
┌──────────────────────────┬─────────────┐                                                                                                                                                                                                           
│          string          │    int      │
├─────────────┬────────────┼─────────────┤
│    data     │     len    │             │
└─────────────┴────────────┴─────────────┘
``` 

这样看来结构体中的成员变量就是顺序排列的，如果是引用类型，那就是**引用类型的数据结构，不是一个指针**  

## 结构体如何调用方法  

```shell
(dlv) c
> main.main() ./base-struct.go:20 (hits goroutine(1):1 total:1) (PC: 0x10cbbea)
    15:	func main() {
    16:		man := Man{name: "xiaoming", age: 18}
    17:	
    18:		var walkFun func() int
    19:		walkFun = man.Walk
=>  20:		walkFun()
    21:	}
```  
把函数作为参数传递给`walkFun`，查看`walkFun`的内存视图: 

```shell
(dlv) print &man
(*main.Man)(0xc00009df60)
(dlv) print &walkFun
(*func() int)(0xc00009df40)
(dlv) print walkFun
main.(*Man).Walk-fm
```  

waklFun的地址为`0xc00009df40`,指向的地址为`0x0c0009df50`->`0x10cbc20`
```shell
(dlv) x -fmt hex -count 32 -size 1 0xc00009df40
0xc00009df40:   0x50   0xdf   0x09   0x00   0xc0   0x00   0x00   0x00   
0xc00009df48:   0x50   0xdf   0x09   0x00   0xc0   0x00   0x00   0x00   
0xc00009df50:   0x20   0xbc   0x0c   0x01   0x00   0x00   0x00   0x00   
0xc00009df58:   0x60   0xdf   0x09   0x00   0xc0   0x00   0x00   0x00

(dlv) x -fmt hex -count 32 -size 1 0xc00009df50
0xc00009df50:   0x20   0xbc   0x0c   0x01   0x00   0x00   0x00   0x00   
0xc00009df58:   0x60   0xdf   0x09   0x00   0xc0   0x00   0x00   0x00   
0xc00009df60:   0x0f   0x12   0x0f   0x01   0x00   0x00   0x00   0x00   
0xc00009df68:   0x08   0x00   0x00   0x00   0x00   0x00   0x00   0x00 

(dlv) x -fmt hex -count 32 -size 1 0x10cbc20
0x10cbc20:   0x65   0x48   0x8b   0x0c   0x25   0x30   0x00   0x00   
0x10cbc28:   0x00   0x48   0x3b   0x61   0x10   0x76   0x51   0x48   
0x10cbc30:   0x83   0xec   0x28   0x48   0x89   0x6c   0x24   0x20   
0x10cbc38:   0x48   0x8d   0x6c   0x24   0x20   0x48   0x8b   0x59 
```    

到这里并不能看出`方法`是什么？结构体`Man`并未拥有方法，但是可以调用它？？ 从这看出`0x10cbc20`
的内存并没有什么特殊，有没有可能是`代码地址`呢？  

在函数`func (man *Man) Walk() int `里增加个断点
```shell
(dlv) b base-struct.go:11
(dlv) c # 调到下个断点

(dlv) ls
> main.(*Man).Walk() ./base-struct.go:11 (hits goroutine(1):1 total:1) (PC: 0x10cbaca)
     6:		name string
     7:		age  int
     8:	}
     9:	
    10:	func (man *Man) Walk() int {
=>  11:		fmt.Println("man Walk")
    12:		return 0
    13:	}
    14:	
    15:	func main() {
    16:		man := Man{name: "xiaoming", age: 18}
    
```
查看汇编代码，确认代码块位置`base-struct.go:10`  
```shell
(dlv) disassemble -a 0x10cbc00  0x10cbc40
TEXT main.main(SB) /Users/ymm/work/mygithub/golang-cookbook/code/base/object/struct/base-struct.go
	base-struct.go:15	0x10cbc00	e81b13faff		call $runtime.morestack_noctxt
	.:0			0x10cbc05	e956ffffff		jmp $main.main
	.:0			0x10cbc0a	cc			int3
	.:0			0x10cbc0b	cc			int3
	.:0			0x10cbc0c	cc			int3
	.:0			0x10cbc0d	cc			int3
	.:0			0x10cbc0e	cc			int3
	.:0			0x10cbc0f	cc			int3
	.:0			0x10cbc10	cc			int3
	.:0			0x10cbc11	cc			int3
	.:0			0x10cbc12	cc			int3
	.:0			0x10cbc13	cc			int3
	.:0			0x10cbc14	cc			int3
	.:0			0x10cbc15	cc			int3
	.:0			0x10cbc16	cc			int3
	.:0			0x10cbc17	cc			int3
	.:0			0x10cbc18	cc			int3
	.:0			0x10cbc19	cc			int3
	.:0			0x10cbc1a	cc			int3
	.:0			0x10cbc1b	cc			int3
	.:0			0x10cbc1c	cc			int3
	.:0			0x10cbc1d	cc			int3
	.:0			0x10cbc1e	cc			int3
	.:0			0x10cbc1f	cc			int3
=>	base-struct.go:10	0x10cbc20	65488b0c2530000000	mov rcx, qword ptr gs:[0x30]
	base-struct.go:10	0x10cbc29	483b6110		cmp rsp, qword ptr [rcx+0x10]
	base-struct.go:10	0x10cbc2d	7651			jbe 0x10cbc80
	base-struct.go:10	0x10cbc2f	4883ec28		sub rsp, 0x28
	base-struct.go:10	0x10cbc33	48896c2420		mov qword ptr [rsp+0x20], rbp
	base-struct.go:10	0x10cbc38	488d6c2420		lea rbp, ptr [rsp+0x20]
	base-struct.go:10	0x10cbc3d	48			rex.w
	base-struct.go:10	0x10cbc3e	8b			prefix(0x8b)
	base-struct.go:10	0x10cbc3f	59			pop rcx
```  

这样就可以确定函数代表的是`代码地址`，也就是要执行和跳转的地址，这里`walkFun()`代表要执行
`base-struct.go:10	0x10cbc20`代码块的内容，也就是指定`func (man *Man) Walk() int`
这个`(man *Man)`会把`Man`的实现当做参数传入函数，相当于`func Walk(man *Man) int`  

> 这就理解`man.walk()`的含义了，相当于吧`man`作为`walk()int`的参数而已，并不向Java那样
> 对象是和方法绑定在一起的，这里是完全分开的，在运行时通过参数传入，效果是一样的。  

备注: 如果想要查看结构体调用方法是如何被当做参数传入的，可查看汇编。  

> 也可以参看 [go常用语句对应的汇编指令之函数调用](https://github.com/ymm135/go-build/blob/master/gouse-assembly.md)    

```
(dlv) disassemble
TEXT main.(*Man).Walk(SB) /Users/ymm/work/mygithub/golang-cookbook/code/base/object/struct/base-struct.go
	base-struct.go:10	0x10cbaa0	65488b0c2530000000	mov rcx, qword ptr gs:[0x30]
	base-struct.go:10	0x10cbaa9	483b6110		cmp rsp, qword ptr [rcx+0x10]
	base-struct.go:10	0x10cbaad	0f868d000000		jbe 0x10cbb40
	base-struct.go:10	0x10cbab3	4883ec68		sub rsp, 0x68
	base-struct.go:10	0x10cbab7	48896c2460		mov qword ptr [rsp+0x60], rbp
	base-struct.go:10	0x10cbabc	488d6c2460		lea rbp, ptr [rsp+0x60]
	base-struct.go:10	0x10cbac1	48c744247800000000	mov qword ptr [rsp+0x78], 0x0
=>	base-struct.go:11	0x10cbaca*	0f57c0			xorps xmm0, xmm0
	base-struct.go:11	0x10cbacd	0f11442438		movups xmmword ptr [rsp+0x38], xmm0
	base-struct.go:11	0x10cbad2	488d442438		lea rax, ptr [rsp+0x38]
	base-struct.go:11	0x10cbad7	4889442430		mov qword ptr [rsp+0x30], rax
	base-struct.go:11	0x10cbadc	8400			test byte ptr [rax], al
	base-struct.go:11	0x10cbade	488d0d5bb40000		lea rcx, ptr [rip+0xb45b]
	base-struct.go:11	0x10cbae5	48894c2438		mov qword ptr [rsp+0x38], rcx
	base-struct.go:11	0x10cbaea	488d0d17200300		lea rcx, ptr [rip+0x32017]
	base-struct.go:11	0x10cbaf1	48894c2440		mov qword ptr [rsp+0x40], rcx
```  

















