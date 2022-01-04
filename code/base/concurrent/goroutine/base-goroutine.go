package main

import (
	"fmt"
	"time"
)

func main() {
	for i := 0; i < 10; i++ {
		tag := i

		go func() {
			for {
				fmt.Println("goroutine:", tag)
				time.Sleep(time.Microsecond * 100)
			}
		}()
	}
	time.Sleep(time.Minute)
}
