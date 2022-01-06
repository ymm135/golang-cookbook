package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	mutex := &sync.Mutex{}

	go func() {
		fmt.Println("goroutine 1 enter")
		mutex.Lock()
		defer mutex.Unlock()
		time.Sleep(time.Second * 2)
	}()

	go func() {
		time.Sleep(time.Second)
		fmt.Println("goroutine 2 enter")

		mutex.Lock()
		defer mutex.Unlock()
		time.Sleep(time.Second * 4)
	}()

	go func() {
		time.Sleep(time.Second)
		fmt.Println("goroutine 3 enter")

		mutex.Lock()
		defer mutex.Unlock()
		time.Sleep(time.Second * 5)
	}()

	time.Sleep(time.Minute)
}
