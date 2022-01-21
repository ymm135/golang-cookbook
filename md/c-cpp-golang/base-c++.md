# c++基础
[学习课程](https://coding.imooc.com/learn/list/414.html)  
## 基础语法
### c++概况
- 大型桌面应用 PS/Chrome/Microsoft Office
- 大型网站后台 搜索引擎
- 大型游戏后台 王章荣耀
- 大型游戏引擎 Unity
- 编译器/解释器 Java虚拟机/JS引擎  

2011年中期C++标准(C++11)完成新的标准,Boost库项目对新标准
产生了相当大的影响。  

### 编译型语言
源程序 => 编译器 => 目标程序 => 链接器 => 可执行程序

### 数据类型

| 类型  | 位   | 范围  |
|-----|-----|-----|
|  char   |  1 个字节   |  -128 到 127 或者 0 到 255   |
|  unsigned char   |  1 个字节   |  0 到 255   |
|  signed char   |  1 个字节   |  -128 到 127   |
|   int  |   4 个字节  |  -2147483648 到 2147483647   |
|  unsigned int   |  4 个字节   |  0 到 4294967295   |
|  signed int   |  4 个字节   |  -2147483648 到 2147483647   |
|   short int  |  	2 个字节   |  -32768 到 32767   |
|  long int   |  8 个字节   |  	-9,223,372,036,854,775,808 到 9,223,372,036,854,775,807   |
|   float  |   4 个字节  |  精度型占4个字节（32位）内存空间，+/- 3.4e +/- 38 (~7 个数字)   |
|   double  |  8 个字节   |   双精度型占8 个字节（64位）内存空间，+/- 1.7e +/- 308 (~15 个数字)  |
|  long double   |  16 个字节   |   长双精度型 16 个字节（128位）内存空间，可提供18-19位有效数字。  |
|  wchar_t   |  	2 或 4 个字节   |   个宽字符   |


### 常量与变量  
- define 编译时替换  
- const 编译时检测，推荐使用该方法    


```c++
#include <iostream>

#define NUMBER2 2

int main() {
    const int NUMBER1 = 1;
    int a = NUMBER2;
    std::cout << a << std::endl;
    return 0;
}
```
通过clion内存视图看到`NUMBER1`的内存为`01 00 00 00`，占用4个字节的变量。通过gdb窗口`disass /m`查看汇编  
> 通过`~/.gdbinit`设置汇编格式为inter，设置参数为:`set disassembly-flavor intel`  

```sh
6	    const int NUMBER1 = 1;
=> 0x00005596401801b5 <+12>:	mov    DWORD PTR [rbp-0x8],0x1   #汇编可以看出常量与变量没有区别,都仅仅是一个地址

7	    int a = NUMBER2;  
   0x00005596401801bc <+19>:	mov    DWORD PTR [rbp-0x4],0x2  # define宏定义在编译时已经被替换为2  
```


## 运算符与表达式

### 算数运算  
| 运算符 | 描述               | 实例             |
|-----|------------------|----------------|
| +	| 把两个操作数相加         | A + B 将得到 30   |
| -	| 从第一个操作数中减去第二个操作数 | A - B 将得到 -10  |
| *	| 把两个操作数相乘         | 	A * B 将得到 200 |
 |   /	| 分子除以分母           | B / A 将得到 2    |  
 |   %	| 取模运算符，整除后的余数	    | B % A 将得到 0    |
 |   ++	| 自增运算符，整数值增加 1	   | A++ 将得到 11     |
 |   --	| 自减运算符，整数值减少 1	   | A-- 将得到 9      |  

实例  
```c++
#include <iostream>

int main()
{
    int a = 10, b= 20;
    int r1 = a + b;
    int r2 = a - b;
    int r3 = a * b;   
    int r4 = a / b;  // 从汇编中可以看出eax存储整除结果,edx存储余数
    int r5 = a % b;  // 获取edx的值
    int r6 = a++;    // 先加载a的值到eax，再把eax的值给r6. 最终把a的值+1 
    int r7 = ++a;    // 先把a的值加1，再把a的值加载到eax中赋给r7  
    int r8 = a--;
    int r9 = --a;

    return 0;
}
```

汇编指令
```c++
Dump of assembler code for function main():
4	{
   0x000055d35ad4b25f <+0>:	endbr64 
   0x000055d35ad4b263 <+4>:	push   rbp
   0x000055d35ad4b264 <+5>:	mov    rbp,rsp                        //函数栈操作

5	    int a = 10, b= 20;
   0x000055d35ad4b267 <+8>:	mov    DWORD PTR [rbp-0x2c],0xa      //a变量地址 [rbp-0x2c]
   0x000055d35ad4b26e <+15>:	mov    DWORD PTR [rbp-0x28],0x14 //b变量地址 [rbp-0x28]

6	    int r1 = a + b;
   0x000055d35ad4b275 <+22>:	mov    edx,DWORD PTR [rbp-0x2c]
   0x000055d35ad4b278 <+25>:	mov    eax,DWORD PTR [rbp-0x28]
   0x000055d35ad4b27b <+28>:	add    eax,edx                  // a + b
   0x000055d35ad4b27d <+30>:	mov    DWORD PTR [rbp-0x24],eax // 结果赋值给r1 [rbp-0x24]

7	    int r2 = a - b;
   0x000055d35ad4b280 <+33>:	mov    eax,DWORD PTR [rbp-0x2c]
   0x000055d35ad4b283 <+36>:	sub    eax,DWORD PTR [rbp-0x28] // sub 减法
   0x000055d35ad4b286 <+39>:	mov    DWORD PTR [rbp-0x20],eax

8	    int r3 = a * b;
   0x000055d35ad4b289 <+42>:	mov    eax,DWORD PTR [rbp-0x2c]
   0x000055d35ad4b28c <+45>:	imul   eax,DWORD PTR [rbp-0x28] // imul 乘法 
   0x000055d35ad4b290 <+49>:	mov    DWORD PTR [rbp-0x1c],eax

9	    int r4 = a / b;   // divisor 除数 signed divide (idic) 有符号的整除  
   0x000055d35ad4b293 <+52>:	mov    eax,DWORD PTR [rbp-0x2c] // 加载变量a->eax=0xa; rdz=0xa
   0x000055d35ad4b296 <+55>:	cdq                             // edx=0x0 : eax=0xa  
   0x000055d35ad4b297 <+56>:	idiv   DWORD PTR [rbp-0x28]     // edx=0xa(余数) : eax=0x0(除数)
   0x000055d35ad4b29a <+59>:	mov    DWORD PTR [rbp-0x18],eax // 整除结果在eax中 

10	    int r5 = a % b;
   0x000055d35ad4b29d <+62>:	mov    eax,DWORD PTR [rbp-0x2c]
   0x000055d35ad4b2a0 <+65>:	cdq    
   0x000055d35ad4b2a1 <+66>:	idiv   DWORD PTR [rbp-0x28]
   0x000055d35ad4b2a4 <+69>:	mov    DWORD PTR [rbp-0x14],edx // 整除结果在edx中 

11	    int r6 = a++;  // 先加载a的值，再把a的值给r6. 最终把a的值+1
   0x000055d35ad4b2a7 <+72>:	mov    eax,DWORD PTR [rbp-0x2c] // 加载变量a->eax
   0x000055d35ad4b2aa <+75>:	lea    edx,[rax+0x1]            // 0xa + 0x1 = 0xb
   0x000055d35ad4b2ad <+78>:	mov    DWORD PTR [rbp-0x2c],edx // a = 0xb
   0x000055d35ad4b2b0 <+81>:	mov    DWORD PTR [rbp-0x10],eax // r6 = 0xa

12	    int r7 = ++a;  // 先把a的值加1，再把a的值赋给r7 
   0x000055d35ad4b2b3 <+84>:	add    DWORD PTR [rbp-0x2c],0x1 // 先把a的值加1 
   0x000055d35ad4b2b7 <+88>:	mov    eax,DWORD PTR [rbp-0x2c]
   0x000055d35ad4b2ba <+91>:	mov    DWORD PTR [rbp-0xc],eax

13	    int r8 = a--;
   0x000055d35ad4b2bd <+94>:	mov    eax,DWORD PTR [rbp-0x2c] // 加载a的值在eax中
   0x000055d35ad4b2c0 <+97>:	lea    edx,[rax-0x1]            // a的值减一存储到edx
   0x000055d35ad4b2c3 <+100>:	mov    DWORD PTR [rbp-0x2c],edx // 把edx的值赋给a变量
   0x000055d35ad4b2c6 <+103>:	mov    DWORD PTR [rbp-0x8],eax // 把eax的值赋给r8 

14	    int r9 = --a;
=> 0x000055d35ad4b2c9 <+106>:	sub    DWORD PTR [rbp-0x2c],0x1
   0x000055d35ad4b2cd <+110>:	mov    eax,DWORD PTR [rbp-0x2c]
   0x000055d35ad4b2d0 <+113>:	mov    DWORD PTR [rbp-0x4],eax

15	
16	    return 0;
   0x000055d35ad4b2d3 <+116>:	mov    eax,0x0

17	}
   0x000055d35ad4b2d8 <+121>:	pop    rbp
   0x000055d35ad4b2d9 <+122>:	ret  
```

> 正如广泛宣传的那样，现代的x86_64处理器具有64位寄存器，可以以向后兼容的方式用作32位寄存器，16位寄存器甚至8位寄存器，例如：  
```
0x1122334455667788
  ================ rax (64 bits)
          ======== eax (32 bits)
              ====  ax (16 bits)
              ==    ah (8 bits)
                ==  al (8 bits)
```

> 可以从字面上采用这样的方案，即，为了读取或写入的目的，总是可以使用指定的名称始终仅访问寄存器的一部分，这将是高度逻辑的。实际上，对于所有32位以下的系统都是这样：  
```
mov  eax, 0x11112222 ; eax = 0x11112222
mov  ax, 0x3333      ; eax = 0x11113333 (works, only low 16 bits changed)
mov  al, 0x44        ; eax = 0x11113344 (works, only low 8 bits changed)
mov  ah, 0x55        ; eax = 0x11115544 (works, only high 8 bits changed)
xor  ah, ah          ; eax = 0x11110044 (works, only high 8 bits cleared)
mov  eax, 0x11112222 ; eax = 0x11112222
xor  al, al          ; eax = 0x11112200 (works, only low 8 bits cleared)
mov  eax, 0x11112222 ; eax = 0x11112222
xor  ax, ax          ; eax = 0x11110000 (works, only low 16 bits cleared)
```

> [CBW, CWD, CDQ 介绍](http://www.c-jump.com/CIS77/MLabs/M11arithmetic/M11_0110_cbw_cwd_cdq.htm)  

### 关系运算符

| 运算符	 | 描述                              | 	实例 |    
|------|---------------------------------|-----|  
| ==	  | 检查两个操作数的值是否相等，如果相等则条件为真。        | 	(A == B) 不为真。 |
| !=	  | 检查两个操作数的值是否相等，如果不相等则条件为真。       |	(A != B) 为真。 |
| \>   | 检查左操作数的值是否大于右操作数的值，如果是则条件为真。    |	(A > B) 不为真。 |
| <	   | 检查左操作数的值是否小于右操作数的值，如果是则条件为真。    |	(A < B) 为真。 |
| >=	  | 检查左操作数的值是否大于或等于右操作数的值，如果是则条件为真。 |	(A >= B) 不为真。 |
| <=	  | 检查左操作数的值是否小于或等于右操作数的值，如果是则条件为真。 |	(A <= B) 为真。 |


实例  
```c++
#include <iostream>
using namespace std;

int main()
{
    int a = 21;
    int b = 10;
    int c ;

    if( a == b )
    {
        cout << "Line 1 - a 等于 b" << endl ;
    }

    if ( a < b )
    {
        cout << "Line 2 - a 小于 b" << endl ;
    }

    return 0;
}
```

汇编指令:
```c++
10	    if( a == b )
=> 0x000055b442ab4299 <+26>:	mov    eax,DWORD PTR [rbp-0x8]
   0x000055b442ab429c <+29>:	cmp    eax,DWORD PTR [rbp-0x4]
   0x000055b442ab429f <+32>:	jne    0x55b442ab42cc <main()+77>  // jne

15	    if ( a < b )
   0x000055b442ab42cc <+77>:	mov    eax,DWORD PTR [rbp-0x8]
   0x000055b442ab42cf <+80>:	cmp    eax,DWORD PTR [rbp-0x4]
   0x000055b442ab42d2 <+83>:	jge    0x55b442ab42ff <main()+128>  // jge 
```

- cmp 进行比较两个操作数的大小
- test 两个操作数的按位AND运算，并根据结果设置标志寄存器  
- je 当ZF标志位为零,就跳转 jump When Equal  
- jge 大于或者等于 jump when greater or equal 
- jb 低于,即不高于且不等于则转移


### 逻辑运算符
| 运算符	         | 描述	                                                | 实例                          |  
|--------------|----------------------------------------------------|-----------------------------|  
| &&	          | 称为逻辑与运算符。如果两个操作数都 true，则条件为 true。                  | 	(A && B) 为 false。          |  
| &#124;&#124; |	称为逻辑或运算符。如果两个操作数中有任意一个 true，则条件为 true。| 	(A &#124;&#124; B) 为 true。 |  
| !	           | 称为逻辑非运算符。用来逆转操作数的逻辑状态，如果条件为 true 则逻辑非运算符将使其为 false。 | 	!(A && B) 为 true。          |  


实例
```c++
int main()
{
    int a = 21;
    int b = 10;
    int c ;

    c = a && b;
    c = a || b;
    c = !a;

    return 0;
}
```

汇编实现
```c++
6	    int a = 21;
   0x00005642644cd267 <+8>:	mov    DWORD PTR [rbp-0xc],0x15

7	    int b = 10;
   0x00005642644cd26e <+15>:	mov    DWORD PTR [rbp-0x8],0xa
   
10	    c = a && b; 
=> 0x00005642644cd275 <+22>:	cmp    DWORD PTR [rbp-0xc],0x0     // 把a与0进行比较 
   0x00005642644cd279 <+26>:	je     0x5642644cd288 <main()+41>  // 如果a等于0, 跳转到<main()+41>
   0x00005642644cd27b <+28>:	cmp    DWORD PTR [rbp-0x8],0x0     // 把b与0进行比较
   0x00005642644cd27f <+32>:	je     0x5642644cd288 <main()+41>  
   0x00005642644cd281 <+34>:	mov    eax,0x1
   0x00005642644cd286 <+39>:	jmp    0x5642644cd28d <main()+46>
   0x00005642644cd288 <+41>:	mov    eax,0x0
   0x00005642644cd28d <+46>:	movzx  eax,al
   0x00005642644cd290 <+49>:	mov    DWORD PTR [rbp-0x4],eax

11	    c = a || b;
   0x00005642644cd293 <+52>:	cmp    DWORD PTR [rbp-0xc],0x0
   0x00005642644cd297 <+56>:	jne    0x5642644cd29f <main()+64>
   0x00005642644cd299 <+58>:	cmp    DWORD PTR [rbp-0x8],0x0
   0x00005642644cd29d <+62>:	je     0x5642644cd2a6 <main()+71>
   0x00005642644cd29f <+64>:	mov    eax,0x1
   0x00005642644cd2a4 <+69>:	jmp    0x5642644cd2ab <main()+76>
   0x00005642644cd2a6 <+71>:	mov    eax,0x0
   0x00005642644cd2ab <+76>:	movzx  eax,al
   0x00005642644cd2ae <+79>:	mov    DWORD PTR [rbp-0x4],eax

12	    c = !a;
   0x00005642644cd2b1 <+82>:	cmp    DWORD PTR [rbp-0xc],0x0
   0x00005642644cd2b5 <+86>:	sete   al
   0x00005642644cd2b8 <+89>:	movzx  eax,al
   0x00005642644cd2bb <+92>:	mov    DWORD PTR [rbp-0x4],eax
```


## 容器
## 指针
## 基础句法
## 高级语法
## 编程思想
## 进阶编程
## GUI开发
## 陷阱与经验