package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	mutex := sync.RWMutex{}
	go func() {
		mutex.RLock()
		fmt.Println("RLock")
	}()

	go func() {
		time.Sleep(time.Second)
		mutex.Lock()
		fmt.Println("Lock")
	}()

	time.Sleep(time.Minute)
}
