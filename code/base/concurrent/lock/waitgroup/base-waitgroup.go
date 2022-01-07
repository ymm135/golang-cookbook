package main

import (
	"fmt"
	"sync"
)

func main() {
	group := sync.WaitGroup{}
	num := 5
	group.Add(num)

	for i := 0; i < num; i++ {
		index := i
		go func() {
			fmt.Println("run goroutine ", index)
			group.Done()
		}()
	}

	group.Wait()
}
