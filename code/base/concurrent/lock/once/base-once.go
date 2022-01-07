package main

import (
	"fmt"
	"sync"
)

func main() {
	once := sync.Once{}
	for i := 0; i < 10; i++ {
		once.Do(func() {
			fmt.Println("once run ")
		})
	}
}
