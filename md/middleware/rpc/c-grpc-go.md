# c语言通过grpc与go通信  
[参考链接](https://github.com/ymm135/grpc-c)  

## 环境搭建  
### [protobuf](https://github.com/protocolbuffers/protobuf)    
[安装文档](https://github.com/protocolbuffers/protobuf/blob/bba446bbf2ac7b0b9923d4eb07d5acd0665a8cf0/src/README.md)  

依赖安装:  
```shell
# Ubuntu
sudo apt-get install autoconf automake libtool curl make g++ unzip

# Centos
sudo yum install autoconf automake libtool curl make gcc gcc-c++ unzip
```

> curl会下载gmock, 如果虚拟机无法访问，请手动下载:  
```
# autogen.sh
# Check that gmock is present.  Usually it is already there since the
# directory is set up as an SVN external.
if test ! -e gmock; then
  echo "Google Mock not present.  Fetching gmock-1.7.0 from the web..."
  curl $curlopts -L -O https://github.com/google/googlemock/archive/release-1.7.0.zip
  unzip -q release-1.7.0.zip
  rm release-1.7.0.zip
  mv googlemock-release-1.7.0 gmock

  curl $curlopts -L -O https://github.com/google/googletest/archive/release-1.7.0.zip
  unzip -q release-1.7.0.zip
  rm release-1.7.0.zip
  mv googletest-release-1.7.0 gmock/gtest
fi
```

protobuf安装
```shell
$ ./autogen.sh

$ ./configure
$ make
$ make check
$ sudo make install
$ sudo ldconfig # refresh shared library cache.
```

增加两个环境变量  
```
# 依赖检测
export PKG_CONFIG_PATH=$PKG_CONFIG_PATH:/usr/local/lib/pkgconfig

# 依赖库路径
export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/usr/local/lib
```

安装成功测试:  
```shell
$ pkg-config --cflags --libs protobuf
-pthread -I/usr/local/include  -pthread -L/usr/local/lib -lprotobuf -lpthread
```

### [protobuf-c](https://github.com/protobuf-c/protobuf-c)  

安装步骤:  
```
./configure && make && make install

# 如果git仓库
./autogen.sh && ./configure && make && make install
```

### [grpc](https://github.com/grpc/grpc)  
[安装文档](https://github.com/grpc/grpc/tree/master/src/cpp)  

安装依赖:  
```shell
# Ubuntu  
$ sudo apt-get install build-essential autoconf libtool pkg-config

# Centos
$ sudo yum install make automake gcc gcc-c++ kernel-devel autoconf libtool  pkgconfig.x86_64 
```

其中依赖很多第三方包`third_party`,需要下载依赖，在根目录执行`git submodule update --init`,不然会提示`cares/cares does not contain a CMakeLists.txt file.`  


使用cmake编译:  
```shell
# 下载依赖  
$ git submodule update --init

# 编译  
$ mkdir -p cmake/build
$ cd cmake/build
$ cmake ../..
$ make
$ sudo make install 
```

> 不同grpc需要指定cmake版本  

错误解决:  
- `grpc-c/third_party/grpc/third_party/zlib/zlib.h`:1758:44: 错误：`va_list`未声  

文件中增加头文件`#include <stdarg.h>`  


### [grpc-c](https://github.com/Juniper/grpc-c)   

编译:  
```shell
autoreconf --install
$ mkdir build && cd build
$ ../configure
$ make
$ sudo make install
```

安装的程序有:  
```shell
/usr/local/bin/protoc-gen-grpc-c
/usr/local/lib/libgrpc-c.so
```

构建example:  
```shell
cd build/examples
make gencode
make
```

出现错误:
```shell
[root@sd1 examples]# make gencode
--grpc-c_out: Unimplemented GenerateAll() method.
--grpc-c_out: Unimplemented GenerateAll() method.
--grpc-c_out: Unimplemented GenerateAll() method.
--grpc-c_out: Unimplemented GenerateAll() method.
make: *** [gencode] 错误 1
```

可以使用makefile(GUN Make)调试模式`SHELL="/bin/bash -vx"`  
```shell
[root@sd1 examples]# make SHELL="/bin/bash -vx" gencode 
(for protofile in `ls -1 ../../examples/*.proto` ; do \
    protoc -I ../../examples --grpc-c_out=. \
    --plugin=protoc-gen-grpc-c=../compiler/protoc-gen-grpc-c \
    $protofile; \
done)

+ for protofile in '`ls -1 ../../examples/*.proto`'
+ protoc -I ../../examples --grpc-c_out=. --plugin=protoc-gen-grpc-c=../compiler/protoc-gen-grpc-c ../../examples/server_streaming.proto
--grpc-c_out: Unimplemented GenerateAll() method.
make: *** [gencode] 错误 1
```

这里会告诉你makefile执行的语句是什么，相当于shell调试信息，可以看到错误是`protoc`  
```shell
protoc -I .. --grpc-c_out=. --plugin=protoc-gen-grpc-c=.. 
```

提示`--grpc-c_out`未实现`GenerateAll()`，通过`protoc --help`查看确实没有，有`--cpp_out=OUT_DIR`   
grpc的代码生成器目前有两个虚函数`Generate()`和`GenerateAll()`,目前版本需要实现`GenerateAll()`,应该和grpc版本有关系? 把/usr/local/lib和/usr/local/lib64下所有与protobuf相关的库都删除，重新install [参考链接](https://github.com/grpc/grpc/issues/10941)  

另外需要增加一些库的引用及升级 [openssl](https://github.com/openssl/openssl) 版本,需要在make时指定依赖库,但是在链接openssl会出现问题:    

```shell
[root@sd1 examples]# /bin/sh ../libtool  --tag=CC   --mode=link gcc -I. -I../../examples/../lib/h/ -I../../examples/../third_party/protobuf-c -I../../examples/../third_party/grpc/include -g -O2   -o foo_client foo_client.o foo.grpc-c.o ../lib/libgrpc-c.la -lgrpc -lgpr -lprotobuf-c -lpthread -lz -lcares -lm -lssl -lcrypto 
libtool: link: gcc -I. -I../../examples/../lib/h/ -I../../examples/../third_party/protobuf-c -I../../examples/../third_party/grpc/include -g -O2 -o .libs/foo_client foo_client.o foo.grpc-c.o  ../lib/.libs/libgrpc-c.so -lgrpc -lgpr /usr/local/lib/libprotobuf-c.so -lpthread -lz /usr/local/lib/libcares.so -lm -lssl -lcrypto
//usr/local/lib64/libgrpc.a(ssl_transport_security.c.o): In function `init_openssl':
ssl_transport_security.c:(.text+0x8d): undefined reference to `OpenSSL_add_all_algorithms'
//usr/local/lib64/libgrpc.a(ssl_transport_security.c.o): In function `add_pem_certificate':
ssl_transport_security.c:(.text+0x65e): undefined reference to `BIO_get_mem_data'
//usr/local/lib64/libgrpc.a(ssl_transport_security.c.o): In function `ssl_ctx_use_certificate_chain':
...
collect2: error: ld returned 1 exit status
```

这里出现问题是编译`/usr/local/lib64/libgrpc.a`时无法找到`OpenSSL_add_all_algorithms`定义，而不是使用`libgrpc.so`动态库，如果是动态库找不到，那说明是`libgrpc.so`编译有问题。  


### [protobuf go](https://github.com/golang/protobuf) 
如果使用protobuf 3.0.x的版本，没有内置gen-go插件，需要自定编译并安装  













