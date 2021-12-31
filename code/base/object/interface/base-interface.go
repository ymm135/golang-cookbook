package main

import (
	"fmt"
	"unsafe"
)

type IMan interface {
	walk() int
}

type Man struct {
	name string
	age  int
}

func (man *Man) walk() int {
	fmt.Println("walk name:", man.name, ",age:", man.age)
	return man.age
}

func main() {
	var test_interface interface{}
	man := Man{name: "xiaoming", age: 18}
	var iman IMan

	iman = &man
	iman.walk()
	man.walk()

	fmt.Printf("iman  %T 占中的字节数是 %d \n", iman, unsafe.Sizeof(iman))
	fmt.Println(man, iman, test_interface)
}
