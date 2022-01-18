package main

import "fmt"

func main() {
	// make
	slice := make([]int, 0, 100)
	hash := make(map[int]bool, 10)
	ch := make(chan int, 5)

	// new
	i := new(int)
	var v int
	i = &v

	fmt.Println(slice, hash, ch, i)
}
