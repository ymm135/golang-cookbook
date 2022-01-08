# 定时器
## 测试代码
```go
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
```

输出结果:
```shell
currTime= 2022-01-08 09:56:28
timeAfter= 2022-01-08 09:56:30
timeTicker= 2022-01-08 09:56:30
timeTicker= 2022-01-08 09:56:32
timeTicker= 2022-01-08 09:56:34
timeTicker= 2022-01-08 09:56:36
timeTicker= 2022-01-08 09:56:38
timeTicker= 2022-01-08 09:56:40
```

从测试代码可以看出计时触发通过`管道`而不是回调函数，延时/定时任务通过向管道发送数据，另一端在管道接收数据。 


## 数据结构
源码位置`go/src/time/sleep.go` 
```
// The Timer type represents a single event.
// When the Timer expires, the current time will be sent on C,
// unless the Timer was created by AfterFunc.
// A Timer must be created with NewTimer or AfterFunc.
type Timer struct {
	C <-chan Time
	r runtimeTimer
}
```

定时任务
```
// A Ticker holds a channel that delivers ``ticks'' of a clock
// at intervals.
type Ticker struct {
	C <-chan Time // The channel on which the ticks are delivered.
	r runtimeTimer
}
```

首先查看延时任务`timeAfterTrigger := time.After(time.Second * 2)`
源码路径:`go/src/time/sleep.go`
```
// After waits for the duration to elapse and then sends the current time
// on the returned channel.
// It is equivalent to NewTimer(d).C.
// The underlying Timer is not recovered by the garbage collector
// until the timer fires. If efficiency is a concern, use NewTimer
// instead and call Timer.Stop if the timer is no longer needed.
func After(d Duration) <-chan Time {
	return NewTimer(d).C
}

//创建一个定时器Timer，包含通知的管道和计时模块`runtimeTimer`
// NewTimer creates a new Timer that will send
// the current time on its channel after at least duration d.
func NewTimer(d Duration) *Timer {
	c := make(chan Time, 1)
	t := &Timer{
		C: c,
		r: runtimeTimer{
			when: when(d), // runtimeNano() + int64(d) 当下时间+延迟时间=触发时间
			f:    sendTime,
			arg:  c,
		},
	}
	startTimer(&t.r)
	return t
}

//到时间后，会往管道发送当下时间
func sendTime(c interface{}, seq uintptr) {
	// Non-blocking send of time on c.
	// Used in NewTimer, it cannot block anyway (buffer).
	// Used in NewTicker, dropping sends on the floor is
	// the desired behavior when the reader gets behind,
	// because the sends are periodic.
	select {
	case c.(chan Time) <- Now():
	default:
	}
}
```

`startTimer(&t.r)`的实现在 `go/src/runtime/time.go:208` 
```
// startTimer adds t to the timer heap.
//go:linkname startTimer time.startTimer
func startTimer(t *timer) {
	if raceenabled {
		racerelease(unsafe.Pointer(t))
	}
	addtimer(t) //添加到计时器
}

// addtimer adds a timer to the current P.
// This should only be called with a newly created timer.
// That avoids the risk of changing the when field of a timer in some P's heap,
// which could cause the heap to become unsorted.
func addtimer(t *timer) {
	// when must be positive. A negative value will cause runtimer to
	// overflow during its delta calculation and never expire other runtime
	// timers. Zero will cause checkTimers to fail to notice the timer.
	if t.when <= 0 {
		throw("timer when must be positive")
	}
	if t.period < 0 {
		throw("timer period must be non-negative")
	}
	if t.status != timerNoStatus {
		throw("addtimer called with initialized timer")
	}
	t.status = timerWaiting

	when := t.when

	// Disable preemption while using pp to avoid changing another P's heap.
	mp := acquirem()

	pp := getg().m.p.ptr()
	lock(&pp.timersLock)
	cleantimers(pp)
	doaddtimer(pp, t) // doaddtimer adds t to the current P's heap.
	unlock(&pp.timersLock)
	
	// 如果它不会在 when 参数之前唤醒，则wakeNetPoller 唤醒在网络轮询器中休眠的线程； 或者它唤醒一个空闲的 P 以服务计时器和网络轮询器（如果还没有的话）。
	wakeNetPoller(when) 

	releasem(mp)
}

``` 




