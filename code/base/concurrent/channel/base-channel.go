package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan string, 1)

	go func() {
		ch <- "hello channel"
		ch <- "!"
		ch <- "quit"
	}()

	//需要延迟一段时间，要不然接收不到数据
	time.Sleep(time.Second)

	str := <-ch
	fmt.Println("recv:", str)

	for { //
		select {
		case str := <-ch:
			fmt.Println(str)
			if str == "quit" {
				goto end
			}
		default:
			fmt.Println("default")
		}
	}
end:
}
