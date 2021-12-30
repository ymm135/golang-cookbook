# 数组/切片  

## [从汇编角度理解数组/切片](https://github.com/ymm135/TD4-4BIT-CPU/blob/master/go-asm.md#go%E6%B1%87%E7%BC%96%E6%8C%87%E4%BB%A4%E5%AD%A6%E4%B9%A0)  

## 数组及切片的数据结构  

数组的数据结构就是连续的内存块,如果是值类型，内存块存储的就是数据。  


`runtime/slice.go`文件中  
```go
type slice struct {
	array unsafe.Pointer
	len   int
	cap   int
}
``` 

从数据结构中可以看出切片包含一个指针，指向数组，另外包含length及cap两个计数器。  

## 内存中的数组/切片  

测试代码
```go
package main

import "fmt"

func main() {
	var a [5]int
	a[0] = 1
	a[3] = 10
	a[4] = 15
	aLen := len(a)

	s := make([]int, 5, 10)
	s[2] = 2
	s[3] = 256
	length := len(s)
	fmt.Println(a, "arrayLen=", aLen, s, "sliceLen=", length)
}
```

输出结果为:  
```text
[1 0 0 10 15] arrayLen= 5 [0 0 2 256 0] sliceLen= 5  
```

首先编译代码
```shell
go mod init gotest  
go mod tidy  

go build // 生成gotest二进制文件  
```

vscode调试配置
不能使用golang的配置
```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Launch Package",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${file}"
        }
    ]
}
```  

需要使用`gdb`调试配置
```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "(gdb) 启动",
            "type": "cppdbg",
            "request": "launch",
            "program": "${workspaceFolder}/gotest",
            "args": [],
            "stopAtEntry": false,
            "cwd": "${fileDirname}",
            "environment": [],
            "externalConsole": false,
            "MIMode": "gdb",
            "setupCommands": [
                {
                    "description": "为 gdb 启用整齐打印",
                    "text": "-enable-pretty-printing",
                    "ignoreFailures": true
                }
            ]
        }
    ]
}
```

数组的内存视图:
```shell
-exec x/24x &a
0xc000088ee0:	0x00000001	0x00000000	0x00000000	0x00000000
0xc000088ef0:	0x00000000	0x00000000	0x0000000a	0x00000000
0xc000088f00:	0x0000000f	0x00000000	0x00000060	0x00000000
0xc000088f10:	0x00082000	0x000000c0	0x0003a748	0x000000c0
0xc000088f20:	0x0054b520	0x00000000	0x00000000	0x00000000
0xc000088f30:	0x0003a778	0x000000c0	0x00088f78	0x000000c0
```
`-exec x/24x &a` 使用gdb查看内存视图，`24x`代表是显示24个16进制，每个`int`占用8字节，
5个`int`占用5个字节，从`0xc000088ee0`-`0xc000088f08`，存储的值依次是`[0x01 0x00 0x00 0x0a 0x0f]`  


类似于C语言的调试，通过gdb可以清晰看到slice的内部实现(数据结构),并且可以查看slice的内存视图

![slice的调试切片](../../../res/slice的调试切片.png)   

通过`-exec x/64x &slice` 可以看到slice的内存视图，总占用32个字节,和数据结构是对应的。 
