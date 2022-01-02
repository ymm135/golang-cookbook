package main

import "fmt"

func main() {
	defer fmt.Println("in main")
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	panic("unknown err")
}
