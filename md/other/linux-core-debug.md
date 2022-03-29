# Linux 内核调试  

## 内核编译  

centos7 内核版本  
```shell
$ uname -a 
Linux d1.localdomain 3.10.0-1127.el7.x86_64 #1 SMP Tue Mar 31 23:36:51 UTC 2020 x86_64 x86_64 x86_64 GNU/Linux  
```


- ### 安装环境: 
#### centos  
```shell
yum install ncurses-devel bison flex elfutils-libelf-devel openssl-devel
``` 

下载源码:  
内核版本为: `3.10.0-1127.el7.x86_64`  
[下载地址 git clone https://github.com/torvalds/linux.git ](https://github.com/torvalds/linux) ,也可以使用`sudo apt-get source linux-image-$(uname -r)`下载当前内核版本或更小的发行版，缺点：版本不全  

```
git clone https://github.com/torvalds/linux.git
git tag | grep 3.10.0
git checkout 
```

更新GCC版本，最低5.1.0,这里升级到gcc7  
```shell
sudo yum install centos-release-scl
sudo yum install devtoolset-7-gcc   # devtoolset-8-gcc 

# 切换对应版本  
scl enable devtoolset-7 bash  

# 测试 
gcc -v

# 直接替换旧的gcc（终极解决方案）
mv /usr/bin/gcc /usr/bin/gcc-4.8.5
ln -s /opt/rh/devtoolset-7/root/bin/gcc /usr/bin/gcc
mv /usr/bin/g++ /usr/bin/g++-4.8.5
ln -s /opt/rh/devtoolset-7/root/bin/g++ /usr/bin/g++
gcc --version
g++ --version
```

> 也可以 [下载源码](https://ftp.gnu.org/gnu/gcc/) 编译，比较麻烦  

#### ubuntu 
安装依赖
```
apt install 
```

gcc9编译有问题，更换为gcc8 
```
# 安装
sudo apt install gcc-8 g++-8  

# 切换为gcc8 
sudo update-alternatives --install /usr/bin/gcc gcc /usr/bin/gcc-8 100
sudo update-alternatives --install /usr/bin/g++ g++ /usr/bin/g++-8 100

# 多版本切换
sudo update-alternatives --config gcc
sudo update-alternatives --config g++
``` 

> 找不到库 `libpixman-1.so.0`，可以使用 `apt-cache search pixman` 查找  



- ### 编译内核  
```
make x86_64_defconfig # 选择对应平台
make menuconfig
make -j8
```

> 要进行打断点调试，需要关闭系统的随机化和开启调试信息 
```shell
Processor type and features  ---> 
    [ ] Build a relocatable kernel                                               
        [ ]  Randomize the address of the kernel image (KASLR) (NEW)  # 按键N关闭 


Kernel hacking  --->
    Compile-time checks and compiler options  --->  
        [*] Compile the kernel with debug info                                                                  
        [ ]   Reduce debugging information                                                                      
        [ ]   Produce split debuginfo in .dwo files                                                             
        [*]   Generate dwarf4 debuginfo                                         
        [*]   Provide GDB scripts for kernel debugging  
```

修改配置会保存在.config中，可以自行查看  
```shell
# grep CONFIG_DEBUG_INFO .config

CONFIG_DEBUG_INFO=y
# CONFIG_DEBUG_INFO_NONE is not set
CONFIG_DEBUG_INFO_DWARF_TOOLCHAIN_DEFAULT=y
# CONFIG_DEBUG_INFO_DWARF4 is not set
# CONFIG_DEBUG_INFO_DWARF5 is not set
# CONFIG_DEBUG_INFO_REDUCED is not set
# CONFIG_DEBUG_INFO_COMPRESSED is not set
# CONFIG_DEBUG_INFO_SPLIT is not set
```  

> 配置文件和之前有差异  


编译成功结果:  
```shell
  OBJCOPY arch/x86/boot/setup.bin
  BUILD   arch/x86/boot/bzImage
Kernel: arch/x86/boot/bzImage is ready  (#1)
```

查看编译完成后的文件 
```shell
# 未压缩的内核文件，这个在 gdb 的时候需要加载，用于读取 symbol 符号信息，由于包含调试信息所以比较大
 $ ls -hl vmlinux 
-rwxr-xr-x. 1 root root 348M 3月  28 20:56 vmlinux

# 压缩后的镜像文件 
$ ls -hl ./arch/x86_64/boot/bzImage
lrwxrwxrwx. 1 root root 22 3月  28 20:56 ./arch/x86_64/boot/bzImage -> ../../x86/boot/bzImage

$ ls -hl ./arch/x86/boot/bzImage
-rw-r--r--. 1 root root 9.7M 3月  28 20:56 ./arch/x86/boot/bzImage
```

- ### 启动内存文件系统制作  
[busybox](https://busybox.net/about.html)  

BusyBox 将许多常见 UNIX 实用程序的微小版本组合成一个小型可执行文件。它为您通常在 GNU fileutils、shellutils 等中找到的大多数实用程序提供了替代品。BusyBox 中的实用程序通常比它们功能齐全的 GNU 表亲具有更少的选项；但是，包含的选项提供了预期的功能，并且其行为与 GNU 对应项非常相似。BusyBox 为任何小型或嵌入式系统提供了一个相当完整的环境。  

```shell
# 首先安装静态依赖，否则会有报错，参见后续的排错章节
$ yum install -y glibc-static.x86_64 -y

$ wget https://busybox.net/downloads/busybox-1.32.1.tar.bz2
$ tar -xvf busybox-1.32.1.tar.bz2
$ cd busybox-1.32.1/

$ make menuconfig

# 安装完成后生成的相关文件会在 _install 目录下
$ make && make install   

$ cd _install
$ mkdir proc
$ mkdir sys
$ touch init  

#  init 内容见后续章节，为内核启动的初始化程序
$ vim init   

# 必须设置成可执行文件
$ chmod +x init  

$ find . | cpio -o --format=newc > ./rootfs.img
cpio: 文件 ./rootfs.img 增长，34304 新字节未被拷贝
2055 块 

$ ls -hl rootfs.img
-rw-r--r--. 1 root root 1.1M 3月  28 21:54 rootfs.img
```

其中上述的 `init` 文件内容如下，打印启动日志和系统的整个启动过程花费的时间  
```
#!/bin/sh
echo "{==DBG==} INIT SCRIPT"
mkdir /tmp
mount -t proc none /proc
mount -t sysfs none /sys
mount -t debugfs none /sys/kernel/debug
mount -t tmpfs none /tmp

mdev -s 
echo -e "{==DBG==} Boot took $(cut -d' ' -f1 /proc/uptime) seconds"

# normal user
setsid /bin/cttyhack setuidgid 1000 /bin/sh
```

到此为止我们已经编译了好了 Linux 内核（vmlinux 和 bzImage）和启动的内存文件系统（rootfs.img）  

> 如果在编译时，没有找到指定库，可以使用 `yum provides */libm.a` 类似语句查询   

#### **现在内核源码已经编译完成，可通过多种方式学习内核，如果仅仅了解每个功能的流程，可以通过vscode及插件阅读源码，另外可替换当前系统内核进行调试。如果想了解内核流程，可通过虚拟机启动内核进行调试。**

- ### vscode 查看内核源码  
通过vscode的 `Remote SSH` 连接虚拟机，打开编译好的内核源码目录。  

- ### macos 启动内核调试 

```shell
brew install qemu 
```

编译好的文件，直接启动即可:  
```shell
qemu-system-x86_64 -kernel bzImage -initrd rootfs.img
```


- ### 替换内核进行调试  


- ### Qemu 启动内核调试   

> 无图像启动: -nographic  

你可能使用的是ubuntu或者centos，都需要搭建 KVM (Kernel-based Virtual Machine).  

KVM 是基于 x86 虚拟化扩展(Intel VT 或者 AMD-V) 技术的虚拟机软件，所以查看 CPU 是否支持 VT 技术，就可以判断是否支持KVM。有返回结果，如果结果中有vmx（Intel）或svm(AMD)字样，就说明CPU的支持的。  

```shell
cat /proc/cpuinfo | egrep 'vmx|svm'

flags   : fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush dts acpi mmx fxsr sse sse2 ss ht tm pbe syscall nx pdpe1gb rdtscp lm constant_tsc arch_perfmon pebs bts rep_good nopl xtopology nonstop_tsc aperfmperf eagerfpu pni pclmulqdq dtes64 monitor ds_cpl vmx smx est tm2 ssse3 fma cx16 xtpr pdcm pcid dca sse4_1 sse4_2 x2apic movbe popcnt tsc_deadline_timer aes xsave avx f16c rdrand lahf_lm abm arat epb pln pts dtherm tpr_shadow vnmi flexpriority ept vpid fsgsbase tsc_adjust bmi1 avx2 smep bmi2 erms invpcid cqm xsaveopt cqm_llc cqm_occup_llc
```  

#### centos  
macos不支持，但是需要增加额外参数  
kvm is the linux hypervisor implementation, that isn't going to work. Recent qemu version have support for the macos hypervisor framework, use `accel=hvf` for that.  

```
qemu-system-x86_64 -m 2G -hda ubuntu.20.qcow2 -accel hvf
```

> ERROR    主机不支持 虚拟化类型 'hvm' 架构 'x86_64' 的虚拟机 kvm  


#### ubuntu 
需要在调试的 Ubuntu 20.04 的系统中安装 Qemu 工具，其中调测的 Ubuntu 系统使用 VirtualBox 安装。   

```shell
cp busybox-1.32.1/_install/rootfs.img .
cp linux//arch/x86/boot/bzImage .

# 启动
apt install qemu qemu-utils qemu-kvm virt-manager libvirt-daemon-system libvirt-clients bridge-utils
```

把上述编译好的 vmlinux、bzImage、rootfs.img 和编译的源码拷贝到我们当前 Unbuntu 机器中。

拷贝 Linux 编译的源码主要是在 gdb 的调试过程中查看源码，其中 vmlinux 和 linux 源码处于相同的目录，本例中 vmlinux 位于 linux-4.19.172 源目录中。  

```shell
qemu-system-x86_64 -kernel ./bzImage -initrd  ./rootfs.img -append "nokaslr console=ttyS0" -s -S -nographic
```

使用上述命令启动调试，启动后会停止在界面处，并等待远程 gdb 进行调试，在使用 GDB 调试之前，可以先使用以下命令进程测试内核启动是否正常。  

```shell
qemu-system-x86_64 -kernel ./bzImage -initrd  ./rootfs.img -append "nokaslr console=ttyS0" -nographic
```

> qemu: could not load PC BIOS 'bios-256k.bin' 

在安装后seabios，发现有`/usr/share/seabios/bios-256k.bin`, 需要增加路径 `-L /usr/share/seabios`  

> qemu-system-x86_64 -kernel ../../arch/x86/boot/bzImage -initrd ../rootfs.img  
> Unable to init server: Could not connect: Connection refused  
> gtk initialization failed  

不能在远程用命令行运行，需要在有图形界面的系统运行该指令。  



- ### GDB 调试内核  



