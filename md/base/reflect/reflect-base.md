# Go 反射及动态代理  
**什么是反射(Reflect)?** [维基百科](https://zh.wikipedia.org/zh-hans/%E5%8F%8D%E5%B0%84%E5%BC%8F%E7%BC%96%E7%A8%8B)    
在计算机学中，反射式编程（英语：reflective programming）或反射（英语：reflection），是指计算机程序在运行时（runtime）
[**可以访问、检测和修改它本身状态或行为的一种能力**]()。
用比喻来说，反射就是程序在运行的时候能够“观察”并且修改自己的行为。  

要注意术语“反射”和“内省”（type introspection）的关系。
内省（或称“自省”）机制仅指程序在运行时对自身信息（称为元数据）的检测；
反射机制不仅包括要能在运行时对程序自身信息进行检测，还要求程序能进一步根据这些信息改变程序状态或结构。

### 优点
支持反射的语言提供了一些在早期高级语言中难以实现的运行时特性。

- 可以在一定程度上避免硬编码，提供灵活性和通用性。
- 可以作为一个第一类物件发现并修改源代码的结构（如代码块、类、方法、协议等）。
- 可以在运行时像对待源代码语句一样动态解析字符串中可执行的代码（类似JavaScript的eval()函数），进而可将跟class或function匹配的字符串转换成class或function的调用或引用。
- 可以创建一个新的语言字节码解释器来给编程结构一个新的意义或用途。


### 缺点   
- 此技术的学习成本高。面向反射的编程需要较多的高级知识，包括框架、关系映射和对象交互，以实现更通用的代码执行。
- 同样因为反射的概念和语法都比较抽象，过多地滥用反射技术会使得代码难以被其他人读懂，不利于合作与交流。
- 由于将部分信息检查工作从编译期推迟到了运行期，此举在提高了代码灵活性的同时，牺牲了一点点运行效率。  

通过深入学习反射的特性和技巧，它的劣势可以尽量避免，但这需要许多时间和经验的积累。  

示例  
```
// Go 
import "reflect"

// Without reflection
f := Foo{}
f.Hello()

// With reflection
fT := reflect.TypeOf(Foo{})
fV := reflect.New(fT)

m := fV.MethodByName("Hello")
if m.IsValid() {
    m.Call(nil)
}

// Java
import java.lang.reflect.Method;

// Without reflection
Foo foo = new Foo();
foo.hello();

// With reflection
try {
    // Alternatively: Object foo = Foo.class.newInstance();
    Object foo = Class.forName("complete.classpath.and.Foo").newInstance();

    Method m = foo.getClass().getDeclaredMethod("hello", new Class<?>[0]);
    m.invoke(foo);
} catch (Exception e) {
    // Catching ClassNotFoundException, NoSuchMethodException
    // InstantiationException, IllegalAccessException
}
```

## Go 反射基础  

通过反射获取结构体的信息，包含tag，通过匹配tag与map中的key，把map的值填充到结构体中  
```go 
type Foo struct {
	Name string `ziduan:"name"`
	Age  int
}

func main() {
	var values map[string]interface{}
	values = make(map[string]interface{})

	values["name"] = "xiaoming"
	foo := Foo{}
	fV := reflect.ValueOf(&foo)
	fT := reflect.TypeOf(&foo)
	fmt.Println(fV)

	mapValues := reflect.ValueOf(values)

	if mapValues.Kind() == reflect.Ptr {
		mapValues = mapValues.Elem()
	}

	// Go语言程序中对指针获取反射对象时，可以通过 reflect.Elem() 方法获取这个指针指向的元素类型，这个获取过程被称为取元素
	if fV.Kind() == reflect.Ptr {
		fV = fV.Elem()
	}

	if fT.Kind() == reflect.Ptr {
		fT = fT.Elem()
	}

	//对象本身是个指针,需要获取原有类型
	// Elem returns a type's element type.
	numField := fV.NumField() // 使用反射获取结构体的成员类型 NumField() 和 Field()
	for i := 0; i < numField; i++ {
		fvFieldVal := fV.Field(i)
		ftFieldVal := fT.Field(i)

		for key, value := range values {
			mapVal := reflect.ValueOf(value)

			//{Name  string ziduan:"name" 0 [0] false}
			fmt.Println(ftFieldVal)

			// 这里使用TypeOf 获取Tag,
			if strings.Compare(ftFieldVal.Tag.Get("ziduan"), key) == 0 { // 结构体标签（Struct Tag）
				if mapVal.Type().Kind() == reflect.String { // TypeOf Kind用于判断类型
					fvFieldVal.Set(reflect.ValueOf(value))
				}
			}
		}
	}

	// {xiaoming 0}
	fmt.Println(foo)
}
```

通过反射可以获取结构体信息: 
```
// A StructField describes a single field in a struct.
type StructField struct {
	// Name is the field name.
	Name string
	// PkgPath is the package path that qualifies a lower case (unexported)
	// field name. It is empty for upper case (exported) field names.
	// See https://golang.org/ref/spec#Uniqueness_of_identifiers
	PkgPath string

	Type      Type      // field type
	Tag       StructTag // field tag string
	Offset    uintptr   // offset within struct, in bytes
	Index     []int     // index sequence for Type.FieldByIndex
	Anonymous bool      // is an embedded field
}

// A StructTag is the tag string in a struct field.
//
// By convention, tag strings are a concatenation of
// optionally space-separated key:"value" pairs.
// Each key is a non-empty string consisting of non-control
// characters other than space (U+0020 ' '), quote (U+0022 '"'),
// and colon (U+003A ':').  Each value is quoted using U+0022 '"'
// characters and Go string literal syntax.
type StructTag string
```

样例: 
```
type Foo struct {
	Name string `ziduan:"name"`
	Age  int
}
```

| 字段 | Name 成员 | Age 成员 | 
| ---- | -------- | ------- | 
| field name | Name | Age | 
| pkg path | | 
| field type | string | int |
| field tag | ziduan:"name" | | 
| offset within struct, in bytes | 0 | 16 | 
| index sequence for Type.FieldByIndex | [0] | [1] | 
| anonymous bool | false | false | 

## Go 反射核心  
`reflect/type.go`中定义类型

```
// Type is the representation of a Go type.
//
// Not all methods apply to all kinds of types. Restrictions,
// if any, are noted in the documentation for each method.
// Use the Kind method to find out the kind of type before
// calling kind-specific methods. Calling a method
// inappropriate to the kind of type causes a run-time panic.
//
// Type values are comparable, such as with the == operator,
// so they can be used as map keys.
// Two Type values are equal if they represent identical types.
type Type interface {
	// Methods applicable to all types.

	// Align returns the alignment in bytes of a value of
	// this type when allocated in memory.
	Align() int

	// FieldAlign returns the alignment in bytes of a value of
	// this type when used as a field in a struct.
	FieldAlign() int

	// Method returns the i'th method in the type's method set.
	// It panics if i is not in the range [0, NumMethod()).
	//
	// For a non-interface type T or *T, the returned Method's Type and Func
	// fields describe a function whose first argument is the receiver,
	// and only exported methods are accessible.
	//
	// For an interface type, the returned Method's Type field gives the
	// method signature, without a receiver, and the Func field is nil.
	//
	// Methods are sorted in lexicographic order.
	Method(int) Method

	// MethodByName returns the method with that name in the type's
	// method set and a boolean indicating if the method was found.
	//
	// For a non-interface type T or *T, the returned Method's Type and Func
	// fields describe a function whose first argument is the receiver.
	//
	// For an interface type, the returned Method's Type field gives the
	// method signature, without a receiver, and the Func field is nil.
	MethodByName(string) (Method, bool)

	// NumMethod returns the number of methods accessible using Method.
	//
	// Note that NumMethod counts unexported methods only for interface types.
	NumMethod() int

	// Name returns the type's name within its package for a defined type.
	// For other (non-defined) types it returns the empty string.
	Name() string

	// PkgPath returns a defined type's package path, that is, the import path
	// that uniquely identifies the package, such as "encoding/base64".
	// If the type was predeclared (string, error) or not defined (*T, struct{},
	// []int, or A where A is an alias for a non-defined type), the package path
	// will be the empty string.
	PkgPath() string

	// Size returns the number of bytes needed to store
	// a value of the given type; it is analogous to unsafe.Sizeof.
	Size() uintptr

	// String returns a string representation of the type.
	// The string representation may use shortened package names
	// (e.g., base64 instead of "encoding/base64") and is not
	// guaranteed to be unique among types. To test for type identity,
	// compare the Types directly.
	String() string

	// Kind returns the specific kind of this type.
	Kind() Kind

	// Implements reports whether the type implements the interface type u.
	Implements(u Type) bool

	// AssignableTo reports whether a value of the type is assignable to type u.
	AssignableTo(u Type) bool

	// ConvertibleTo reports whether a value of the type is convertible to type u.
	ConvertibleTo(u Type) bool

	// Comparable reports whether values of this type are comparable.
	Comparable() bool

	// Methods applicable only to some types, depending on Kind.
	// The methods allowed for each kind are:
	//
	//	Int*, Uint*, Float*, Complex*: Bits
	//	Array: Elem, Len
	//	Chan: ChanDir, Elem
	//	Func: In, NumIn, Out, NumOut, IsVariadic.
	//	Map: Key, Elem
	//	Ptr: Elem
	//	Slice: Elem
	//	Struct: Field, FieldByIndex, FieldByName, FieldByNameFunc, NumField

	// Bits returns the size of the type in bits.
	// It panics if the type's Kind is not one of the
	// sized or unsized Int, Uint, Float, or Complex kinds.
	Bits() int

	// ChanDir returns a channel type's direction.
	// It panics if the type's Kind is not Chan.
	ChanDir() ChanDir

	// IsVariadic reports whether a function type's final input parameter
	// is a "..." parameter. If so, t.In(t.NumIn() - 1) returns the parameter's
	// implicit actual type []T.
	//
	// For concreteness, if t represents func(x int, y ... float64), then
	//
	//	t.NumIn() == 2
	//	t.In(0) is the reflect.Type for "int"
	//	t.In(1) is the reflect.Type for "[]float64"
	//	t.IsVariadic() == true
	//
	// IsVariadic panics if the type's Kind is not Func.
	IsVariadic() bool

	// Elem returns a type's element type.
	// It panics if the type's Kind is not Array, Chan, Map, Ptr, or Slice.
	Elem() Type

	// Field returns a struct type's i'th field.
	// It panics if the type's Kind is not Struct.
	// It panics if i is not in the range [0, NumField()).
	Field(i int) StructField

	// FieldByIndex returns the nested field corresponding
	// to the index sequence. It is equivalent to calling Field
	// successively for each index i.
	// It panics if the type's Kind is not Struct.
	FieldByIndex(index []int) StructField

	// FieldByName returns the struct field with the given name
	// and a boolean indicating if the field was found.
	FieldByName(name string) (StructField, bool)

	// FieldByNameFunc returns the struct field with a name
	// that satisfies the match function and a boolean indicating if
	// the field was found.
	//
	// FieldByNameFunc considers the fields in the struct itself
	// and then the fields in any embedded structs, in breadth first order,
	// stopping at the shallowest nesting depth containing one or more
	// fields satisfying the match function. If multiple fields at that depth
	// satisfy the match function, they cancel each other
	// and FieldByNameFunc returns no match.
	// This behavior mirrors Go's handling of name lookup in
	// structs containing embedded fields.
	FieldByNameFunc(match func(string) bool) (StructField, bool)

	// In returns the type of a function type's i'th input parameter.
	// It panics if the type's Kind is not Func.
	// It panics if i is not in the range [0, NumIn()).
	In(i int) Type

	// Key returns a map type's key type.
	// It panics if the type's Kind is not Map.
	Key() Type

	// Len returns an array type's length.
	// It panics if the type's Kind is not Array.
	Len() int

	// NumField returns a struct type's field count.
	// It panics if the type's Kind is not Struct.
	NumField() int

	// NumIn returns a function type's input parameter count.
	// It panics if the type's Kind is not Func.
	NumIn() int

	// NumOut returns a function type's output parameter count.
	// It panics if the type's Kind is not Func.
	NumOut() int

	// Out returns the type of a function type's i'th output parameter.
	// It panics if the type's Kind is not Func.
	// It panics if i is not in the range [0, NumOut()).
	Out(i int) Type

	common() *rtype
	uncommon() *uncommonType
}
```

## 实际开发中应用  
- ### 每次都需要写sql查询语句条件，能不能根据结构体中tag及数据类型，自动生成(Int/String/Slice)数据类型sql条件    

```go
package main

import (
	"fmt"
	"reflect"
	"strconv"
)

type Level int

//威胁规则列表展示内容
type ThreatRuleList struct {
	ThreatName   string `json:"threatName" where:"threatName"`                    // 威胁名称
	Category     string `json:"category" gorm:"column:category" where:"category"` // 分类
	PublishDate  string `json:"publishDate" gorm:"column:publishDate"`            // 发布日期
	RiskLevel    Level  `json:"riskLevel" gorm:"column:riskLevel"`                // 威胁等级
	Status       int    `json:"status" sql:"status"`                              // 规则开启状态
	SignatureIds []int  `json:"signatureId" where:"signatureId"`                  // sid
}

func main() {
	s := make([]int, 2)
	s[0] = 2
	s[1] = 4

	t := ThreatRuleList{
		ThreatName:   "TestTreat",
		Status:       1,
		Category:     "ddd",
		SignatureIds: s,
	}

	fmt.Println(GenerateWhereSql(t))
}

const TAG_NAME = "where"

// 通过遍历结构体内每个字段，生成where sql语句
// 比如结构体内有A,B两个成员，where A = ? And B = ?
// 优先获取gorm:"column:threatName" 再是 json
func GenerateWhereSql(s interface{}) string {
	// 首先获取结构体内所有的成员
	if s == nil {
		return ""
	}

	sV := reflect.ValueOf(s)
	sT := reflect.TypeOf(s)

	if sV.Kind() == reflect.Ptr {
		sV = sV.Elem()
	}

	if sT.Kind() == reflect.Ptr {
		sT = sT.Elem()
	}

	sql := " where 1=1 "
	numField := sV.NumField() // 使用反射获取结构体的成员类型 NumField() 和 Field()
	for i := 0; i < numField; i++ {
		svFieldVal := sV.Field(i)
		stFieldVal := sT.Field(i)

		// 先获取gorm的字段
		gormTag := stFieldVal.Tag.Get(TAG_NAME)
		if len(gormTag) == 0 {
			continue
		}

		kind := svFieldVal.Kind()
		value := ""
		switch kind {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			intVal := svFieldVal.Int()
			if intVal != -1 { // 如果值为-1， 那就不作为where条件
				value = strconv.Itoa(int(intVal))
			}
		case reflect.String:
			value = "'" + svFieldVal.String() + "'"
		case reflect.Slice:
			length := svFieldVal.Len()
			if length == 0 {
				continue
			}

			sliceValType := reflect.Invalid

			value = " ( "
			for i := 0; i < length; i++ {

				sliceVal := svFieldVal.Index(i) // 获取切片存储的数据类型
				if sliceValType == reflect.Invalid {
					sliceValType = sliceVal.Kind()
				}
				switch sliceValType {
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					value += strconv.Itoa(int(sliceVal.Int()))
				case reflect.String:
					strVal := sliceVal.String()
					if len(strVal) != 0 {
						value += "'" + strVal + "'"
					}
				}

				if i < length-1 {
					value += ", "
				}
			}
			value += " ) "
		}

		if len(value) != 0 {
			if kind == reflect.Slice {
				sql += " and " + gormTag + " in " + value
			} else {
				sql += " and " + gormTag + " = " + value
			}
		}
	}

	return sql
}

``` 

运行结果: 

```sql
 where 1=1  and threatName = 'TestTreat' and category = 'ddd' and signatureId in  ( 2, 4 )   
```




