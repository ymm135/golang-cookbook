package main

import (
	"fmt"
	"time"
)

func main() {
	testParam1()
	testParam2()
}

func testParam1() {
	startedAt := time.Now()
	defer fmt.Println("testParam1=", time.Since(startedAt))

	time.Sleep(time.Second)
}

func testParam2() {
	startedAt := time.Now()
	defer func() { fmt.Println("testParam1=", time.Since(startedAt)) }()

	time.Sleep(time.Second)
}
