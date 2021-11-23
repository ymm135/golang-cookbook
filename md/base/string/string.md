# string  
[参考文章](https://draveness.me/golang/docs/part2-foundation/ch03-datastructure/golang-string/)  
## 数据结构  
源码在`go1.16.9/src/reflect/value.go`  

```
// StringHeader is the runtime representation of a string.
// It cannot be used safely or portably and its representation may
// change in a later release.
// Moreover, the Data field is not sufficient to guarantee the data
// it references will not be garbage collected, so programs must keep
// a separate, correctly typed pointer to the underlying data.
type StringHeader struct {
	Data uintptr
	Len  int
}
```  
与切片的结构体相比，字符串只少了一个表示容量的 Cap 字段，而正是因为切片在 Go 语言的运行时表示与字符串高度相似，所以我们经常会说字符串是一个只读的切片类型。  
```
type SliceHeader struct {
	Data uintptr
	Len  int
	Cap  int
}
```
因为字符串作为只读的类型，我们并不会直接向字符串直接追加元素改变其本身的内存空间，所有在字符串上的写入操作都是通过拷贝实现的。  

