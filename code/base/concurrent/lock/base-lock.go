package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	mutex := sync.Mutex{}

	go func() {
		// 不可重入锁
		mutex.Lock() // 同一个协程不能加锁多次
		mutex.Lock() // 一直循环获取锁

		fmt.Println("Write ...")
	}()

	go func() {
		time.Sleep(time.Second * 2)
		defer mutex.Unlock() // 不管Lock多少次，只需要解锁一次即可

		fmt.Println("Read ...")
	}()

	time.Sleep(time.Second * 4)
}
