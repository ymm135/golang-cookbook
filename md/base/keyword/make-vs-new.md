# make vs new 
- make 的作用是初始化内置的数据结构，也就是我们在前面提到的切片、哈希表和 Channel2；
- new 的作用是根据传入的类型分配一片内存空间并返回指向这片内存空间的指针3；

示例  
```go
// make 
slice := make([]int, 0, 100)
hash := make(map[int]bool, 10)
ch := make(chan int, 5)

// new
i := new(int)
var v int
i := &v
```  

