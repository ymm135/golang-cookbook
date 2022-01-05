package main

import (
	"fmt"
	"time"
)

func main() {
	// string
	str := "I am String!"
	for i, s := range str {
		fmt.Println("[string](", i, ")=", string(s))
	}

	// array slice
	array := []int{1, 3, 5, 7, 9}
	for i, v := range array {
		// 也可以使用array[i]
		fmt.Println("array(", i, ")=", v)
	}

	// hash
	hashTable := make(map[string]string, 10)
	hashTable["a"] = "array"
	hashTable["b"] = "bar"
	hashTable["c"] = "car"
	for k, v := range hashTable {
		fmt.Println("[hash]", k, ":", v)
	}

	//channel
	ch := make(chan string, 10)
	go func() {
		ch <- "hello"
		ch <- "go"
		ch <- "!"
	}()

	time.Sleep(time.Second)

	// 如果不在协程中开启, fatal error: all goroutines are asleep - deadlock!
	go func() {
		for c := range ch {
			fmt.Println("[channel]", c)
		}
	}()

	time.Sleep(time.Second)
}
