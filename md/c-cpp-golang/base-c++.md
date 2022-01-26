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

数据结构如下:
```
(vector)._M_impl
((vector)._M_impl)._M_start
((vector)._M_impl)._M_finish
((vector)._M_impl)._M_end_of_storage

```

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

初始化和插入元素的汇编实现`macos clion`
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

从clion中可以看出方法的原始名称不需要转换，如果使用vscode连接centos7虚拟机的汇编实现如下:
```c++
-exec disass /m
Dump of assembler code for function main():
23	int main() {
   0x0000000000400dc6 <+0>:	push   rbp
   0x0000000000400dc7 <+1>:	mov    rbp,rsp
   0x0000000000400dca <+4>:	push   r13
   0x0000000000400dcc <+6>:	push   r12
   0x0000000000400dce <+8>:	push   rbx
   0x0000000000400dcf <+9>:	sub    rsp,0x28

24	    vector<int> vs = {1,2,3,4};
   0x0000000000400dd3 <+13>:	lea    rax,[rbp-0x25]
   0x0000000000400dd7 <+17>:	mov    rdi,rax
   0x0000000000400dda <+20>:	call   0x400f58 <_ZNSaIiEC2Ev>  // std::allocator<int>::allocator()
   0x0000000000400ddf <+25>:	mov    r12d,0x401ec0
   0x0000000000400de5 <+31>:	mov    r13d,0x4
   0x0000000000400deb <+37>:	lea    rdi,[rbp-0x25]
   0x0000000000400def <+41>:	mov    rcx,r12
   0x0000000000400df2 <+44>:	mov    rbx,r13
   0x0000000000400df5 <+47>:	mov    rax,r12
   0x0000000000400df8 <+50>:	mov    rdx,r13
   0x0000000000400dfb <+53>:	mov    rsi,rcx
   0x0000000000400dfe <+56>:	lea    rax,[rbp-0x40]
   0x0000000000400e02 <+60>:	mov    rcx,rdi
   0x0000000000400e05 <+63>:	mov    rdi,rax
   0x0000000000400e08 <+66>:	call   0x400fe6 <_ZNSt6vectorIiSaIiEEC2ESt16initializer_listIiERKS0_>  //std::vector<int, std::allocator<int> >::vector(std::initializer_list<int>, std::allocator<int> const&)
   0x0000000000400e0d <+71>:	lea    rax,[rbp-0x25]
   0x0000000000400e11 <+75>:	mov    rdi,rax
   0x0000000000400e14 <+78>:	call   0x400f72 <_ZNSaIiED2Ev> // std::allocator<int>::~allocator() 

25	    vs.push_back(5);
   0x0000000000400e19 <+83>:	mov    DWORD PTR [rbp-0x24],0x5
   0x0000000000400e20 <+90>:	lea    rdx,[rbp-0x24]
   0x0000000000400e24 <+94>:	lea    rax,[rbp-0x40]
   0x0000000000400e28 <+98>:	mov    rsi,rdx
   0x0000000000400e2b <+101>:	mov    rdi,rax
   0x0000000000400e2e <+104>:	call   0x4010c8 <_ZNSt6vectorIiSaIiEE9push_backEOi> // std::vector<int, std::allocator<int> >::push_back(int&&)
   
26	
27	    return 0;
=> 0x0000000000400e33 <+109>:	mov    ebx,0x0
   0x0000000000400e38 <+114>:	lea    rax,[rbp-0x40]
   0x0000000000400e3c <+118>:	mov    rdi,rax
   0x0000000000400e3f <+121>:	call   0x401076 <_ZNSt6vectorIiSaIiEED2Ev>   // std::vector<int, std::allocator<int> >::~vector()  
   0x0000000000400e44 <+126>:	mov    eax,ebx
   0x0000000000400e46 <+128>:	jmp    0x400e7c <main()+182>
   0x0000000000400e48 <+130>:	mov    rbx,rax
   0x0000000000400e4b <+133>:	lea    rax,[rbp-0x25]
   0x0000000000400e4f <+137>:	mov    rdi,rax
   0x0000000000400e52 <+140>:	call   0x400f72 <_ZNSaIiED2Ev>  // std::allocator<int>::~allocator()
   0x0000000000400e57 <+145>:	mov    rax,rbx
   0x0000000000400e5a <+148>:	mov    rdi,rax
   0x0000000000400e5d <+151>:	call   0x400b80 <_Unwind_Resume@plt>
   0x0000000000400e62 <+156>:	mov    rbx,rax
   0x0000000000400e65 <+159>:	lea    rax,[rbp-0x40]
   0x0000000000400e69 <+163>:	mov    rdi,rax
   0x0000000000400e6c <+166>:	call   0x401076 <_ZNSt6vectorIiSaIiEED2Ev> // std::vector<int, std::allocator<int> >::~vector()  
   0x0000000000400e71 <+171>:	mov    rax,rbx
   0x0000000000400e74 <+174>:	mov    rdi,rax
   0x0000000000400e77 <+177>:	call   0x400b80 <_Unwind_Resume@plt>

28	}
   0x0000000000400e7c <+182>:	add    rsp,0x28
   0x0000000000400e80 <+186>:	pop    rbx
   0x0000000000400e81 <+187>:	pop    r12
   0x0000000000400e83 <+189>:	pop    r13
   0x0000000000400e85 <+191>:	pop    rbp
   0x0000000000400e86 <+192>:	ret    
End of assembler dump.
```

可以使用`c++filt`工具，centos安装指令`yum install binutils`
```shell
c++filt _ZNSt6vectorIiSaIiEEC2ESt16initializer_listIiERKS0_
std::vector<int, std::allocator<int> >::vector(std::initializer_list<int>, std::allocator<int> const&)
```
做一个 [小工具]() 自动转换该指令，翻译后的结果就是正常的函数名了  
```go
package main

import (
	"container/list"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

func main() {

	// 读取输入文件
	currPath, _ := os.Getwd()

	// 输入文件路径, 使用os.Stdin vscode命令行无法交互输入
	file, err := os.Open(currPath + "/inputfile")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	contents, err := ioutil.ReadAll(file)
	strs := string(contents)
	lines := strings.Split(strs, "\n")
	results := list.New()

	// 解析数据
	for _, line := range lines {
		if strings.Contains(line, "<_ZN") {
			start := strings.Index(line, "<_ZN")
			end := strings.LastIndex(line, ">")

			funcSymbol := line[start+1 : end]
			// 执行c++filt 命令
			cmd := exec.Command("/usr/bin/c++filt", funcSymbol)
			out, _ := cmd.CombinedOutput()
			funcName := string(out)

			line = line[:start] + funcName
		}
		results.PushBack(line)
	}

	for i := results.Front(); i != nil; i = i.Next() {
		fmt.Println(i.Value)
	}
```

目前可以看到`vector`的内存创建与`std::allocator<int>::allocator()`和`std::initializer_list<int>`有关系  

> 需要注意c++源码的版本，目前在centos使用`gcc-g++ 4.8.5`版本，对应的源码为gcc-g++-4.8.2-x86_64-1.txz
> 路径是"/usr/include/c++/4.8.2/bits/stl_vector.h"  

> vs版本、gcc版本、c++版本之间的关系   
> c++版本是一个标准，需要编译器支持。比如c++11标准, `gcc4.8.1`及以上可以完全支持。`vs2015`可以完全支持。  


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
   0x0000000100003dfc <+652>:	mov    rax,QWORD PTR [rip+0x3fd4]       # 0x100007dd7 => "HelloWorld\0" 
   0x0000000100003e03 <+659>:	mov    QWORD PTR [rbp-0x23],rax
   0x0000000100003e07 <+663>:	mov    cx,WORD PTR [rip+0x3fd1]        # 0x100007ddf => allocator<T>...
   0x0000000100003e0e <+670>:	mov    WORD PTR [rbp-0x1b],cx
   0x0000000100003e12 <+674>:	mov    dl,BYTE PTR [rip+0x3fc9]        # 0x100007de1
   0x0000000100003e18 <+680>:	mov    BYTE PTR [rbp-0x19],dl
```

查看内存地址
```shell 
(gdb) x/32c 0x100007dd7    // 字符串结束时有个结束符'\0'  
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

查看不同0值的含义  
```c++
#include <iostream>
using namespace std;

int main()
{
    char a1 = 0;    // 0x00
    char a2= '\0';  // 0x00
    char a3 = '0';  // 0x30
    return 0;
}
```

通过gdb查看内存值  
```
-exec x/bb &a1
0x7fffffffdcff:	0x00

-exec x/bb &a2
0x7fffffffdcfe:	0x00

-exec x/bb &a3
0x7fffffffdcfd:	0x30

-exec x/db &a3        // d 十进制显示   b显示单位为byte  
0x7fffffffdcfd:	48
```

从中可以看出`0`与`'\0'`是等价的  

#### unicode编码  
- Unicode编码:最初的目的是把世界上的文字都映射到一套字符空间中  
- 为了表示Unicode字符集，有3种(确切的说是5种)Unicode的编码方式:  
  - UTF-8: 1 byte来表示字符，可以兼容ASCII码; 特点是存储效率高，变长(不方便内部随机访问)，无字节序问题(可作为外部编码)  
  - UTF-16:分为UTF-16BE(big endian), UTF-16LE(ittle endian),特点是定长(方便内部随机访问), 有字节序问题(不可作为外部编码)  
  - UTF-32:分为UTF-32BE(big endian), UTF-32LE(ittle endian),特点是定长(方便内部随机访问), 有字节序问题(不可作为外部编码)
- 编码错误的根本原因在于编码方式和解码方式的不统一;  

> Windows的文件可能有BOM(byte order mark)如要在其他平台使用，可以去掉BOM  

```
// utf-8 格式
68 65 6C 6C 6F 77 6F 72 6C 64   |   helloworld  

//utf-16 LE   在utf-8的基础上增加一个byte 00 
68 00 65 00 6C 00 6C 00 6F 00 77 00 6F 00 72 00 6C 00 64 00   |     h⊙e⊙l⊙l⊙o⊙w⊙o⊙r⊙l⊙d⊙  

//utf-16 BE
00 68 00 65 00 6C 00 6C 00 6F 00 77 00 6F 00 72 00 6C 00 64   |     ⊙h⊙e⊙l⊙l⊙o⊙w⊙o⊙r⊙l⊙d  
```


#### 字符串指针  
> 指向常量区的字符串指针，内容不能被修改。(常量区只读)。 golang中的字符串内容是可以修改的，没有常量区不能修改这一说。  

```c++
#include <string.h>
#include <iostream>
using namespace std;
int main()
{
    // 定义一个数组
    char strHelloWorld[11] = {"helloworld"}; // 这个定义可以, 字符存放到数组中  
    char *pStrHelloWrold = "helloworld";     // 常量区的值是不可以改变的
    pStrHelloWrold = strHelloWorld;
    //strHelloWorld = pStrHelloWrold;        // 数组变量的值不允许改变

    // 通过数组变量遍历修改数组中的元素值
    for (int index = 0; index < strlen(strHelloWorld); ++index)
    {
        strHelloWorld[index] += 1;
        std::cout << strHelloWorld[index];
    }
    cout << endl; // 换行  

    // 通过指针变量遍历修改数组中的元素值  
    for (int index = 0; index < strlen(strHelloWorld); ++index)
    {
        pStrHelloWrold[index] += 1;  
        std::cout << pStrHelloWrold[index];  
    }

    cout << endl; // 换行
    
    // 计算字符串长度
    cout << "字符串长度为: " << strlen(strHelloWorld) << endl;
    cout << "字符串占用空间为:  " << sizeof(strHelloWorld) << endl;

    return 0;
}
```

输出结果 
```
ifmmpxpsme
jgnnqyqtnf
字符串长度为: 10
字符串占用空间为:  11
```

### 字符串基本操作  
```c++
/* Copy SRC to DEST.  */
extern char *strcpy (char *__restrict __dest, const char *__restrict __src)
/* Copy no more than N characters of SRC to DEST.  */
extern char *strncpy (char *__restrict __dest,const char *__restrict __src, size_t __n)

/* Append SRC onto DEST.  */
extern char *strcat (char *__restrict __dest, const char *__restrict __src)

/* Compare S1 and S2.  */  // 0 相等; 正数 S1 > S1; 负数 S1 < S1;  
extern int strcmp (const char *__s1, const char *__s2)

/* Return the length of S.  */
extern size_t strlen (const char *__s)

/* Find the first occurrence of C in S.  */   // 第一次出现char的位置  
extern char *strchr (char *__s, int __c)

/* Find the first occurrence of NEEDLE in HAYSTACK.  */  // 字符串查找  
extern char *strstr (char *__haystack, const char *__needle)



```

基础操作示例`<string.h> 使用C库的头文件`  
```
#include <string.h> //使用C库的头文件
#include <iostream>
using namespace std;
const unsigned int MAX_LEN_NUM = 16;
const unsigned int STR_LEN_NUM = 7;
const unsigned int NUM_TO_COPY = 2;
int main()
{
    char strHelloWorld1[] = {"hello"};
    char strHelloWorld2[STR_LEN_NUM] = {"world1"};
    char strHelloWorld3[MAX_LEN_NUM] = {0};

    // 字符串拷贝  (dest, src)
    strcpy(strHelloWorld3, strHelloWorld1); // hello
    // strcpy_s(strHelloWorld3, MAX_LEN_NUM, strHelloWorld1);

    int res = strcmp(strHelloWorld1, strHelloWorld2);
    cout << "cmp result:" << res << ";" << strHelloWorld1 << ":" << strHelloWorld2 << endl;

    strncpy(strHelloWorld3, strHelloWorld2, NUM_TO_COPY); // wollo
    // strncpy_s(strHelloWorld3, MAX_LEN_NUM,  strHelloWorld2, NUM_TO_COPY);

    // 字符串拼接(dest, src), 追加的方式    
    strcat(strHelloWorld3, strHelloWorld2); //  wolloworld1
    // strcat_s(strHelloWorld3, MAX_LEN_NUM, strHelloWorld2);

    unsigned int len = strlen(strHelloWorld3);
    // unsigned int len = strnlen_s(strHelloWorld3, MAX_LEN_NUM);
    for (unsigned int index = 0; index < len; ++index)
    {
        cout << strHelloWorld3[index] << "";
    }
    cout << endl;

    // 小心缓冲区溢出
    strcat(strHelloWorld2, "Welcome to C++");
    // strcat_s(strHelloWorld2, STR_LEN_NUM, "Welcome to C++");

    return 0;
}
```

> c中原始字符串的安全性和效率存在一定问题  
> 缓冲区溢出、strlen的效率可以提高(空间换时间)  
> redis字符串中增加了长度信息，可以直接查询    

```c
typedef char *sds;


struct sdshdr {

    // buf 已占用长度
    int len;

    // buf 剩余可用长度
    int free;

    // 实际保存字符串数据的地方
    char buf[];
};
```

另一种字符串`<string> Forward declarations -*- C++ -*-`文件名为`stringfwd.h`    
```
  /// A string of @c char
  typedef basic_string<char>    string;  
```

- c++标准库中提供了string类型专门表示字符串 `#include <string>` 
- 使用string可以更为方便和安全的管理字符串  
- 定义字符串变量的方式: `string s` 、`string s = "hello"`、`string s("helloworld")`  

基础操作示例:
```
#include <iostream>
#include <string>
using namespace std;
int main()
{
    // 字符串定义
    string s1;                //定义空字符串
    string s2 = "helloworld"; //定义并初始化
    string s3("helloworld");
    string s4 = string("helloworld");

    // 获取字符串长度  length与size 等价的
    cout << "length:" << s2.length() << endl;     // _M_length
    cout << "size:" << s2.size() << endl;         // _M_length
    cout << "capacity:" << s2.capacity() << endl; // _M_capacity

    // 字符串比较
    s1 = "hello", s2 = "world";
    cout << "(s1==s2):" << (s1 == s2) << endl;
    cout << "(s1!=s2):" << (s1 != s2) << endl;
    cout << "(s1<s2):" << (s1 < s2) << endl;

    //  转换成C风格的字符串
    const char *c_str1 = s1.c_str();
    cout << "The C-style string c_str1 is: " << c_str1 << endl;
    //  随机访问
    for (unsigned int index = 0; index < s1.length(); ++index)
    {
        cout << c_str1[index] << " ";
    }
    cout << endl;
    for (unsigned int index = 0; index < s1.length(); ++index)
    {
        cout << s1[index] << " ";
    }
    cout << endl;

    // 字符串拷贝
    s1 = "helloworld";
    s2 = s1;

    // 字符串连接
    s1 = "helllo", s2 = "world";
    s3 = s1 + s2; //s3: helloworld
    s1 += s2;     //s1: helloworld
    return 0;
}
```

> string 结合了C++的新特性，使用起来比原始的C风格方法更安全和方便，对性能要求不是特别高的场景可以使用。  

## 指针
### 数组指针和指针数组 
- 数组的指针 `T (*pA)[]`   
- 指针的数组 `T* a[]`  

示例代码
```c++
#include <iostream>
using namespace std;

int main()
{
    int a[4] = {0x80000000, 0x000000FF, 3, 4};
    int *pA = a; // &a[0]  // 指针指向数组首地址
    cout << "a数组地址:" << a << ",pA[0]:" << pA[0] << ",pA[1]:" << pA[1] << endl;

    int* pArr[4] = {&a[0], &a[1], &a[2], &a[3]};  // 指针的数组   T pArr[4]  => T = int* 
    // 输出的是指针的地址 0x7fffffffdcd0  
    cout << "pArr[0]:" << pArr[0] << ",pArr[1]:" << pArr[1] << endl;
    cout << "*(pArr[0]):" << *(pArr[0]) << ",*(pArr[1):" << *(pArr[1]) << endl;

    // 指向int[4]数组的指针  a pointer to an array 
    int (*arrP)[4];    // 数组的指针  T (*pArr)[4] => T = int
    arrP = &a;         // 数组的个数必须一致  

    // arrP[0] 相当于指针的偏移操作，偏移量为0,那就指向a数组的地址  //arrP[1] 偏移16个字节(4字节x4个元素)  
    cout << "arrP[0]:" << arrP[0] << ",arrP[1]:" << arrP[1] << endl;
    // (*arrP) 是指向数组的地址，(*arrP)[0] 是指向数组地址的首个元素  
    cout << "(*arrP)[0]:" << (*arrP)[0] << ",(*arrP)[1]:" << (*arrP)[1] << endl;

}
```

打印输出 
```shell
a数组地址:0x7fffffffdcd0,pA[0]:-2147483648,pA[1]:255
pArr[0]:0x7fffffffdcd0,pArr[1]:0x7fffffffdcd4
arrP[0]:0x7fffffffdcd0,arrP[1]:0x7fffffffdce0       
(*arrP)[0]:-2147483648,(*arrP)[1]:255

a数组地址:0x7fffffffdcd0,pA[0]:-2147483648,pA[1]:255
pArr[0]:0x7fffffffdcd0,pArr[1]:0x7fffffffdcd4
*(pArr[0]):-2147483648,*(pArr[1):255                // pArr[0] 是数组的地址, *(pArr[0]) 取值  
arrP[0]:0x7fffffffdcd0,arrP[1]:0x7fffffffdce0       // arrP[0] 数组a的地址, arrP[0] 偏移16个字节的地址
(*arrP)[0]:-2147483648,(*arrP)[1]:255
```  

> 需要特别注意的是`c/c++`中指针(内存地址)是可以进行运算的，一种方式是`arrP + 1`，另一种方式是`arrP[1]` 含义是相同的  

### const与指针  
- const pointer  
- pointer to const 

> 关于const修改部分，看左侧最近的部分，如果左侧没有，则看右侧。  
> **主要确认修改的部分是`(*p) 指向的内容`还是`(p) 指针`**    

```c++
#include <iostream>
#include <string.h>
using namespace std;
unsigned int MAX_LEN = 11;

int main()
{
    char strHelloworld[] = {"helloworld"};
    // const修饰谁，谁的内容就不可变，其他的都可变
    const char *pStr1 = "helloworld";       // 修饰的是(*pStr1) 指向的内容不能修改。 主要用于函数的形参，防止被修改  
    char *const pStr2 = strHelloworld;      // 修饰的是(pStr2)  指向不能修改
    const char *const pStr3 = "helloworld"; // 修饰的是(*pStr3)和(pStr3)  指向和指向的内容都不能修改

    pStr1 = strHelloworld; // 指向可以变，指向的内容不能变
    //pStr2 = strHelloworld;                // pStr2不可改
    //pStr3 = strHelloworld;                // pStr3不可改

    unsigned int len = strnlen(pStr2, MAX_LEN);
    cout << len << endl;
    for (unsigned int index = 0; index < len; ++index)
    {
        //pStr1[index] += 1;                               // pStr1里的值不可改
        pStr2[index] += 1;
        //pStr3[index] += 1;                               // pStr3里的值不可改
    }

    char a = 'a';
    const char *pA = &a;

    // 指针指向的内容不能修改， 但是a的值仍是可以修改的  
    // *pA = 'b'; // assignment of read-only location ‘* pA’
    a = 'c';      // (*pA)指向的内容就是变量a的地址，a变量仍然是可以修改的  
    cout << *pA << endl;

    return 0;
}
```

查看内存值
```shell
-exec x/16xb &strHelloworld
0x7fffffffdcd0:	0x68	0x65	0x6c	0x6c	0x6f	0x77	0x6f	0x72
0x7fffffffdcd8:	0x6c	0x64	0x00	0x00	0x00	0x00	0x00	0x00

-exec x/16xb &pStr1
0x7fffffffdcf0:	0xd0	0xdc	0xff	0xff	0xff	0x7f	0x00	0x00
0x7fffffffdcf8:	0x00	0x00	0x00	0x00	0x00	0x00	0x00	0x00

-exec x/16xb &pStr2
0x7fffffffdce8:	0xd0	0xdc	0xff	0xff	0xff	0x7f	0x00	0x00
0x7fffffffdcf0:	0xd0	0xdc	0xff	0xff	0xff	0x7f	0x00	0x00

-exec x/16xb &pStr3
0x7fffffffdce0:	0x81	0x09	0x40	0x00	0x00	0x00	0x00	0x00
0x7fffffffdce8:	0xd0	0xdc	0xff	0xff	0xff	0x7f	0x00	0x00

-exec x/16cb 0x400981
0x400981:	104 'h'	101 'e'	108 'l'	108 'l'	111 'o'	119 'w'	111 'o'	114 'r'
0x400989:	108 'l'	100 'd'	0 '\000'	1 '\001'	27 '\033'	3 '\003'	59 ';'	64 '@'
```

通过查看内存值可以画出如下关系:  

<br>
<div align=center>
    <img src="../../res/const与指针.jpg" width="60%" height="60%" title="队列参数设置"></img>  
</div>
<br>

`pStr1`、`pStr2`、`pStr3`变量的地址是连续的，`pStr1`、`pStr2`指向  


`const`还是针对于编译器，可以从汇编实现看出，运行时与其他变量无异!  
- const 是由编译器进行处理，执行类型检查和作用域的检查；
- define 是由预处理器进行处理，只做简单的文本替换工作而已。

### 二级指针和野指针  

指向指针的指针
```c++
#include <iostream>
using namespace std;

int main()
{
    int a = 123;
    int* b = &a;
    int** c = &b;

    cout << "&a=" << &a << ",b=" << b << ",c=" << c << endl;
}
```
输出结果  
```
&a=0x7fffffffdce4,b=0x7fffffffdce4,c=0x7fffffffdcd8
```

如果一个指针在声明时没有指定初始值，那系统分配的初始值可能是异常数据，最坏的情况是，地址可以访问，导致修改了其他有效数据。  

```
#include <iostream>
using namespace std;

int main()
{
    int *a;  // 分配的地址为 0x400730 ,可能是非法地址  
    cout << "a=" << a << ",&a=" << &a << endl;

    /*
    -exec x/16xb 0x400730
    0x400730 <_start>:	    0x31	0xed	0x49	0x89	0xd1	0x5e	0x48	0x89
    0x400738 <_start+8>:	0xe2	0x48	0x83	0xe4	0xf0	0x50	0x54	0x49
    */
    return 0;
}
```

> 指针在不适用的时候要把指针置空(NULL)  
> 在C++中建议使用nullptr替代NULL，因为在C++中NULL是: #define NULL 0 这样在整型重载的时候可能会有问题。而C++11加入了nullptr，可以保证在任何情况下都代表空指针， 所以比较安全。

```c++
#include <iostream>
using namespace std;
int main()
{
    // 指针的指针
    int a = 123;
    int *b = &a;
    int **c = &b;

    // NULL 的使用
    int *pA = NULL;
    pA = &a;
    if (pA != NULL) //  判断NULL指针
    {
        cout << (*pA) << endl;
    }
    pA = NULL; //  pA不用时，置为NULL

    return 0;
}
```

**野指针**  
- 指向`垃圾`内存的指针， if判断对它没有作用，因为没有置空

一般有三种情况:  
- 指针变量没有初始化;
- 已经释放不用的指针没有置NULL, 如`delete`与`free`之后的指针;
- 指针操作超越了变量的作用范围;  

> 没有初始化的指针，不用或者超出作用范围的指针请把值置为NULL  


### 指针的基本操作  
- `&`和`*`
- `++`、`--`  
- 

```c++
#include <iostream>
using namespace std;
int main()
{
    char ch = 'a';

    // &操作符
    //&ch = 97;     // &ch左值不合法
    char *cp = &ch; // &ch右值
    //&cp = 97;     // &cp左值不合法
    char **cpp = &cp; // &cp右值

    // *操作符
    *cp = 'a';      // *cp左值取变量ch位置
    char ch2 = *cp; // *cp右值取变量ch存储的值
    //*cp + 1 = 'a'; //  *cp+1左值不合法的位置
    ch2 = *cp + 1;   //  *cp+1右值取到的字符做ASCII码+1操作
    *(cp + 1) = 'a'; //  *(cp+1)左值语法上合法，取ch后面位置
    ch2 = *(cp + 1); //  *(cp+1)右值语法上合法，取ch后面位置的值

    return 0;
}
```

```c++
int main()
{
    char ch = 'a';
    char *cp = &ch;
    // ++,--操作符
    char *cp2 = ++cp;
    char *cp3 = cp++;
    char *cp4 = --cp;
    char *cp5 = cp--;

    // ++ 左值
    //++cp2 = 97;
    //cp2++ = 97;

    // *++, ++*
    *++cp2 = 98;
    char ch3 = *++cp2;
    *cp2++ = 98;
    char ch4 = *cp2++;

    // ++++, ----操作符等
    int a = 1, b = 2, c, d;
    //c = a++b;               // error
    c = a++ + b;
    //d = a++++b;             // error
    char ch5 = ++*++cp;

    return 0;
}
```

### CPP程序的存储区域划分  
- (stack)栈区  
- 常量区  
- (heap)堆区
- (text)代码区, 调用函数时使用代码区的地址  
- (GVAR)全局初始化区
- (bss)全局未初始化区

```c++
#include <string.h>
#include <iostream>
using namespace std;

int a = 0; //(GVAR)全局初始化区
int *p1;   //(bss)全局未初始化区

int main() //(text)代码区
{
    int b = 1;                //(stack)栈区变量
    char s[] = "abc";         //(stack)栈区变量
    int *p2 = NULL;           //(stack)栈区变量
    char *p3 = "123456";      //123456\0在常量区, p3在(stack)栈区
    static int c = 0;         //(GVAR)全局(静态)初始化区
    p1 = new int(10);         //(heap)堆区变量
    p2 = new int(20);         //(heap)堆区变量
    char *p4 = new char[7];   //(heap)堆区变量
    strncpy(p4, "123456", 7); //(text)代码区   strncpy方法名是代码区的地址  

    cout << "b=" << &b << ",s=" << &s << ",p1=" << p1 << ",p2=" << p2 << endl;
    cout << "&p3=" << &p3 << ",&p4=" << &p4 << ",(void*)p3=" << (void *)p3 << ",(void*)p4=" << (void *)p4 << ",p3=" << p3 << ",p4=" << p4 << endl;

    //(text)代码区
    if (p1 != NULL)
    {
        delete p1;
        p1 = NULL;
    }
    if (p2 != NULL)
    {
        delete p2;
        p2 = NULL;
    }
    if (p4 != NULL)
    {
        delete[] p4;
        p4 = NULL;
    }
    //(text)代码区
    return 0; //(text)代码区
}
```

输出结果:  
```shell
b=0x7fffffffdcd4,s=0x7fffffffdcd0,p1=0x603010,p2=0x603030
&p3=0x7fffffffdcc8,&p4=0x7fffffffdcc0,(void*)p3=0x400c61,(void*)p4=0x603050,p3=123456,p4=123456
```

`b`和`s`在同一栈区(0x7fffffff)，`p1`、`p2`、`p4`指向的内存地址在堆区(0x6030)，`p3`指向的内存地址在常量区(0x400c)  

> `p3`指向常量字符串，如何打印字符串的地址呢?  `(void*)p3`
```
-exec p p3
$1 = 0x400c61 "123456"

-exec p p4
$2 = 0x603050 "123456"
```

<br>
<div align=center>
    <img src="../../res/cpp代码存储区域.png" width="80%" height="80%" title="队列参数设置"></img>  
</div>
<br>


## 基础句法
## 高级语法
## 编程思想

### 泛型编程思想  
- 如果说面向对象是一种通过间接层来调用函数，以换取一种抽象，那么泛型编程则是更直接的抽象，它不会因为间接层而损失效率;
- 不同于面向对象的动态期多态，泛型编程是一种静态期多态，通过编译器生成最直接的代码；
- 泛型编程可以将算法与特定类型、结构剥离，尽可能复用代码；

样例
```c++
#include <iostream>
#include <string.h>
// 不能使用 using namespace std; 
template <typename V> // template <class V>  两种方式一样
V max(V a, V b)
{
    return a > b ? a : b;
}

// 特例  
template<>
char* max(char* a, char* b)
{
    return (strcmp(a, b) > 0 ? a : b);
}
// 两种不同类型比较  
template<typename V1, typename V2>
int max(V1 a, V2 b)
{
    return static_cast<int>(a > b ? a : b);
}

int main() 
{
    std::cout << max(2, 4) << std::endl;
    std::cout << max(7.1, 5.6) << std::endl;
    std::cout << max('a', 'b') << std::endl;

    char* x = "1234";
    char* y = "2345";
    std::cout << max(x, y) << std::endl;

    std::cout << max(10, 12.8) << std::endl;
    return 0;
}
```

输出结果:
```
4
7.1
b
2345
12
```

从汇编实现中可以看出`max`函数在编译时被替换成对应的类型   
比如:
- max(2, 4) => int max<int>(int, int)  
- max('a', 'b') => char max<char>(char, char)  
- max(10, 12.8) => int max<int, double>(int, double)  

以下是汇编实现  
```c++
Dump of assembler code for function main():
24	{
   0x0000000000400960 <+0>:	push   rbp
   0x0000000000400961 <+1>:	mov    rbp,rsp
   0x0000000000400964 <+4>:	sub    rsp,0x20

25	    std::cout << max(2, 4) << std::endl;
   0x0000000000400968 <+8>:	mov    esi,0x4
   0x000000000040096d <+13>:	mov    edi,0x2
   0x0000000000400972 <+18>:	call   0x400adb int max<int>(int, int)

   0x0000000000400977 <+23>:	mov    esi,eax
   0x0000000000400979 <+25>:	mov    edi,0x602080
   0x000000000040097e <+30>:	call   0x400790 _ZNSolsEi@plt

   0x0000000000400983 <+35>:	mov    esi,0x400830
   0x0000000000400988 <+40>:	mov    rdi,rax
   0x000000000040098b <+43>:	call   0x400820 _ZNSolsEPFRSoS_E@plt


26	    std::cout << max(7.1, 5.6) << std::endl;
   0x0000000000400990 <+48>:	movabs rdx,0x4016666666666666
   0x000000000040099a <+58>:	movabs rax,0x401c666666666666
   0x00000000004009a4 <+68>:	mov    QWORD PTR [rbp-0x18],rdx
   0x00000000004009a8 <+72>:	movsd  xmm1,QWORD PTR [rbp-0x18]
   0x00000000004009ad <+77>:	mov    QWORD PTR [rbp-0x18],rax
   0x00000000004009b1 <+81>:	movsd  xmm0,QWORD PTR [rbp-0x18]
   0x00000000004009b6 <+86>:	call   0x400af7 double max<double>(double, double)

   0x00000000004009bb <+91>:	movsd  QWORD PTR [rbp-0x18],xmm0
   0x00000000004009c0 <+96>:	mov    rax,QWORD PTR [rbp-0x18]
   0x00000000004009c4 <+100>:	mov    QWORD PTR [rbp-0x18],rax
   0x00000000004009c8 <+104>:	movsd  xmm0,QWORD PTR [rbp-0x18]
   0x00000000004009cd <+109>:	mov    edi,0x602080
   0x00000000004009d2 <+114>:	call   0x400780 _ZNSolsEd@plt

   0x00000000004009d7 <+119>:	mov    esi,0x400830
   0x00000000004009dc <+124>:	mov    rdi,rax
   0x00000000004009df <+127>:	call   0x400820 _ZNSolsEPFRSoS_E@plt


27	    std::cout << max('a', 'b') << std::endl;
   0x00000000004009e4 <+132>:	mov    esi,0x62
   0x00000000004009e9 <+137>:	mov    edi,0x61
   0x00000000004009ee <+142>:	call   0x400b26 char max<char>(char, char)

   0x00000000004009f3 <+147>:	movsx  eax,al
   0x00000000004009f6 <+150>:	mov    esi,eax
   0x00000000004009f8 <+152>:	mov    edi,0x602080
   0x00000000004009fd <+157>:	call   0x4007e0 _ZStlsISt11char_traitsIcEERSt13basic_ostreamIcT_ES5_c@plt

   0x0000000000400a02 <+162>:	mov    esi,0x400830
   0x0000000000400a07 <+167>:	mov    rdi,rax
   0x0000000000400a0a <+170>:	call   0x400820 _ZNSolsEPFRSoS_E@plt

28	
29	    char* x = "1234";
   0x0000000000400a0f <+175>:	mov    QWORD PTR [rbp-0x8],0x400c11

30	    char* y = "2345";
   0x0000000000400a17 <+183>:	mov    QWORD PTR [rbp-0x10],0x400c16

31	    std::cout << max(x, y) << std::endl;
   0x0000000000400a1f <+191>:	mov    rdx,QWORD PTR [rbp-0x10]
   0x0000000000400a23 <+195>:	mov    rax,QWORD PTR [rbp-0x8]
   0x0000000000400a27 <+199>:	mov    rsi,rdx
   0x0000000000400a2a <+202>:	mov    rdi,rax
   0x0000000000400a2d <+205>:	call   0x40092d char* max<char*>(char*, char*)

   0x0000000000400a32 <+210>:	mov    rsi,rax
   0x0000000000400a35 <+213>:	mov    edi,0x602080
   0x0000000000400a3a <+218>:	call   0x400800 _ZStlsISt11char_traitsIcEERSt13basic_ostreamIcT_ES5_PKc@plt

   0x0000000000400a3f <+223>:	mov    esi,0x400830
   0x0000000000400a44 <+228>:	mov    rdi,rax
   0x0000000000400a47 <+231>:	call   0x400820 _ZNSolsEPFRSoS_E@plt


32	
33	    std::cout << max(10, 12.8) << std::endl;
=> 0x0000000000400a4c <+236>:	movabs rax,0x402999999999999a
   0x0000000000400a56 <+246>:	mov    QWORD PTR [rbp-0x18],rax
   0x0000000000400a5a <+250>:	movsd  xmm0,QWORD PTR [rbp-0x18]
   0x0000000000400a5f <+255>:	mov    edi,0xa
   0x0000000000400a64 <+260>:	call   0x400b49 int max<int, double>(int, double)

   0x0000000000400a69 <+265>:	mov    esi,eax
   0x0000000000400a6b <+267>:	mov    edi,0x602080
   0x0000000000400a70 <+272>:	call   0x400790 _ZNSolsEi@plt

   0x0000000000400a75 <+277>:	mov    esi,0x400830
   0x0000000000400a7a <+282>:	mov    rdi,rax
   0x0000000000400a7d <+285>:	call   0x400820 _ZNSolsEPFRSoS_E@plt


34	    return 0;
   0x0000000000400a82 <+290>:	mov    eax,0x0

35	}
   0x0000000000400a87 <+295>:	leave  
   0x0000000000400a88 <+296>:	ret    

End of assembler dump.
```

泛型的优点是在编译期完成的，可以减少运行期的时间。  
```
#include <iostream>
using namespace std;
// 1+2+3...+100 ==> n*(n+1)/2 

template<int n>
struct Sum
{
	enum Value {N = Sum<n-1>::N+n}; // Sum(n) = Sum(n-1)+n
};
template<>
struct Sum<1>
{
	enum Value {N = 1};    // n=1
};

int main()
{
	cout << Sum<100>::N << endl;

    return 0;
}
```

汇编实现
```c++
Dump of assembler code for function main():
17	{
   0x00000000004007ad <+0>:	push   rbp
   0x00000000004007ae <+1>:	mov    rbp,rsp

18		cout << Sum<100>::N << endl;
=> 0x00000000004007b1 <+4>:	mov    esi,0x13ba    //编译时就已经计算完毕了， 5050  
   0x00000000004007b6 <+9>:	mov    edi,0x601060
   0x00000000004007bb <+14>:	call   0x400640 _ZNSolsEi@plt

   0x00000000004007c0 <+19>:	mov    esi,0x4006b0
   0x00000000004007c5 <+24>:	mov    rdi,rax
   0x00000000004007c8 <+27>:	call   0x4006a0 _ZNSolsEPFRSoS_E@plt

19	
20	    return 0;
   0x00000000004007cd <+32>:	mov    eax,0x0

21	}
   0x00000000004007d2 <+37>:	pop    rbp
   0x00000000004007d3 <+38>:	ret    

End of assembler dump.
```

## 进阶编程
### STL标准模板库(Standard Template Library) 

- STL算法是泛型的(generic), 不与任何特定的数据结构和对象绑定，不必在环境类似的环境下重写代码；    
- STL算法可以量身定做，并且具有很高的效率；  
- STL可以进行扩充，你可以编写自己的组件并且能与STL标准的组件进行很好地融合；

> STL六大组件  

![stl标准库](../../res/stl标准库.png)
<br>

![stl标准库组件之间的关系](../../res/stl-relationship.png)

### 容器  
容器用于存放数据; STL的容器分为两大类:
- 序列式容器(Sequence Containers):  
其中的元素都是可排序的(ordered),STL提供了vector, list, deque等序列式容器,而stack, queue, priority_ queue则是容器适配器;    
- 关联式容器(Associative Containers):  
每个数据元素都是由一-个键(key)和值(Value)组成，当元素被插入到容器时，按基键以某种特定规则放入适当位置;常见的STL关联容器如: set, multiset, map, multimap;

### 序列式容器的基本使用  
```c++
#include <vector>
#include <list>
#include <queue>
#include <stack>
#include <map>
#include <string>
#include <functional>
#include <algorithm>
#include <utility>
#include <iostream>
using namespace std;

struct Display
{
    void operator()(int i)
    {
        cout << i << " ";
    }
};

struct Display2
{
    void operator()(pair<string, double> info)
    {
        cout << info.first << ":  " << info.second << "  ";
    }
};

int main()
{
    int iArr[] = {1, 2, 3, 4, 5};

    vector<int> iVector(iArr, iArr + 4);
    list<int> iList(iArr, iArr + 4);
    deque<int> iDeque(iArr, iArr + 4);
    queue<int> iQueue(iDeque);                   // 队列 先进先出
    stack<int> iStack(iDeque);                   // 栈 先进后出
    priority_queue<int> iPQueue(iArr, iArr + 4); // 优先队列，按优先权

    for_each(iVector.begin(), iVector.end(), Display());
    cout << endl;
    for_each(iList.begin(), iList.end(), Display());
    cout << endl;
    for_each(iDeque.begin(), iDeque.end(), Display());
    cout << endl;

    while (!iQueue.empty())
    {
        cout << iQueue.front() << " "; // 1  2 3 4
        iQueue.pop();
    }
    cout << endl;

    while (!iStack.empty())
    {
        cout << iStack.top() << " "; // 4 3  2  1
        iStack.pop();
    }
    cout << endl;

    while (!iPQueue.empty())
    {
        cout << iPQueue.top() << " "; // 4 3 2 1
        iPQueue.pop();
    }
    cout << endl;

    return 0;
}
```

## GUI开发
## 陷阱与经验