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
> 在macos系统的clion中，可以通过设置toolchains的Debugger，设置调试器为GDB而不是LLDB。   

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
   0x00005642644cd279 <+26>:	je     0x5642644cd288 <main()+41>  // 如果a等于0, 跳转到<main()+41> 结果为0
   0x00005642644cd27b <+28>:	cmp    DWORD PTR [rbp-0x8],0x0     // 把b与0进行比较
   0x00005642644cd27f <+32>:	je     0x5642644cd288 <main()+41>  // 如果b等于0, 跳转到<main()+41> 结果为0
   0x00005642644cd281 <+34>:	mov    eax,0x1                     // 如果a与b都不等于0，结果为1  
   0x00005642644cd286 <+39>:	jmp    0x5642644cd28d <main()+46>  // 结束了，结果为1
   0x00005642644cd288 <+41>:	mov    eax,0x0                     // 如果a、b有一个为0，就跳转到这，结果为0
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

从汇编实现中可以看出，两个逻辑运算数，都会检验各自是否为0(false)，如果满足条件跳转到对应位置即可。  

### 赋值运算符  
| 运算符 | 	描述 | 实例                        |  
|-------|---------------|---------------------------|
| =	| 简单的赋值运算符，把右边操作数的值赋给左边操作数 | C = A + B 将把 A + B 的值赋给 C |
| +=	 | 加且赋值运算符，把右边操作数加上左边操作数的结果赋值给左边操作数	| C += A 相当于 C = C + A      |
| -=	 | 减且赋值运算符，把左边操作数减去右边操作数的结果赋值给左边操作数	| C -= A 相当于 C = C - A      |
| *=	 | 乘且赋值运算符，把右边操作数乘以左边操作数的结果赋值给左边操作数	| C *= A 相当于 C = C * A      |
| /=	 | 除且赋值运算符，把左边操作数除以右边操作数的结果赋值给左边操作数	| C /= A 相当于 C = C / A      |
| %=	 | 求模且赋值运算符，求两个操作数的模赋值给左边操作数	C %= A | 相当于 C = C % A             |
| <<= | 	左移且赋值运算符	 | C <<= 2 等同于 C = C << 2    |
| &gt;&gt;= | 	右移且赋值运算符	 | C >>= 2 等同于 C = C >> 2    |
| &=	 | 按位与且赋值运算符	 | C &= 2 等同于 C = C & 2      |
| ^=	 | 按位异或且赋值运算符	 | C ^= 2 等同于 C = C ^ 2      |
| &#124;=	 | 按位或且赋值运算符	 | C &#124;= 2 等同于 C = C     | 2 |  

实例及汇编实现
```c++
6	    int a = 21;
   0x0000556c23ef6267 <+8>:	mov    DWORD PTR [rbp-0xc],0x15

7	    int b = 10;
   0x0000556c23ef626e <+15>:	mov    DWORD PTR [rbp-0x8],0xa

8	    int c ;
9	
10	    c += a;
=> 0x0000556c23ef6275 <+22>:	mov    eax,DWORD PTR [rbp-0xc]  // 加载a
   0x0000556c23ef6278 <+25>:	add    DWORD PTR [rbp-0x4],eax  // 计算c + a

11	    c &= a;
   0x0000556c23ef627b <+28>:	mov    eax,DWORD PTR [rbp-0xc]
   0x0000556c23ef627e <+31>:	and    DWORD PTR [rbp-0x4],eax
```

### 位运算符  
| 运算符 | 描述 | 实例 |
|----------|----------|--------|
| &        | 	按位与操作，按二进制位进行"与"运算。 | A & B) 将得到 12，即为 0000 1100 |
| &#124;   | 按位或运算符，按二进制位进行"或"运算。| (A | B) 将得到 61，即为 0011 1101 |
| ^	       | 异或运算符，按二进制位进行"异或"运算。 | (A ^ B) 将得到 49，即为 0011 0001 |
| ~	       | 取反运算符，按二进制位进行"取反"运算。 | (~A ) 将得到 -61，即为 1100 0011，一个有符号二进制数的补码形式。 |
| <<       | 二进制左移运算符。将一个运算对象的各二进制位全部左移若干位（左边的二进制位丢弃，右边补0）| A << 2 将得到 240，即为 1111 0000 |
| &gt;&gt; | 二进制右移运算符。将一个数的各二进制位全部右移若干位，正数左补0，负数左补1，右边丢弃。| A >> 2 将得到 15，即为 0000 1111 |


实例及汇编实现
```c++
6	    int a = 21;
   0x00005636333c6267 <+8>:	mov    DWORD PTR [rbp-0xc],0x15

7	    int b = 10;
   0x00005636333c626e <+15>:	mov    DWORD PTR [rbp-0x8],0xa

8	    int c ;
9	
10	    c = a >> 2;
=> 0x00005636333c6275 <+22>:	mov    eax,DWORD PTR [rbp-0xc]  
   0x00005636333c6278 <+25>:	sar    eax,0x2                  // SAR destination, count 右移操作符
   0x00005636333c627b <+28>:	mov    DWORD PTR [rbp-0x4],eax

11	    c = a & b;
   0x00005636333c627e <+31>:	mov    eax,DWORD PTR [rbp-0xc]
   0x00005636333c6281 <+34>:	and    eax,DWORD PTR [rbp-0x8]  // 与 
   0x00005636333c6284 <+37>:	mov    DWORD PTR [rbp-0x4],eax
```

### 算数优先级
- 一元运算符的优先级高于二级运算符。
- 弄不清优先级的，加括号。

### 补码  
> 有没有一种办法用加法来计算减法 ？ 


## 容器
### 概念  
- 代表内存里一组连续的同类型存储区  
- 可以用来把多个存储区合并成一个整体   

### 数组  

```c++
6	    int a[] = {1,2,3,4,5,6,7,8};
=> 0x000055f61bcd829a <+27>:	mov    DWORD PTR [rbp-0x30],0x1  // [rbp-0x30] 既是数组a的地址，也是a[0]的地址 
   0x000055f61bcd82a1 <+34>:	mov    DWORD PTR [rbp-0x2c],0x2  // 以后每个元素4个字节 
   0x000055f61bcd82a8 <+41>:	mov    DWORD PTR [rbp-0x28],0x3  
   0x000055f61bcd82af <+48>:	mov    DWORD PTR [rbp-0x24],0x4
   0x000055f61bcd82b6 <+55>:	mov    DWORD PTR [rbp-0x20],0x5
   0x000055f61bcd82bd <+62>:	mov    DWORD PTR [rbp-0x1c],0x6
   0x000055f61bcd82c4 <+69>:	mov    DWORD PTR [rbp-0x18],0x7
   0x000055f61bcd82cb <+76>:	mov    DWORD PTR [rbp-0x14],0x8

7	    int b = a[2];
   0x000055f61bcd82d2 <+83>:	mov    eax,DWORD PTR [rbp-0x28]  // 直接使用数组偏移后的位置[rbp-0x30] [rbp-0x28] 
   0x000055f61bcd82d5 <+86>:	mov    DWORD PTR [rbp-0x34],eax
```

> 数组下标从0开始的原因就是地址的偏移量为0  

*二维数组*

实例及汇编实现  
```
6	    int a[2][4] = {{1,2,3,4},{5,6,7,8}};
=> 0x0000560568ded29a <+27>:	mov    DWORD PTR [rbp-0x30],0x1
   0x0000560568ded2a1 <+34>:	mov    DWORD PTR [rbp-0x2c],0x2
   0x0000560568ded2a8 <+41>:	mov    DWORD PTR [rbp-0x28],0x3
   0x0000560568ded2af <+48>:	mov    DWORD PTR [rbp-0x24],0x4
   0x0000560568ded2b6 <+55>:	mov    DWORD PTR [rbp-0x20],0x5
   0x0000560568ded2bd <+62>:	mov    DWORD PTR [rbp-0x1c],0x6
   0x0000560568ded2c4 <+69>:	mov    DWORD PTR [rbp-0x18],0x7
   0x0000560568ded2cb <+76>:	mov    DWORD PTR [rbp-0x14],0x8

7	    int b = a[2][3];
   0x0000560568ded2d2 <+83>:	mov    eax,DWORD PTR [rbp-0x4]
   0x0000560568ded2d5 <+86>:	mov    DWORD PTR [rbp-0x34],eax
```

可以看出一维数组与二维数组的内存排列是一致的，也是线性排列的，只有有两个索引位置，需要计算行列的值。  

### 动态数组Vector  
- 面向对象方式的动态数组

#### vector的数据结构
不同操作系统表示的方式不同。以下是MACOS系统Clion的  
源码路径:`/Library/Developer/CommandLineTools/SDKs/MacOSX11.3.sdk/usr/include/c++/v1/vector`    
![vector的数据结构](../../res/vector_struct.png)  

Centos的VSCODE的  
![Centos的VSCODE的](../../res/centos-clion.png)  

Ubuntu的Clion的, Ubuntu的vscode也是这样的:scream:      
![Ubuntu的Clion的](../../res/ubuntu-clion.png)

> 不同系统及不同版本的汇编实现都是有差异的  


实例
```c++
#include <iostream>
#include <vector>
using namespace std;

void printVec(string tag, vector<int> vs)
{
    cout << tag << " ";
    for(int i=0; i < vs.size(); i++)
    {
        if(i == vs.size() - 1)
        {
            cout << vs[i] << endl;
        }else
        {
            cout << vs[i] << ",";
        }
    }

    cout << "size:" << vs.size() << ",cap:" << vs.capacity() << endl;
}


int main() {
    vector<int> vs = {1,2,3,4};
    vs.push_back(5);
    printVec("init", vs);

    vs.insert(--vs.end(), 6);
    printVec("insert", vs);

    vs.pop_back();
    vs.erase(vs.end() - 2);
    printVec("delete", vs);

    return 0;
}
```

初始化和插入元素的汇编实现
```c++
Dump of assembler code for function main():
4	int main() {
   0x0000000100002260 <+0>:	push   rbp
   0x0000000100002261 <+1>:	mov    rbp,rsp
   0x0000000100002264 <+4>:	sub    rsp,0x70
   0x0000000100002268 <+8>:	mov    rax,QWORD PTR [rip+0x1db1]        # 0x100004020
   0x000000010000226f <+15>:	mov    rax,QWORD PTR [rax]
   0x0000000100002272 <+18>:	mov    QWORD PTR [rbp-0x8],rax
   0x0000000100002276 <+22>:	mov    DWORD PTR [rbp-0x1c],0x0

5	    vector<int> vs = {1,2,3,4};
=> 0x000000010000227d <+29>:	mov    DWORD PTR [rbp-0x18],0x1  // [rbp-0x18]是数组的地址
   0x0000000100002284 <+36>:	mov    DWORD PTR [rbp-0x14],0x2
   0x000000010000228b <+43>:	mov    DWORD PTR [rbp-0x10],0x3
   0x0000000100002292 <+50>:	mov    DWORD PTR [rbp-0xc],0x4
   0x0000000100002299 <+57>:	lea    rax,[rbp-0x18]            //[rbp-0x18] 赋值给rax 
   0x000000010000229d <+61>:	mov    QWORD PTR [rbp-0x48],rax  // rax-> [rbp-0x48]
   0x00000001000022a1 <+65>:	mov    QWORD PTR [rbp-0x40],0x4  // 元素个数为4 
   0x00000001000022a9 <+73>:	mov    rsi,QWORD PTR [rbp-0x48]
   0x00000001000022ad <+77>:	mov    rdx,QWORD PTR [rbp-0x40]
   0x00000001000022b1 <+81>:	lea    rax,[rbp-0x38]            // 第三个变量 [rbp-0x38] 
   0x00000001000022b5 <+85>:	mov    rdi,rax
   0x00000001000022b8 <+88>:	mov    QWORD PTR [rbp-0x68],rax  // 本地变量vs (vector<int> vs)
   0x00000001000022bc <+92>:	call   0x100002340 <std::__1::vector<int, std::__1::allocator<int> >::vector(std::initializer_list<int>)>

6	    vs.push_back(5);
   0x00000001000022c1 <+97>:	mov    DWORD PTR [rbp-0x4c],0x5  // 加载5为[rbp-0x4c]
   0x00000001000022c8 <+104>:	lea    rsi,[rbp-0x4c]
   0x00000001000022cc <+108>:	mov    rdi,QWORD PTR [rbp-0x68]  // 变量vs , push_back方法传递两个变量 this和5  
   0x00000001000022d0 <+112>:	call   0x100002370 <std::__1::vector<int, std::__1::allocator<int> >::push_back(int&&)>
   0x00000001000022d5 <+117>:	jmp    0x1000022da <main()+122>

7	
8	    return 0;
   0x00000001000022da <+122>:	mov    DWORD PTR [rbp-0x1c],0x0

9	}  // 在方法结束时，调用了vector的析构函数  
   0x00000001000022e1 <+129>:	lea    rdi,[rbp-0x38]
   0x00000001000022e5 <+133>:	call   0x1000023e0 <std::__1::vector<int, std::__1::allocator<int> >::~vector()>
   0x00000001000022ea <+138>:	mov    eax,DWORD PTR [rbp-0x1c]
   0x00000001000022ed <+141>:	mov    rcx,QWORD PTR [rip+0x1d2c]        # 0x100004020
   0x00000001000022f4 <+148>:	mov    rcx,QWORD PTR [rcx]
   0x00000001000022f7 <+151>:	mov    rdx,QWORD PTR [rbp-0x8]
   0x00000001000022fb <+155>:	cmp    rcx,rdx
   0x00000001000022fe <+158>:	mov    DWORD PTR [rbp-0x6c],eax
   0x0000000100002301 <+161>:	jne    0x10000232b <main()+203>
   0x0000000100002307 <+167>:	mov    eax,DWORD PTR [rbp-0x6c]
   0x000000010000230a <+170>:	add    rsp,0x70
   0x000000010000230e <+174>:	pop    rbp
   0x000000010000230f <+175>:	ret    
   0x0000000100002310 <+176>:	mov    QWORD PTR [rbp-0x58],rax
   0x0000000100002314 <+180>:	mov    DWORD PTR [rbp-0x5c],edx
   0x0000000100002317 <+183>:	lea    rdi,[rbp-0x38]
   0x000000010000231b <+187>:	call   0x1000023e0 <std::__1::vector<int, std::__1::allocator<int> >::~vector()>
   0x0000000100002320 <+192>:	mov    rdi,QWORD PTR [rbp-0x58]
   0x0000000100002324 <+196>:	call   0x100003cfe
   0x0000000100002329 <+201>:	ud2    
   0x000000010000232b <+203>:	call   0x100003d5e
   0x0000000100002330 <+208>:	ud2    
   0x0000000100002332 <+210>:	nop    WORD PTR cs:[rax+rax*1+0x0]
   0x000000010000233c <+220>:	nop    DWORD PTR [rax+0x0]

End of assembler dump.
```



`c++/v1/vector`文件中vector的创建及`push_back`方法  
创建方法
```
template <class _Tp, class _Allocator>
inline _LIBCPP_INLINE_VISIBILITY
vector<_Tp, _Allocator>::vector(initializer_list<value_type> __il)
{
#if _LIBCPP_DEBUG_LEVEL >= 2
    __get_db()->__insert_c(this);
#endif
    if (__il.size() > 0)
    {
        __vallocate(__il.size());
        __construct_at_end(__il.begin(), __il.end(), __il.size());
    }
}
```

添加元素方法  
```c++
template <class _Tp, class _Allocator>
inline _LIBCPP_INLINE_VISIBILITY
void
vector<_Tp, _Allocator>::push_back(value_type&& __x)
{
    if (this->__end_ < this->__end_cap())
    {
        __construct_one_at_end(_VSTD::move(__x));
    }
    else
        __push_back_slow_path(_VSTD::move(__x));
}
```

### 字符串  
- 以字符'\0'结束  

汇编实现
```c++
35	    char str[] = {"HelloWorld"};
                                // 编译器在解释字符串时，就在最后增加了'\0'结束符
   0x0000000100003dfc <+652>:	mov    rax,QWORD PTR [rip+0x3fd4]        # 0x100007dd7 => "HelloWorld\0" 
   0x0000000100003e03 <+659>:	mov    QWORD PTR [rbp-0x23],rax
   0x0000000100003e07 <+663>:	mov    cx,WORD PTR [rip+0x3fd1]        # 0x100007ddf => allocator<T>...
   0x0000000100003e0e <+670>:	mov    WORD PTR [rbp-0x1b],cx
   0x0000000100003e12 <+674>:	mov    dl,BYTE PTR [rip+0x3fc9]        # 0x100007de1
   0x0000000100003e18 <+680>:	mov    BYTE PTR [rbp-0x19],dl
```

查看内存地址
```shell 
(gdb) x/32c 0x100007dd7
0x100007dd7:	72 'H'	101 'e'	108 'l'	108 'l'	111 'o'	87 'W'	111 'o'	114 'r'
0x100007ddf:	108 'l'	100 'd'	0 '\000'	97 'a'	108 'l'	108 'l'	111 'o'	99 'c'
0x100007de7:	97 'a'	116 't'	111 'o'	114 'r'	60 '<'	84 'T'	62 '>'	58 ':'
0x100007def:	58 ':'	97 'a'	108 'l'	108 'l'	111 'o'	99 'c'	97 'a'	116 't'

(gdb) x/32s 0x100007ddf
0x100007ddf:	"ld"
0x100007de2:	"allocator<T>::allocate(size_t n) 'n' exceeds maximum supported size"
0x100007e26:	""
0x100007e27:	""
0x100007e28:	"\001"

(gdb) x/32s 0x100007de1
0x100007de1:	""
0x100007de2:	"allocator<T>::allocate(size_t n) 'n' exceeds maximum supported size"
0x100007e26:	""
0x100007e27:	""
0x100007e28:	"\001"
```

> 字符串创建时"HelloWorld\0"已经存在，接着调用`allocator<T>::allocate(size_t n)`，`vector创建时`也会调用该方法 ？  



## 指针
## 基础句法
## 高级语法
## 编程思想
## 进阶编程
## GUI开发
## 陷阱与经验