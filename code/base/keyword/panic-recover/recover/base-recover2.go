package main

import "fmt"

func main() {
	defer fmt.Println("in main")
	if err := recover(); err != nil { //不生效
		fmt.Println(err)
	}

	panic("unknown err")
}
