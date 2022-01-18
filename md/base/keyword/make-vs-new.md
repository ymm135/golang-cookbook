# make vs new 
- make 的作用是初始化内置的数据结构，也就是我们在前面提到的切片、哈希表和 Channel2；
- new 的作用是根据传入的类型分配一片内存空间并返回指向这片内存空间的指针3；

示例  
```go
     5:	func main() {
     6:		// make
     7:		slice := make([]int, 0, 100)      //调用runtime.makeslice
     8:		hash := make(map[int]bool, 10)    //调用runtime.makemap
=>   9:		ch := make(chan int, 5)           //调用runtime.makechan
    10:	
    11:		// new
    12:		i := new(int)        //调用runtime.newobject
    13:		var v int            //调用runtime.newobject
    14:		i = &v
}
```  

对应的汇编程序
```
(dlv) disass
TEXT main.main(SB) /Users/ymm/work/mygithub/golang-cookbook/code/base/keyword/make-new/base-make-new.go
	base-make-new.go:5	0x10cbaa0	65488b0c2530000000		mov rcx, qword ptr gs:[0x30]
	base-make-new.go:5	0x10cbaa9	488d442490			lea rax, ptr [rsp-0x70]
	base-make-new.go:5	0x10cbaae	483b4110			cmp rax, qword ptr [rcx+0x10]
	base-make-new.go:5	0x10cbab2	0f869c020000			jbe 0x10cbd54
	base-make-new.go:5	0x10cbab8*	4881ecf0000000			sub rsp, 0xf0
	base-make-new.go:5	0x10cbabf	4889ac24e8000000		mov qword ptr [rsp+0xe8], rbp
	base-make-new.go:5	0x10cbac7	488dac24e8000000		lea rbp, ptr [rsp+0xe8]
	base-make-new.go:7	0x10cbacf	488d056aae0000			lea rax, ptr [rip+0xae6a]
	base-make-new.go:7	0x10cbad6	48890424			mov qword ptr [rsp], rax
	base-make-new.go:7	0x10cbada	48c744240800000000		mov qword ptr [rsp+0x8], 0x0
	base-make-new.go:7	0x10cbae3	48c744241064000000		mov qword ptr [rsp+0x10], 0x64
	base-make-new.go:7	0x10cbaec	e8af50f8ff			call $runtime.makeslice
	base-make-new.go:7	0x10cbaf1	488b442418			mov rax, qword ptr [rsp+0x18]
	base-make-new.go:7	0x10cbaf6	4889442478			mov qword ptr [rsp+0x78], rax
	base-make-new.go:7	0x10cbafb	48c784248000000000000000	mov qword ptr [rsp+0x80], 0x0
	base-make-new.go:7	0x10cbb07	48c784248800000064000000	mov qword ptr [rsp+0x88], 0x64
	base-make-new.go:8	0x10cbb13	488d0526ef0000			lea rax, ptr [rip+0xef26]
	base-make-new.go:8	0x10cbb1a	48890424			mov qword ptr [rsp], rax
	base-make-new.go:8	0x10cbb1e	48c74424080a000000		mov qword ptr [rsp+0x8], 0xa
	base-make-new.go:8	0x10cbb27	48c744241000000000		mov qword ptr [rsp+0x10], 0x0
	base-make-new.go:8	0x10cbb30	e8eb30f4ff			call $runtime.makemap
	base-make-new.go:8	0x10cbb35	488b442418			mov rax, qword ptr [rsp+0x18]
	base-make-new.go:8	0x10cbb3a	4889442438			mov qword ptr [rsp+0x38], rax
=>	base-make-new.go:9	0x10cbb3f	488d05faa70000			lea rax, ptr [rip+0xa7fa]
	base-make-new.go:9	0x10cbb46	48890424			mov qword ptr [rsp], rax
	base-make-new.go:9	0x10cbb4a	48c744240805000000		mov qword ptr [rsp+0x8], 0x5
	base-make-new.go:9	0x10cbb53	e84890f3ff			call $runtime.makechan
	base-make-new.go:9	0x10cbb58	488b442410			mov rax, qword ptr [rsp+0x10]
	base-make-new.go:9	0x10cbb5d	4889442440			mov qword ptr [rsp+0x40], rax
	base-make-new.go:12	0x10cbb62	488d05d7ad0000			lea rax, ptr [rip+0xadd7]
	base-make-new.go:12	0x10cbb69	48890424			mov qword ptr [rsp], rax
	base-make-new.go:12	0x10cbb6d	e8ee21f4ff			call $runtime.newobject
	base-make-new.go:12	0x10cbb72	488b442408			mov rax, qword ptr [rsp+0x8]
	base-make-new.go:12	0x10cbb77	4889442430			mov qword ptr [rsp+0x30], rax
	base-make-new.go:13	0x10cbb7c	488d05bdad0000			lea rax, ptr [rip+0xadbd]
	base-make-new.go:13	0x10cbb83	48890424			mov qword ptr [rsp], rax
	base-make-new.go:13	0x10cbb87	e8d421f4ff			call $runtime.newobject
	base-make-new.go:13	0x10cbb8c	488b442408			mov rax, qword ptr [rsp+0x8]
	base-make-new.go:13	0x10cbb91	4889442470			mov qword ptr [rsp+0x70], rax
	base-make-new.go:13	0x10cbb96	48c70000000000			mov qword ptr [rax], 0x0
	base-make-new.go:14	0x10cbb9d	488b442470			mov rax, qword ptr [rsp+0x70]
	base-make-new.go:14	0x10cbba2	4889442430			mov qword ptr [rsp+0x30], rax
```

从汇编实现中可以看出`new(int)`和`var v int`都是调用`newobject`, 最终返回的是内存地址  
```go
// implementation of new builtin
// compiler (both frontend and SSA backend) knows the signature
// of this function
func newobject(typ *_type) unsafe.Pointer {
	return mallocgc(typ.size, typ, true)
}
```


