# ubuntu20 搭建内核调试环境  

搭建 linux 内核的调试环境
- 编译内核  
- 制作文件系统及启动程序 
- qemu 模拟器运行linux 
- 通过gdb或vscode调试内核源码

整体环境: macos + pd(vmware) + ubuntu + gdb + qemu + linux kernel 

## 安装ubuntu系统  
> 目前需要安装桌面版本，服务器版本不能启动虚拟机界面。  

下载镜像，安装系统。 我这里使用的是 `parallels desktop`,其他虚拟机也行。  

## 配置基础环境  

```shell
# 设置 root 密码。
sudo passwd root
# 切换 root 用户。
su root 

# 安装部分工具。
apt-get install vim git tmux openssh-server -y

vi /etc/ssh/sshd_config

# root用户远程登录。
PermitRootLogin yes 

# 重启sshd 
service sshd restart 

```

## 下载编译 linux 内核 
[github内核仓库]
```
# 下载源码
mkdir /root/work
wget https://github.com/torvalds/linux/archive/refs/tags/v5.13.tar.gz -O linux-5.13.tar.gz
tar -zxvf linux-5.13.tar.gz
cd linux-5.13 

# 安装编译依赖组件。
apt install build-essential flex bison libssl-dev libelf-dev libncurses-dev -y

# 设置调试的编译菜单。
make x86_64_defconfig
make menuconfig

# 下面选项如果没有选上的，选上（点击空格键），然后 save 保存设置，退出 exit。
##################################################################
Kernel hacking  --->
    Compile-time checks and compiler options  ---> 
        [*] Compile the kernel with debug info
            [*] Provide GDB scripts for kernel debugging
##################################################################

# 编译内核。
make -j4
```

## 源码安装 gdb
> 为什么要源码安装呢？后面使用qemu启动内核时，通过gdb调试时，会有一个错误，我们需要屏蔽它。  

<br>
<div align=center>
    <img src="../../res/kernel-debug2.png" width="40%" height="40%"></img>  
</div>
<br>

源码安装高版本的 gdb 8.3.1 

```shell
# 如果已安装gdb，请删除 gdb
gdb -v | grep gdb
apt remove gdb -y

# 下载解压 gdb
cd /root/work
wget http://ftp.gnu.org/gnu/gdb/gdb-8.3.1.tar.gz
tar -zxvf gdb-8.3.1.tar.gz

# 修改 gdb/remote.c 代码。
cd gdb-8.3.1
vim gdb/remote.c

```

注释一部分代码: 
```
/* Further sanity checks, with knowledge of the architecture.  */
// if (buf_len > 2 * rsa->sizeof_g_packet)
//   error (_("Remote 'g' packet reply is too long (expected %ld bytes, got %d "
//      "bytes): %s"),
//    rsa->sizeof_g_packet, buf_len / 2,
//    rs->buf.data ());

// 其他代码不动
  /* Save the size of the packet sent to us by the target.  It is used
     as a heuristic when determining the max size of packets that the
     target can safely receive.  */
  if (rsa->actual_register_packet_size == 0)
    rsa->actual_register_packet_size = buf_len;  
    ...
```

编译安装
```shell
./configure
make -j4
make install

# 校验一下版本对不对 
gcc -v
```

## gdb 调试内核 
因为linux内核是运行在虚拟中的，需要通过gdb远程调试。  

```shell
# 安装 qemu 模拟器，以及相关组件。 
qemu qemu-utils qemu-kvm virt-manager libvirt-daemon-system libvirt-clients bridge-utils -y

# 虚拟机进入 linux 内核源码目录。
cd /root/work/linux-5.13  

# 从 github 下载内核测试源码, 如果虚拟机无法访问github，下载离线包即可  
# git clone https://github.com/ymm135/kernel_test.git
wget https://github.com/ymm135/kernel_test/archive/refs/tags/v1.0.tar.gz -O kernel_test.tar.gz
# tar -zxvf kernel_test.tar.gz
# mv kernel_test

# 进入测试源码目录。
cd kernel_test/test_epoll_thundering_herd
# make 编译
make
# 通过 qemu 启动内核测试用例。
make rootfs
# 在 qemu 窗口输入小写字符 's', 启动测试用例服务程序。
s
# 在 qemu 窗口输入小写字符 'c', 启动测试用例客户端程序。
c
```

<br>
<div align=center>
    <img src="../../res/kernel-debug4.png" width="60%" height="60%"></img>  
</div>
<br>


```shell
# 通过 qemu 命令启动内核测试用例进行调试。
qemu-system-x86_64 -kernel ../../arch/x86/boot/bzImage -initrd ../rootfs.img -append nokaslr -S -s
# 在 qemu 窗口输入小写字符 's', 启动测试用例服务程序。
s
# 在 qemu 窗口输入小写字符 'c', 启动测试用例客户端程序。
c
```  

> -kernel bzImage use 'bzImage' as kernel image  
> -append cmdline use 'cmdline' as kernel command line  
> -initrd file    use 'file' as initial ram disk  



> 界面会出现 `guest has not initialized the display(yet)` ，这个不影响，因为我们增加`-S`参数，所以启动时暂停了  

<br>
<div align=center>
    <img src="../../res/kernel-debug1.png" width="60%" height="60%"></img>  
</div>
<br>

```
# gdb 调试命令。
gdb vmlinux
target remote :1234
b start_kernel
b tcp_v4_connect
c
focus
bt
```

<br>
<div align=center>
    <img src="../../res/kernel-debug5.png" width="100%" height="100%"></img>  
</div>
<br>

> 有时 `c` 会提示`The program is not being run`, 但是使用vscode远程调试就不会，可以直接进入IDE调试  

## vscode 远程调试  
首选vscode需要安装`remote-ssh`插件，配置好链接: 
`10.211.55.13` 虚拟机IP地址  
```shell
Host 10.211.55.13
HostName 10.211.55.13
User root
```

安装调试插件`C/C++ Extension Pack`, 打开ubuntu虚拟机目录: `/root/work/linux-5.13`

写好配置文件: 
```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "gdb内核启动",
            "type": "cppdbg",
            "request": "launch",
            "miDebuggerServerAddress": "127.0.0.1:1234",
            "program": "${workspaceFolder}/vmlinux",
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
                },
                {
                    "description":  "将反汇编风格设置为 Intel",
                    "text": "-gdb-set disassembly-flavor intel",
                    "ignoreFailures": true
                }
            ]
        }
    ]
}
```

增加好断点: `init/main.c:876 start_kernel函数`和`net/ipv4/tcp_ipv4.c:199 tcp_v4_connect函数`  
然后在qemu启动的内核窗口中输入`s`和`c`启动服务端和客户端   

<br>
<div align=center>
    <img src="../../res/kernel-debug3.png" width="100%" height="100%"></img>  
</div>
<br>

> 输入`s`和`c`启动服务端和客户端的含义  
`test_epoll_thundering_herd/main.c`中，最终被打包为: `rootfs.img`
```c
int main(int argc, char **argv) {
    char buf[64] = {0};
    int port = SERVER_PORT;
    const char *ip = SERVER_IP;

    if (argc >= 3) {
        ip = argv[1];
        port = atoi(argv[2]);
    }

    LOG("pls input 's' to run server or 'c' to run client!");

    while (1) {
        scanf("%s", buf);

        if (strcmp(buf, "s") == 0) {
            proc(ip, port);
        } else if (strcmp(buf, "c") == 0) {
            proc_client(ip, port, SEND_DATA);
        } else {
            LOG("pls input 's' to run server or 'c' to run client!");
        }
    }

    return 0;
}
```

把自己的程序编译成`init`二进制文件，然后打包为`rootfs.img`,最终内核运行时，会启动该文件  
```shell
$(CC) $(CFLAGS) $(INC) $(SRCS) -o init -static -lpthread # 把程序编译为init
find init | cpio -o -Hnewc | gzip -9 > ../rootfs.img     # 把init打包为img
```

> cpio -o或--create 　执行copy-out模式，建立备份档, -H<备份格式> 　指定备份时欲使用的文件格式   -Hnewc  SVR4的格式，如果使用ASCII  -H odc
> gzip -<压缩效率> 　压缩效率是一个介于1－9的数值，预设值为"6"，指定愈大的数值，压缩效率就会愈高。  

如果解压`rootfs.img`
```shell
mv rootfs.img rootfs.img.gz 
gzip -d -v rootfs.img.gz

cpio -i < rootfs.img  

# file init 
init: ELF 64-bit LSB executable, x86-64, version 1 (GNU/Linux), statically linked, BuildID[sha1]=88728f3ec8d24dc02b8aaadbaf4a66da499c7a0a, for GNU/Linux 3.2.0, not stripped
```



> vscode 全局配置文件json 调出命令行，输入`Preferences: Configure language specific settings` 可直接编辑setting.json配置文件  

> kernel镜像格式:vmlinux  vmlinuz是可引导的、可压缩的内核镜像，vm代表Virtual Memory.Linux支持虚拟内存，因此得名vm.它是由用户对内核源码编译得到，实质是elf格式的文件.也就是说，vmlinux是编译出来的最原始的内核文件，未压缩.这种格式的镜像文件多存放在PC机上.  
> kernel镜像格式:bzImage  bz表示big zImage,其格式与zImage类似，但采用了不同的压缩算法，注意，bzImage的压缩率更高.  


## 制作文件系统  

- ### busybox编译  
BusyBox 将许多常见 UNIX 实用程序的微小版本组合成一个小型可执行文件。BusyBox 为任何小型或嵌入式系统提供了一个相当完整的环境。  

```shell
$ wget http://busybox.net/downloads/busybox-1.32.0.tar.bz2
$ tar vjxf busybox-1.32.0.tar.bz2
$ cd busybox-1.32.0/
$ make menuconfig
```

相关的配置如下，同样[*]的代表打开，[ ]代表关闭.可以适当增加删减一些配置.  

```shell
Busybox Settings  --->
  [*] Don't use /usr (NEW)
  --- Build Options
  [*] Build BusyBox as a static binary (no shared libs)
  --- Installation Options ("make install" behavior)
  (./_install) BusyBox installation prefix (NEW)

Miscellaneous Utilities  --->
  [ ] flash_eraseall
  [ ] flash_lock
  [ ] flash_unlock
  [ ] flashcp
```

编译busybox，会安装到./_install目录下  

```shell
$ make -j4
$ make install

$ ls _install/
drwxr-xr-x 2 root root 4096 Apr  5 20:31 bin
lrwxrwxrwx 1 root root   11 Apr  5 20:31 linuxrc -> bin/busybox
drwxr-xr-x 2 root root 4096 Apr  5 20:31 sbin
``` 

- ### 制作文件系统  
根文件系统镜像大小256MiB，格式化为ext3文件系统．  

```shell
# in working-dir
$ dd if=/dev/zero of=rootfs.img bs=1M count=256
$ mkfs.ext3 rootfs.img
```

将文件系统mount到本地路径，复制busybox相关的文件，并生成必要的文件和目录
```shell
# in working-dir
$ mkdir /tmp/rootfs-busybox
$ sudo mount -o loop $PWD/rootfs.img /tmp/rootfs-busybox

$ sudo cp -a ../busybox-1.32.0/_install/* /tmp/rootfs-busybox/
$ pushd /tmp/rootfs-busybox/              # 进入到目录: /tmp/rootfs-busybox/
$ sudo mkdir dev sys proc etc lib mnt
$ popd                                    # 返回到源目录 
```

还需要制作系统初始化文件
```shell
# in working-dir
$ sudo cp -a ../busybox-1.32.0/examples/bootfloppy/etc/* /tmp/rootfs-busybox/etc/
```

Busybox所使用的rcS，内容可以写成  
```shell
$ cat /tmp/rootfs-busybox/etc/init.d/rcS
#! /bin/sh

/bin/mount -a
/bin/mount -t sysfs sysfs /sys
/bin/mount -t tmpfs tmpfs /dev
/sbin/mdev -s
```

可以查看文件列表: 
```shell
# ls -l /tmp/rootfs-busybox
drwxr-xr-x 2 root root  4096 Apr  5 20:31 bin
                                          ├── cat -> busybox
drwxr-xr-x 2 root root  4096 Apr  5 20:37 dev
drwxr-xr-x 3 root root  4096 Apr  5 20:39 etc
                                          ├── fstab
                                          ├── init.d
                                          │   └── rcS
                                          ├── inittab
                                          └── profile
drwxr-xr-x 2 root root  4096 Apr  5 20:37 lib
lrwxrwxrwx 1 root root    11 Apr  5 20:31 linuxrc -> bin/busybox
drwx------ 2 root root 16384 Apr  5 20:36 lost+found
drwxr-xr-x 2 root root  4096 Apr  5 20:37 mnt
drwxr-xr-x 2 root root  4096 Apr  5 20:37 proc
drwxr-xr-x 2 root root  4096 Apr  5 20:31 sbin
drwxr-xr-x 2 root root  4096 Apr  5 20:37 sys
```

接下来就不需要挂载的虚拟磁盘了

```shell
$ sudo umount /tmp/rootfs-busybox
```

通过qeum启动 
```shell
sudo qemu-system-x86_64 -kernel linux/arch/x86/boot/bzImage -append 'root=/dev/sda' -boot c -hda rootfs.img -k en-us 
```

> Block device options:  
> -fda/-fdb file  use 'file' as floppy disk 0/1 image  
> -hda/-hdb file  use 'file' as IDE hard disk 0/1 image  
> -hdc/-hdd file  use 'file' as IDE hard disk 2/3 image  
> -cdrom file     use 'file' as IDE cdrom image (cdrom is ide1 master)  

<br>
<div align=center>
    <img src="../../res/kernel-debug2.png" width="40%" height="40%"></img>  
</div>
<br>

> 使用文件系统启动后，会停止kernerl打印，输入指令   
> 使用ctrl+alt+2切换qemu控制台,使用ctrl+alt+1切换回调试kernel  


- ### 向文件系统增加程序并调试  

```shell
# 挂载
$ mkdir /tmp/rootfs-busybox
$ sudo mount -o loop $PWD/rootfs.img /tmp/rootfs-busybox
# 向文件系统增加自己的程序 work/program 
$ pushd /tmp/rootfs-busybox/
$ mkdir work && cd work 
$ cp /path/to/program  .    # /root/work/linux-5.13/kernel_test/test_unix_socket/test_unix_socket  测试程序  
$ popd  

# 卸载
$ sudo umount /tmp/rootfs-busybox  
```


## 调试网络

### 网络设备设置

我们使用qemu的**bridge模式**设置虚机网络，该模式需要在宿主机配置网桥，并使用该网桥配置地址和默认路由．具体见host的`/etc/qemu-ifup`文件．
然后使用`-net`参数启动qemu虚机．

> 如果通过宿主机eth0远程登录，该操作可能导致网络登录中断． 如果使用虚拟机，可以多开几个虚拟网卡，使用非远程登录的接口    

```shell
$ sudo brctl addbr br0
$ sudo brctl addif br0 enp0s6
$ sudo ifconfig enp0s6 0
$ sudo dhclient br0
$ sudo qemu-system-x86_64 -kernel linux/arch/x86/boot/bzImage \
       -append 'root=/dev/sda' -boot c -hda rootfs.img -k en-us \
       -net nic -net tap,ifname=tap0  
```

`tap0`是在宿主机中对应的接口名．我们可以在宿主机中看到网桥及其两个端口．tap设备的另一端是VM的eth0．

```shell
$ brctl show
bridge name	bridge id		STP enabled	interfaces
br0		8000.001c42b3acf8	no		enp0s6
virbr0		8000.5254000ab2bb	yes		virbr0-nic
```

为简单测试VM和宿主机的网络联通性，我们在宿主机的`br0`和虚机的`eth0`分别配置两个私有地址来测试．

```shell
# qemu VM（调试Kernel）
/ # ip addr add 192.168.0.2/24 dev eth0 
/ # ip link set eth0 up
```

```shell
# 宿主机
$ sudo ip addr add 192.168.0.1/24 dev br0
$ ping 192.168.0.2
PING 192.168.0.2 (192.168.0.2) 56(84) bytes of data.
64 bytes from 192.168.0.2: icmp_seq=1 ttl=64 time=0.328 ms
64 bytes from 192.168.0.2: icmp_seq=2 ttl=64 time=0.282 ms
... ...
```

```shell
# qemu VM（调试Kernel）
$ ping 192.168.0.1
PING 192.168.0.2 (192.168.0.1) 56(84) bytes of data.
64 bytes from 192.168.0.2: icmp_seq=1 ttl=64 time=0.328 ms
64 bytes from 192.168.0.2: icmp_seq=2 ttl=64 time=0.282 ms
... ...
```

这样可以基本满足调试Kernel的网络协议栈的环境了.

> 或者在VM中运行`udhcpc eth0`让VM获取和host相同网络的IP，不过这需要DHCP Server的支持．

## 使用nfs挂载rootfs 

> NFS 是Network File System的缩写，即网络文件系统。一种使用于分散式文件系统的协定，由Sun公司开发，于1984年向外公布。功能是通过网络让不同的机器、不同的操作系统能够彼此分享个别的数据，让应用程序在客户端通过网络访问位于服务器磁盘中的数据，是在类Unix系统间实现磁盘文件共享的一种方法。  

调试内核模块（或其他用户态程序的时候），挂载静态的ext3文件系统并不方便．为此我们可以采用nfs的形式挂载qemu kernel的rootfs，这样就能方便的在host中修改，编译内核模块，并在qemu kernel中配合gdb进行调试．根文件系统制作方法和之前相同，  

```
# working dir
$ mkdir rootfs.nfs
$ cp -a ../busybox-1.32.0/_install/* rootfs.nfs/
$ pushd rootfs.nfs/
$ mkdir dev sys proc etc lib mnt
$ popd
$ cp -a ../busybox-1.32.0/examples/bootfloppy/etc/* rootfs.nfs/etc/
$ cat rootfs.nfs/etc/init.d/rcS   # 修改文件如下: 
#! /bin/sh

/bin/mount -a
/bin/mount -t sysfs sysfs /sys
/bin/mount -t tmpfs tmpfs /dev
/sbin/mdev -s

$ chmod -R 777 rootfs.nfs/
```

配置host的nfs服务并启动,

```shell
$ sudo apt-get install nfs-kernel-server
$ cat /etc/exports   # 修改如下
/path/to/working/dir/rootfs.nfs *(rw,insecure,sync,no_root_squash)

$ service nfs-kernel-server restart
```

使用nfs挂载qemu Kernel的根文件系统

```
$ sudo qemu-system-x86_64 -kernel linux/arch/x86/boot/bzImage \
        -append 'root=/dev/nfs nfsroot="192.168.1.1:/path/to/working/dir/rootfs.nfs/" rw ip=192.168.1.2' \
	-boot c -k en-us -net nic -net tap,ifname=tap0
```

其中`nfsroot`为host的IP及要挂载的根文件系统在host中的路径，`ip`参数填写qemu Kernel将使用的IP地址．

### 调试内核模块

> Note: 编译内核模块的时候，源码树和虚机Kernel编译需要是同一份．不然会出现模块版本不匹配无法运行的情况. 编译内核模块的时候，使用`ccflags-y += -g -O0`保留信息避免优化．

Kernel模块每次插入后的内存位置不确定，需要将其各个内存段的位置取出才能按源码单步调试． 首先在`do_init_module`设置断点, insmod的时候会触发断点，

```gdb
(gdb) b do_init_module
```

模块各内存段信息保存在`mod->sect_attrs->attrs[]`数组中，我们需要以下几个字段信息,

* `.text`
* `.rodata`
* `.bss`

分别打印字段的名字和其地址，

```gdb
(gdb) print  mod->sect_attrs->attrs[1].name
$82 = 0xffff880006109ad8 <__this_cpu_preempt_check> ".text"
(gdb) print  mod->sect_attrs->attrs[5].name
$86 = 0xffff880006109ad0 <__phys_addr_nodebug+10> ".rodata"
(gdb) print  mod->sect_attrs->attrs[12].name
$93 = 0xffff880006109ac8 <__phys_addr_nodebug+2> ".bss"

(gdb) print /x  mod->sect_attrs->attrs[1]->address
$96 = 0xffffffffa0005000
(gdb) print /x  mod->sect_attrs->attrs[5]->address
$97 = 0xffffffffa0006040
(gdb) print /x  mod->sect_attrs->attrs[12]->address
$98 = 0xffffffffa0007380
```

然后为gdb设置模块路径和各内存段地址，

```gdb
(gdb) add-symbol-file /path/to/module/xxx.ko 0xffffffffa0005000 \
-s .data 0xffffffffa0006040 \
-s .bss 0xffffffffa0007380
```

接下来，就能为模块的各个函数设置断点进行调试了．


- #### 参考文章1 [搭建 Linux 内核网络调试环境](https://zhuanlan.zhihu.com/p/445453676)  
- #### 参考文章2 [使用 GDB + Qemu 调试 Linux 内核](https://z.itpub.net/article/detail/9CCD29B78F55B5BEA664AD7045915411)  
- #### [linux内核其他调试环境](../../md/other/linux-core-debug.md) 
- #### [gdb-kernel-debugging](https://www.kernel.org/doc/html/v4.11/dev-tools/gdb-kernel-debugging.html)
- #### [kernel-qemu-gdb](https://github.com/beacer/notes/blob/master/kernel/kernel-qemu-gdb.md)  