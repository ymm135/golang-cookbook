package main

import "fmt"

func main() {
	defer fmt.Println("in main")
	defer func() {
		defer func() {
			fmt.Println("panic again and again")
			panic("panic again and again")
		}()
		fmt.Println("panic again")
		panic("panic again")
	}()

	panic("panic once")
}
