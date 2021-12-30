# defer  
## 执行顺序  
先声明后执行，相当于堆栈，`先进后出`  
```go
func main() {
	defer println("A")
	defer println("B")
	defer println("C")
}
```  

输出结果:  
```shell
C
B
A
```  

## 实现原理  