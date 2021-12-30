package main

import "fmt"

type Man struct {
	name string
	age  int
}

func (man *Man) Walk() int {
	fmt.Println("man Walk")
	return 0
}

func main() {
	man := Man{name: "xiaoming", age: 18}

	var walkFun func() int
	walkFun = man.Walk
	walkFun()
}
