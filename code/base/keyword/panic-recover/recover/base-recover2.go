package main

import "fmt"

func main() {
	defer fmt.Println("in main")
	if err := recover(); err != nil { //δΈηζ
		fmt.Println(err)
	}

	panic("unknown err")
}
