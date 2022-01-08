package main

import (
	"fmt"
	"time"
)

func main() {
	// 1.延时执行
	fmt.Println("currTime=", time.Now().Format("2006-01-02 15:04:05"))
	// create a nobuf channel and a goroutine `timer` will write it after 2 seconds
	timeAfterTrigger := time.After(time.Second * 2)
	// will be suspend but we have `timer` so will be not deadlocked
	curTime, _ := <-timeAfterTrigger
	// print current time
	fmt.Println("timeAfter=", curTime.Format("2006-01-02 15:04:05"))

	// 2.定时执行
	// 创建一个计时器
	timeTicker := time.NewTicker(time.Second * 2)
	i := 0
	for {
		if i > 5 {
			break
		}

		fmt.Println("timeTicker=", time.Now().Format("2006-01-02 15:04:05"))
		i++
		<-timeTicker.C // 下次触发的时间

	}
	// 清理计时器
	timeTicker.Stop()
}
