# Socket通信  
[维基百科](https://zh.wikipedia.org/wiki/%E7%B6%B2%E8%B7%AF%E6%8F%92%E5%BA%A7)  
网络套接字（英语：Network socket；又译网络套接字、网络接口、网络插槽）在计算机科学中是电脑网络中进程间资料流的端点。使用以网际协议（Internet Protocol）为通信基础的网络套接字，称为网际套接字（Internet socket）。因为网际协议的流行，现代绝大多数的网络套接字，都是属于网际套接字。

socket是一种操作系统提供的进程间通信机制。  

在操作系统中，通常会为应用程序提供一组应用程序接口（API），称为套接字接口（英语：socket API）。应用程序可以通过套接字接口，来使用网络套接字，以进行资料交换。最早的套接字接口来自于4.2 BSD，因此现代常见的套接字接口大多源自Berkeley套接字（Berkeley sockets）标准。在套接字接口中，以IP地址及端口组成套接字地址（socket address）。远程的套接字地址，以及本地的套接字地址完成连线后，再加上使用的协议（protocol），这个五元组（five-element tuple），作为套接字对（socket pairs），之后就可以彼此交换资料。例如，在同一台计算机上，TCP协议与UDP协议可以同时使用相同的port而互不干扰。 操作系统根据套接字地址，可以决定应该将资料送达特定的行程或线程。这就像是电话系统中，以电话号码加上分机号码，来决定通话对象一般。  

## [python socket](https://docs.python.org/3/library/socket.html)   
不同socket 协议簇地址结构是不同的:  

- ### AF_UNIX
一个绑定在文件系统节点上的 AF_UNIX 套接字的地址表示为一个字符串，使用文件系统字符编码和 'surrogateescape' 错误回调方法（see PEP 383）。一个地址在 Linux 的抽象命名空间被返回为带有初始的 null 字节的 bytes-like object ；注意在这个命名空间中的套接字可能与普通文件系统套接字通信，所以打算运行在 Linux 上的程序可能需要解决两种地址类型。当传递为参数时，一个字符串或字节类对象可以用于任一类型的地址。


- ### AF_INET 
一对 (host, port) 被用于 AF_INET 地址族，host 是一个代表互联网域名表示法之内主机名或者一个 IPv4 地址的字符串，例如 'daring.cwi.nl' 或 '100.50.200.5'，port 是一个整数。

对于 IPv4 地址，有两种可接受的特殊形式被用来代替一个主机地址： '' 代表 INADDR_ANY，用来绑定到所有接口；字符串 '<broadcast>' 代表 INADDR_BROADCAST。此行为不兼容 IPv6，因此，如果你的 Python 程序打算支持 IPv6，则可能需要避开这些。

- ### AF_INET6  
对于 AF_INET6 地址族，使用一个四元组 (host, port, flowinfo, scopeid)， flowinfo 和 scopeid 代表了 C 库里 struct sockaddr_in6 中的 sin6_flowinfo 和 sin6_scope_id 成员。 对于 socket 模块中的方法， flowinfo 和 scopeid 可以被省略，只为了向后兼容。注意，scopeid 的省略可能会导致 IPv6 地址的操作范围问题。  

- ### AF_NETLINK  
AF_NETLINK 套接字由一对 (pid, groups) 表示。  

- ### AF_CAN  
AF_CAN 地址族使用元组 (interface, )，其中 interface 是表示网络接口名称的字符串，如 'can0'。网络接口名 '' 可以用于接收本族所有网络接口的数据包。  

- ### AF_ALG 
AF_ALG 是一个仅 Linux 可用的、基于套接字的接口，用于连接内核加密算法。算法套接字可用包括 2 至 4 个元素的元组来配置 (type, name [, feat [, mask]])

- ### AF_PACKET 
AF_PACKET 是一个底层接口，直接连接至网卡。数据包使用元组 (ifname, proto[, pkttype[, hatype[, addr]]]) 表示，其中：  
- ifname - 指定设备名称的字符串。
- proto - 一个用网络字节序表示的整数，指定以太网协议版本号。
- pkttype - 指定数据包类型的整数（可选）：
- PACKET_HOST （默认） - 数据包寻址到本地宿主机。
- PACKET_BROADCAST - 物理层广播报文。
- PACKET_MULTIHOST - 数据包发送到物理层多播地址。
- PACKET_OTHERHOST - 被（处于混杂模式的）网卡驱动捕获的、发送到其他主机的数据包。
- PACKET_OUTGOING - 来自本地主机的、回环到一个套接字的数据包。
- hatype - 可选整数，指定 ARP 硬件地址类型。
- addr - 可选的类字节串对象，用于指定硬件物理地址，其解释取决于各设备。  

### AF_UNIX 实现原理 
- unix - sockets for local interprocess communication  

```c
#include <sys/socket.h>
#include <sys/un.h>

unix_socket = socket(AF_UNIX, type, 0);
error = socketpair(AF_UNIX, type, 0, int *sv);
```  

AF_UNIX（也称为 AF_LOCAL ）套接字用于同一台机器上的进程之间进行通信。
传统上，UNIX 域套接字可以是未命名的，或者绑定到文件系统路径名（标记为套接字类型）。UNIX 域中的有效套接字类型是：
- SOCK_STREAM 有保障的(即能保证数据正确传送到对方)面向连接的SOCKET，多用于资料(如文件)传送  
- SOCK_DGRAM 是无保障的面向消息的socket，主要用于在网络上发广播信息  

[Demo实例](https://docs.oracle.com/cd/E19504-01/802-5886/6i9k5sgsl/index.html)  

### AF_PACKET 实现原理 

[内核调试方法](../../../md/other/linux-core-debug.md)  
















