# gin参数绑定
## 通过反射解析query参数
### shouldBindQuery源码
比如需要解析query参数，填充的结构体为:  
```
// 周报搜索结构体
type WtReportsSearch struct {
	CurrUserId uint   `form:"currUserId"`
	UserId     uint   `form:"userId"`
	StartTime  string `form:"startTime" example:"2021-11-04 12:36:34"`
	EndTime    string `form:"endTime"`
	Content    string `form:"content" example:"xx项目"`
	request.PageInfo
}
```

gin解析代码和请求参数:  
```
// url参数: userId=1&content=工作&startTime=2021-11-04 01:11:07&endTime=2021-11-04 03:11:08
func (wtReportsApi *WtReportsApi) GetWtReportsList(c *gin.Context) {
	var searchInfo wtReq.WtReportsSearch
	_ = c.ShouldBindQuery(&searchInfo)
	...
}
```

ShouldBindQuery方法，传入实现: **binding.Query**    
```
// gin@v1.7.4/context.go
// ShouldBindQuery is a shortcut for c.ShouldBindWith(obj, binding.Query).
func (c *Context) ShouldBindQuery(obj interface{}) error {
	return c.ShouldBindWith(obj, binding.Query)
}

//调用
// ShouldBindWith binds the passed struct pointer using the specified binding engine.
// See the binding package.
func (c *Context) ShouldBindWith(obj interface{}, b binding.Binding) error {
	return b.Bind(c.Request, obj)
}

//调用接口
// Binding describes the interface which needs to be implemented for binding the
// data present in the request such as JSON request body, query parameters or
// the form POST.
type Binding interface {
	Name() string
	Bind(*http.Request, interface{}) error
}

//方法实现 gin@v1.7.4/binding/query.go
func (queryBinding) Bind(req *http.Request, obj interface{}) error {
	values := req.URL.Query()
	if err := mapForm(obj, values); err != nil {
		return err
	}
	return validate(obj)
}

// gin@v1.7.4/binding/form_mapping.go  
// 需要注意"form"表单, 类型也有json 
func mapForm(ptr interface{}, form map[string][]string) error {
	return mapFormByTag(ptr, form, "form")
}

// 最终把参数转为map, 通过map匹配struct的tag，通过反射Value、Field设定值。  

```

## 自己动手写的demo
```
package main

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type WtReportsSearch struct {
	CurrUserId uint   `form:"currUserId"`
	UserId     uint   `form:"userId"`
	StartTime  string `form:"startTime" example:"2021-11-04 12:36:34"`
	EndTime    string `form:"endTime"`
	Content    string `form:"content" example:"xx项目"`
}

func bindQuery(obj interface{}) {
	var values map[string]interface{}
	values = make(map[string]interface{})

	values["currUserId"] = 1
	values["userId"] = 2
	values["content"] = "项目工作"

	prtObjTyp := reflect.TypeOf(obj)
	prtObjVal := reflect.ValueOf(obj)

	// 必须对象是ValueOf ，如果是TypeOf，不能修改
	//objVal := reflect.ValueOf(obj)
	//field0 := objVal.Elem().Field(0)

	//objVal := reflect.TypeOf(obj)
	//field0 := objVal.Elem().Field(0)
	//filed0Val := reflect.ValueOf(field0)
	//filed0Val.SetUint(2)

	ptrMapVal := reflect.ValueOf(values)

	if ptrMapVal.Kind() == reflect.Ptr {
		ptrMapVal = ptrMapVal.Elem()
	}

	// Go语言程序中对指针获取反射对象时，可以通过 reflect.Elem() 方法获取这个指针指向的元素类型，这个获取过程被称为取元素
	if prtObjVal.Kind() == reflect.Ptr {
		prtObjVal = prtObjVal.Elem()
		prtObjTyp = prtObjTyp.Elem()
	}

	//对象本身是个指针,需要获取原有类型
	// Elem returns a type's element type.
	numField := prtObjVal.NumField() // 使用反射获取结构体的成员类型 NumField() 和 Field()
	for i := 0; i < numField; i++ {

		fieldVal := prtObjVal.Field(i)
		fieldTpy := prtObjTyp.Field(i)
		fmt.Println(fieldVal, fieldTpy) // {CurrUserId  uint form:"currUserId" 0 [0] false}

		for key, value := range values {
			mapVal := reflect.ValueOf(value)
			fmt.Println("mapVal=", mapVal, ", canSet=", mapVal.CanSet())

			if strings.Compare(fieldTpy.Tag.Get("form"), key) == 0 { // 结构体标签（Struct Tag）
				// 值能被修改的条件: 可被寻址, 可被设置
				fmt.Println("fileVal interface=", fieldVal.Interface(), ", canSet=", fieldVal.CanSet(), ",CanAddr=", fieldVal.CanAddr())

				if mapVal.Type().Kind() == reflect.String { // Kind用于判断类型
					fieldVal.Set(reflect.ValueOf(value))
				}

				if mapVal.Type().Kind() == reflect.Int {
					str := fmt.Sprintf("%v", value)
					atoi, _ := strconv.Atoi(str)
					fieldVal.SetUint(uint64(atoi))
				}
				fmt.Println(value)
			}
		}

	}
}

// 通过反射给结构体赋值
func main() {
	var searchInfo WtReportsSearch
	bindQuery(&searchInfo)

	fmt.Println(searchInfo)
}
```  

## reflect.New实现动态代理

首先展示通过方法传入接口依赖 [code](../../../code/reflect/proxy/main.go) :  
```
package main

import "fmt"

type IMan interface {
	Walk() int
}

type ManProxy struct {
	manProxy IMan // 不能是指针
}

func (proxy *ManProxy) setProxy(man IMan) () {
	proxy.manProxy = man
}

func (proxy *ManProxy) Walk() int {
	proxy.manProxy.Walk()
	fmt.Println("proxy Walk")
}

type Man struct {
}

func (man *Man) Walk() int {
	fmt.Println("man Walk")
}

func main() {
	manProxy := &ManProxy{}
	var manImpl IMan
	var man Man
	manImpl = &man //需要取地址, 调用man具体实现,而不是复制
	manProxy.setProxy(manImpl)
}
```  
如果通过结构名或者文件名动态加载， 那就只能通过自定义配置文件匹配静态绑定。  




