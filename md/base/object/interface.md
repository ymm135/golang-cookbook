# 接口    
## 数据结构  
`runtime/runtime2.go`定义`iface`
```go
type iface struct {
	tab  *itab
	data unsafe.Pointer
}

// layout of Itab known to compilers
// allocated in non-garbage-collected memory
// Needs to be in sync with
// ../cmd/compile/internal/gc/reflect.go:/^func.dumptabs.
type itab struct {
inter *interfacetype
_type *_type
hash  uint32 // copy of _type.hash. Used for type switches.
_     [4]byte
fun   [1]uintptr // variable sized. fun[0]==0 means _type does not implement inter.
}
```  

`interfacetype`定义在文件`runtime/type.go`  
```go
type interfacetype struct {
	typ     _type
	pkgpath name
	mhdr    []imethod   // method handler 
```

`接口`数据结构中会存储实现者的`指针`。  



