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
	manProxy := &ManProxy{}
	var manImpl IMan
	var man Man
	manImpl = &man //需要取地址, 调用man具体实现,而不是复制
	manProxy.setProxy(manImpl)
	manProxy.Walk()

	return
}
