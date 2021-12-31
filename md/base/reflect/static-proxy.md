# 静态代理
[静态代理代码](../../../code/reflect/proxy/main.go)  

```go
package main

import "fmt"

type IMan interface {
	Walk() int
}

type ManProxy struct {
	manProxy IMan // 不能是指针
}

func (proxy *ManProxy) setProxy(man IMan) () {
	proxy.manProxy = man
}

func (proxy *ManProxy) Walk() int {
	proxy.manProxy.Walk()
	fmt.Println("proxy Walk")
	return 0
}

type Man struct {
}

func (man *Man) Walk() int {
	fmt.Println("man Walk")
	return 0
}

func main() {
	manProxy := &ManProxy{} // 是不是指针都行
	var manImpl IMan
	var man Man
	manImpl = &man //需要取地址, 调用man具体实现,而不是复制

	manProxy.setProxy(manImpl)
	manProxy.Walk()

	return
}
```  

> 疑问: 为什么需要实现类的地址?(&man)以及`manImpl = &man`的原理是啥?  

这里应该是编译器的解析规则，比如接口包含实现类的指针，我猜`manImpl = &man`语句就是把man结构实现的地址
赋给`IMan`的指针  






