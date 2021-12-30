package main

import "fmt"

type IMan interface {
	Walk() int
}

type ManProxy struct {
	manProxy IMan // 不能是指针
}

func (proxy *ManProxy) setProxy(man IMan) {
	proxy.manProxy = man
}

func (proxy *ManProxy) Walk() int {
	proxy.manProxy.Walk()
	fmt.Println("proxy Walk")
	return 0
}

type Man struct {
	name string
	age  int
}

func (man *Man) Walk() int {
	fmt.Println("man Walk")
	return 0
}

func main() {
	manProxy := &ManProxy{} // 是不是指针都行
	var manImpl IMan
	man := Man{name: "xiaoming", age: 18}
	manImpl = &man //需要取地址, 调用man具体实现,而不是复制

	manProxy.setProxy(manImpl)
	manProxy.Walk()

	return
}
