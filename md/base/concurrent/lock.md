- # 锁

目录:  
- [测试代码](#测试代码)
- [数据结构](#数据结构)
	- [sync.Mutex(互斥锁)](#syncmutex互斥锁)
		- [互斥锁的状态](#互斥锁的状态)
		- [互斥锁的`正常模式`与`饥饿模式`](#互斥锁的正常模式与饥饿模式)
		- [加解锁的实现](#加解锁的实现)
	- [读写锁(sync.RWMutex)](#读写锁syncrwmutex)
	- [集体等待`sync.WaitGroup`](#集体等待syncwaitgroup)
		- [数据结构及实现](#数据结构及实现)
	- [一次就好`sync.Once`](#一次就好synconce)
	- [条件变量`sync.Cond`](#条件变量synccond)
	- [信号量](#信号量)


`锁`是一种并发编程中的`同步原语`（Synchronization Primitives），它能保证多个 `Goroutine` 在访问同一片内存时不会出现`竞争条件`（Race condition）等问题。  
这些基本原语提供了较为基础的同步功能，但是它们是一种相对原始的同步机制，在多数情况下，我们都应该使用抽象层级更高的 `Channel` 实现同步。  
## 测试代码
```go
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
```

输出结果
```shell
Read ...
Write ...
```

## 数据结构
### sync.Mutex(互斥锁) 
源码路径`go/src/sync/mutex.go`  
```go
// A Mutex is a mutual exclusion lock. 
// The zero value for a Mutex is an unlocked mutex.
//
// A Mutex must not be copied after first use.
type Mutex struct {
	state int32   //当前互斥锁的状态
	sema  uint32  //用于控制锁状态的信号量
}
``` 

#### 互斥锁的状态  
```go
const (
	mutexLocked = 1 << iota // mutex is locked
	mutexWoken
	mutexStarving
	mutexWaiterShift = iota
	starvationThresholdNs = 1e6
)
```

- `mutexLocked` 表示互斥锁的锁定状态 
- `mutexWoken` 表示从正常模式被唤醒
- `mutexStarving` 当前的互斥锁进入饥饿状态

#### 互斥锁的`正常模式`与`饥饿模式`
[官方文档](https://github.com/ymm135/go/blob/debug.1.16.9/src/sync/mutex.go)  

> 互斥锁公平机制。  
> 互斥锁可以有两种操作模式：正常和饥饿。  
> 在正常模式下，等待获取锁的协程们按 FIFO 顺序排队，但被唤醒的等待着不拥有互斥锁并与新来的协程争夺所有权。  
> 新来的协程有一个优势——它们已经在 CPU 上运行，而且可能有很多，所以被唤醒的协程很有可能会获取锁失败。  
> 在这种情况下，它排在等待队列的前面。如果一个等待着超过1ms未能获取到mutex，它就会将mutex切换到饥饿模式。  
> 在饥饿模式下，互斥锁的所有权直接从解锁协程移交给队列前面的等待着。  
> 新来的协程不会尝试获取互斥锁，即使它看起来已解锁，也不会尝试自旋。相反，他们在等待队列的尾部排队。  
> 如果队列中的等待着拥有了互斥锁的所有权并发现它是队列中的最后一个等待着，或 它等待的时间少于 1 毫秒，则它将互斥锁切换回正常操作模式。  
> 普通模式的性能要好得多，因为即使有阻塞的等待程序，协程也可以连续多次获取互斥锁。  

总的来说就是正常模式处于抢占模式，谁等得到所有权算谁的，饥饿模式是排队模式，大家老老实实的排队，避免过度竞争，进入自旋模式等过多消耗资源。  

#### 加解锁的实现  

加锁的实现机制是:  
- 如果未加锁, 通过`CAS`操作把`m.state`的状态从`0`修改为`1`，代表已加锁。
- 如果已加锁   
  - 如果处于正常模式，一直循环获取锁(自旋锁)，获取锁的时间超过`1ms`之后，进入饥饿模式(抢占式、非公平锁)  
  - 如果处于饥饿模式，新来的协程不会尝试获取锁，所有的协程都入队，处于等待休眠状态(信号量)，如果被唤醒就依次获取锁(队列式，公平锁)  

```go
// Lock locks m.
// If the lock is already in use, the calling goroutine
// blocks until the mutex is available.
func (m *Mutex) Lock() {
    // 快速途径，获取未加锁的互斥锁 m.state 从 0 -> 1 
	// Fast path: grab unlocked mutex.
	if atomic.CompareAndSwapInt32(&m.state, 0, mutexLocked) {
		if race.Enabled {
			race.Acquire(unsafe.Pointer(m))
		}
		return
	}
	// Slow path (outlined so that the fast path can be inlined)
	m.lockSlow()
}

func (m *Mutex) lockSlow() {
	var waitStartTime int64
	starving := false
	awoke := false
	iter := 0
	old := m.state
	for {
		// 自旋锁， 不影响其他goroutine
		// Don't spin in starvation mode, ownership is handed off to waiters
		// so we won't be able to acquire the mutex anyway.
		if old&(mutexLocked|mutexStarving) == mutexLocked && runtime_canSpin(iter) {
			...
            runtime_doSpin()
            iter++
            old = m.state
            continue
		}
		
		// Don't try to acquire starving mutex, new arriving goroutines must queue.
		if atomic.CompareAndSwapInt32(&m.state, old, new) {
			if old&(mutexLocked|mutexStarving) == 0 {
				break // locked the mutex with CAS
				...
				// 进入休眠状态，等待被唤醒(信号量m.sema)
				runtime_SemacquireMutex(&m.sema, queueLifo, 1)
				// 计算饥饿时间
			    starving = starving || runtime_nanotime()-waitStartTime > starvationThresholdNs
				...
			}
		}
	}
}
``` 

解锁的实现
- 如果未加锁，尝试解锁，抛出异常`sync.Mutex.Unlock` 
- 通过`CAS`操作直接把`m.state`状态变为未加锁状态
- 如果没有修改成功
  - 处于正常模式， 如果没有等待释放的锁或已经被唤醒的协程，直接返回；其他情况通过`runtime_Semrelease`唤醒协程。
  - 处于饥饿模式，将锁所有权会交给队列的下个等待着，等待着会负责设置`mutexLocked`标志位。
  
```go
// Unlock unlocks m.
// It is a run-time error if m is not locked on entry to Unlock.
//
// A locked Mutex is not associated with a particular goroutine.
// It is allowed for one goroutine to lock a Mutex and then
// arrange for another goroutine to unlock it.
func (m *Mutex) Unlock() {
	if race.Enabled {
		_ = m.state
		race.Release(unsafe.Pointer(m))
	}

	// Fast path: drop lock bit.
	new := atomic.AddInt32(&m.state, -mutexLocked)
	if new != 0 {
		// Outlined slow path to allow inlining the fast path.
		// To hide unlockSlow during tracing we skip one extra frame when tracing GoUnblock.
		m.unlockSlow(new)
	}
}

func (m *Mutex) unlockSlow(new int32) {
	if (new+mutexLocked)&mutexLocked == 0 {
		throw("sync: unlock of unlocked mutex")
	}
	if new&mutexStarving == 0 {
		old := new
		for {
          if old>>mutexWaiterShift == 0 || old&(mutexLocked|mutexWoken|mutexStarving) != 0 {
                  return
              }
              // Grab the right to wake someone.
              new = (old - 1<<mutexWaiterShift) | mutexWoken
              if atomic.CompareAndSwapInt32(&m.state, old, new) {
                  runtime_Semrelease(&m.sema, false, 1)
                  return
              }
              old = m.state
          }
    } else {
        // 释放信号量的同时，就会唤醒等待的goroutine
		runtime_Semrelease(&m.sema, true, 1)
	}

```

### 读写锁(sync.RWMutex)
读写互斥锁 `sync.RWMutex` 是细粒度的互斥锁，它不限制资源的并发读，但是`读写、写写`操作无法并行执行。

```go
// There is a modified copy of this file in runtime/rwmutex.go.
// If you make any changes here, see if you should make them there.

// A RWMutex is a reader/writer mutual exclusion lock.
// The lock can be held by an arbitrary number of readers or a single writer.
// The zero value for a RWMutex is an unlocked mutex.
//
// A RWMutex must not be copied after first use.
//
// If a goroutine holds a RWMutex for reading and another goroutine might
// call Lock, no goroutine should expect to be able to acquire a read lock
// until the initial read lock is released. In particular, this prohibits
// recursive read locking. This is to ensure that the lock eventually becomes
// available; a blocked Lock call excludes new readers from acquiring the
// lock.
type RWMutex struct {
	w           Mutex  // held if there are pending writers
	writerSem   uint32 // semaphore for writers to wait for completing readers
	readerSem   uint32 // semaphore for readers to wait for completing writers
	readerCount int32  // number of pending readers
	readerWait  int32  // number of departing readers
}
``` 

测试代码:  
```go
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

	time.Sleep(time.Second * 2)
}
```
输出结果为
```
RLock
```

`写写`操作无法完成好理解，因为有把写锁`w Mutex`,如何实现读写不能同时进行呢?  
现在假设先读后写的流程:  

读锁的实现是把`readerCount`自增，调用一次，增加一个，如果有未释放的写锁，那就等待写锁释放后启动。(通过信号量`readerSem`)
```go
func (rw *RWMutex) RLock() {
	if race.Enabled {
		_ = rw.w.state
		race.Disable()
	}
	if atomic.AddInt32(&rw.readerCount, 1) < 0 {
		// A writer is pending, wait for it.
		runtime_SemacquireMutex(&rw.readerSem, false, 0)
	}
	if race.Enabled {
		race.Enable()
		race.Acquire(unsafe.Pointer(&rw.readerSem))
	}
}
```
下面查看`Lock`的实现:  
大概逻辑为先对写锁进行加锁，如果写锁已经占用，`rw.w.Lock()`就会阻塞。如果没有写锁，再去判断读锁的数量，如果读锁的数量不为0，
那就一直等待，直到读锁完全释放。(通过信号量`writerSem`实现)  
```go
// Lock locks rw for writing.
// If the lock is already locked for reading or writing,
// Lock blocks until the lock is available.
func (rw *RWMutex) Lock() {
	if race.Enabled {
		_ = rw.w.state
		race.Disable()
	}
	// First, resolve competition with other writers.
	rw.w.Lock()
	// Announce to readers there is a pending writer.
	r := atomic.AddInt32(&rw.readerCount, -rwmutexMaxReaders) + rwmutexMaxReaders
	// Wait for active readers.
	if r != 0 && atomic.AddInt32(&rw.readerWait, r) != 0 {
		runtime_SemacquireMutex(&rw.writerSem, false, 0)
	}
	if race.Enabled {
		race.Enable()
		race.Acquire(unsafe.Pointer(&rw.readerSem))
		race.Acquire(unsafe.Pointer(&rw.writerSem))
	}
}
```  

如果读锁个数为清零，写锁一直等待信号量`rw.writerSem`；读锁在调用`RUnlock`会减少读锁个数，如果读锁数量为`小于0`，
会通过`runtime_Semrelease(&rw.writerSem, false, 1)`唤醒信号量`rw.writerSem`，写锁协程就被唤醒了。  

> 写锁与写锁之间通过`Mutex`实现，读写锁之间通过`信号量`实现。  

### 集体等待`sync.WaitGroup`
[测试代码](../../../code/base/concurrent/lock/waitgroup/base-waitgroup.go)  等待所有协程运行结束后再退出  
```go
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
```

#### 数据结构及实现
源码路径``
```go
// A WaitGroup waits for a collection of goroutines to finish.
// The main goroutine calls Add to set the number of
// goroutines to wait for. Then each of the goroutines
// runs and calls Done when finished. At the same time,
// Wait can be used to block until all goroutines have finished.
//
// A WaitGroup must not be copied after first use.
type WaitGroup struct {
	noCopy noCopy // 保证不会被开发者通过再赋值的方式拷贝

	// 64-bit value: high 32 bits are counter, low 32 bits are waiter count.
	// 64-bit atomic operations require 64-bit alignment, but 32-bit
	// compilers do not ensure it. So we allocate 12 bytes and then use
	// the aligned 8 bytes in them as state, and the other 4 as storage
	// for the sema.
	state1 [3]uint32
}
``` 

从官网注释中可以看出`state1`存储这状态及信号量`statep, semap`，等待及唤醒功能通过`信号量`实现的 

`Wait()`代码实现，主要语句是`runtime_Semacquire(semap)`  
```go
// Wait blocks until the WaitGroup counter is zero.
func (wg *WaitGroup) Wait() {
	statep, semap := wg.state()

	for {
		state := atomic.LoadUint64(statep)
		v := int32(state >> 32)
		w := uint32(state)

		// Increment waiters count.
		if atomic.CompareAndSwapUint64(statep, state, state+1) {

			runtime_Semacquire(semap)
			if *statep != 0 {
				panic("sync: WaitGroup is reused before previous Wait has returned")
			}

			return
		}
	}
}
```


### 一次就好`sync.Once`

```go
// Once is an object that will perform exactly one action.
//
// A Once must not be copied after first use.
type Once struct {
	// done indicates whether the action has been performed.
	// It is first in the struct because it is used in the hot path.
	// The hot path is inlined at every call site.
	// Placing done first allows more compact instructions on some architectures (amd64/386),
	// and fewer instructions (to calculate offset) on other architectures.
	done uint32
	m    Mutex
}
```

测试demo
```go
func main() {
	once := sync.Once{}
	for i := 0; i < 10; i++ {
		once.Do(func() {
			fmt.Println("once run ")
		})
	}
}
```
运行结果: 
```shell
once run  
```

`once.Do`源码分析  
实现方式是通过`m    Mutex`互斥锁防止多线程运行，通过状态`done uint32`保存是否运行标记。  
```go
func (o *Once) Do(f func()) {
	if atomic.LoadUint32(&o.done) == 0 {
		// Outlined slow-path to allow inlining of the fast-path.
		o.doSlow(f)
	}
}

func (o *Once) doSlow(f func()) {
	o.m.Lock()
	defer o.m.Unlock()
	if o.done == 0 {
		defer atomic.StoreUint32(&o.done, 1)
		f()
	}
}
```


### 条件变量`sync.Cond` 

测试代码
```go
var status int64

func main() {
	c := sync.NewCond(&sync.Mutex{})
	for i := 0; i < 10; i++ {
		go listen(c)
	}
	time.Sleep(1 * time.Second)
	go broadcast(c)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch
}

func broadcast(c *sync.Cond) {
	c.L.Lock()
	atomic.StoreInt64(&status, 1)
	c.Broadcast()
	c.L.Unlock()
}

func listen(c *sync.Cond) {
	c.L.Lock()
	for atomic.LoadInt64(&status) != 1 {
		c.Wait()
	}
	fmt.Println("listen")
	c.L.Unlock()
}
```

输出结果:  
```shell
listen
listen
listen
listen
listen
listen
listen
listen
listen
listen
```

数据结构: 
```go
// Cond implements a condition variable, a rendezvous point
// for goroutines waiting for or announcing the occurrence
// of an event.
//
// Each Cond has an associated Locker L (often a *Mutex or *RWMutex),
// which must be held when changing the condition and
// when calling the Wait method.
//
// A Cond must not be copied after first use.
type Cond struct {
	noCopy noCopy

	// L is held while observing or changing the condition
	L Locker

	notify  notifyList
	checker copyChecker
}
```
方法及使用说明，源码注释写的很详细
```go
// Wait atomically unlocks c.L and suspends execution
// of the calling goroutine. After later resuming execution,
// Wait locks c.L before returning. Unlike in other systems,
// Wait cannot return unless awoken by Broadcast or Signal.
//
// Because c.L is not locked when Wait first resumes, the caller
// typically cannot assume that the condition is true when
// Wait returns. Instead, the caller should Wait in a loop:
//
//    c.L.Lock()
//    for !condition() {
//        c.Wait()
//    }
//    ... make use of condition ...
//    c.L.Unlock()
//
func (c *Cond) Wait() {
	c.checker.check()
	t := runtime_notifyListAdd(&c.notify)
	c.L.Unlock()
	runtime_notifyListWait(&c.notify, t)
	c.L.Lock()
}
```

单个唤醒和全部唤醒:  
```go
// Signal wakes one goroutine waiting on c, if there is any.
//
// It is allowed but not required for the caller to hold c.L
// during the call.
func (c *Cond) Signal() {
	c.checker.check()
	runtime_notifyListNotifyOne(&c.notify)
}

// Broadcast wakes all goroutines waiting on c.
//
// It is allowed but not required for the caller to hold c.L
// during the call.
func (c *Cond) Broadcast() {
	c.checker.check()
	runtime_notifyListNotifyAll(&c.notify)
}
```

从源码实现中可以看出`互斥锁`为了线程安全，等待及唤醒通过`信号量`实现。

`sync.Cond` 不是一个常用的同步机制，但是在条件长时间无法满足时，与使用 `for {}` 进行忙碌等待相比，
`sync.Cond` 能够让出处理器的使用权，提高 CPU 的利用率。使用时我们也需要注意以下问题：
- `sync.Cond.Wait` 在调用之前一定要使用获取互斥锁，否则会触发程序崩溃；
- `sync.Cond.Signal` 唤醒的 Goroutine 都是队列最前面、等待最久的 Goroutine；
- `sync.Cond.Broadcast` 会按照一定顺序广播通知等待的全部 Goroutine；

### 信号量 

源码路径:`go/src/sync/runtime.go`中的操作方法
```go
// Semacquire waits until *s > 0 and then atomically decrements it.
// It is intended as a simple sleep primitive for use by the synchronization
// library and should not be used directly.
func runtime_Semacquire(s *uint32)
```

具体实现`go/src/runtime/sema.go`  
```go
//go:linkname sync_runtime_Semacquire sync.runtime_Semacquire
func sync_runtime_Semacquire(addr *uint32) {
	semacquire1(addr, false, semaBlockProfile, 0)
}

//go:linkname poll_runtime_Semacquire internal/poll.runtime_Semacquire
func poll_runtime_Semacquire(addr *uint32) {
	semacquire1(addr, false, semaBlockProfile, 0)
}
``` 
`semacquire1`方法实现: 
```go
func semacquire1(addr *uint32, lifo bool, profile semaProfileFlags, skipframes int) {
	gp := getg()
	if gp != gp.m.curg {
		throw("semacquire not on the G stack")
	}

	// Easy case.
	if cansemacquire(addr) {
		return
	}

	// Harder case:
	//	increment waiter count
	//	try cansemacquire one more time, return if succeeded
	//	enqueue itself as a waiter
	//	sleep
	//	(waiter descriptor is dequeued by signaler)
	s := acquireSudog()
	root := semroot(addr)
	t0 := int64(0)
	s.releasetime = 0
	s.acquiretime = 0
	s.ticket = 0
	if profile&semaBlockProfile != 0 && blockprofilerate > 0 {
		t0 = cputicks()
		s.releasetime = -1
	}
	if profile&semaMutexProfile != 0 && mutexprofilerate > 0 {
		if t0 == 0 {
			t0 = cputicks()
		}
		s.acquiretime = t0
	}
	for {
		lockWithRank(&root.lock, lockRankRoot)
		// Add ourselves to nwait to disable "easy case" in semrelease.
		atomic.Xadd(&root.nwait, 1)
		// Check cansemacquire to avoid missed wakeup.
		if cansemacquire(addr) {
			atomic.Xadd(&root.nwait, -1)
			unlock(&root.lock)
			break
		}
		// Any semrelease after the cansemacquire knows we're waiting
		// (we set nwait above), so go to sleep.
		root.queue(addr, s, lifo)
		goparkunlock(&root.lock, waitReasonSemacquire, traceEvGoBlockSync, 4+skipframes)
		if s.ticket != 0 || cansemacquire(addr) {
			break
		}
	}
	if s.releasetime > 0 {
		blockevent(s.releasetime-t0, 3+skipframes)
	}
	releaseSudog(s)
}
```

通过信号量实现休眠的函数为`goparkunlock`唤醒的函数为`goready`  
```go
// Puts the current goroutine into a waiting state and unlocks the lock.
// The goroutine can be made runnable again by calling goready(gp).
func goparkunlock(lock *mutex, reason waitReason, traceEv byte, traceskip int) {
	gopark(parkunlock_c, unsafe.Pointer(lock), reason, traceEv, traceskip)
}

// 唤醒线程
func goready(gp *g, traceskip int) {
	systemstack(func() {
		ready(gp, traceskip, true)
	})
}
```

当没有接收者能够处理数据时，向`channel`的发送数据会被阻塞，用的的就是`goparkunlock`；  
`channel`数据发送也会用到`goready`,唤醒阻塞的接收者们。  

> 调用 `runtime.goparkunlock` 将当前的 `Goroutine` 陷入沉睡等待唤醒；
> 调用 `runtime.goready` 将等待接收数据的 `Goroutine` 标记成可运行状态 `Grunnable` 并把该 
`Goroutine` 放到发送方所在的处理器的 `runnext` 上等待执行，该处理器在下一次调度时会立刻唤醒数据的接收方； 
















