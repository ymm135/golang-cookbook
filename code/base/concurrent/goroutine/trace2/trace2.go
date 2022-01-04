package main

import (
	"fmt"
	"time"
)

func main() {
	for i := 0; i < 20; i++ {
		time.Sleep(time.Second)
		fmt.Println("Hello Go")
	}
}
