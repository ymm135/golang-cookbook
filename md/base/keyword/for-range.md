# for 和 range
[参考文章]()
循环是所有编程语言都有的控制结构，除了使用经典的三段式循环之外，
`Go` 语言还引入了另一个关键字 `range` 帮助我们快速遍历`数组`、`切片`、`哈希表`以及 `Channel` 等集合类型。

## 测试代码 

```go
func main() {
	// string
	str := "I am String!"
	for i, s := range str {
		fmt.Println("[string](", i, ")=", string(s))
	}

	// array slice
	array := []int{1, 3, 5, 7, 9}
	for i, v := range array {
		// 也可以使用array[i]
		fmt.Println("array(", i, ")=", v)
	}

	// hash
	hashTable := make(map[string]string, 10)
	hashTable["a"] = "array"
	hashTable["b"] = "bar"
	hashTable["c"] = "car"
	for k, v := range hashTable {
		fmt.Println("[hash]", k, ":", v)
	}

	//channel
	ch := make(chan string, 10)
	go func() {
		ch <- "hello"
		ch <- "go"
		ch <- "!"
	}()

	time.Sleep(time.Second)

	// 如果不在协程中开启, fatal error: all goroutines are asleep - deadlock!
	go func() {
		for c := range ch {
			fmt.Println("[channel]", c)
		}
	}()

	time.Sleep(time.Second)
}
``` 

结果输出:
```shell
[string]( 0 )= I
[string]( 1 )=  
[string]( 2 )= a
[string]( 3 )= m
[string]( 4 )=  
[string]( 5 )= S
[string]( 6 )= t
[string]( 7 )= r
[string]( 8 )= i
[string]( 9 )= n
[string]( 10 )= g
[string]( 11 )= !
array( 0 )= 1
array( 1 )= 3
array( 2 )= 5
array( 3 )= 7
array( 4 )= 9
[hash] a : array
[hash] b : bar
[hash] c : car
[channel] hello
[channel] go
[channel] !
``` 

## `for`和`range`的实现
在编译阶段，会针对不同场景的`range`做不同的解析 
```go
// cmd/compile/internal/gc/walk.go
// The result of walkstmt MUST be assigned back to n, e.g.
// 	n.Left = walkstmt(n.Left)
func walkstmt(n *Node) *Node {
    ...
    case ORANGE:
		n = walkrange(n)
    ...
}

//cmd/compile/internal/gc/range.go
// walkrange transforms various forms of ORANGE into
// simpler forms.  The result must be assigned back to n.
// Node n may also be modified in place, and may also be
// the returned node.
func walkrange(n *Node) *Node {
    switch t.Etype {
	default:
		Fatalf("walkrange")

	case TARRAY, TSLICE:
	    ...
	case TMAP:
	    ...
	case TCHAN:
	    ...
	case TSTRING:
	    ...
}
```
从`range.go`的实现来看，针对`TARRAY/TSLICE`、`TMAP`、`TCHAN`、`TSTRING` 5中数据类型都有具体的实现  

### `TARRAY/TSLICE`类型处理 

```go
func walkrange(n *Node) *Node {
    ...
    case TARRAY, TSLICE:
            if arrayClear(n, v1, v2, a) {
                lineno = lno
                return n
            }
    
            // order.stmt arranged for a copy of the array/slice variable if needed.
            ha := a
    
            hv1 := temp(types.Types[TINT])
            hn := temp(types.Types[TINT])
    
            init = append(init, nod(OAS, hv1, nil))
            init = append(init, nod(OAS, hn, nod(OLEN, ha, nil)))
    
            n.Left = nod(OLT, hv1, hn)
            n.Right = nod(OAS, hv1, nod(OADD, hv1, nodintconst(1)))
    
            // for range ha { body }
            if v1 == nil {
                break
            }
    
            // for v1 := range ha { body }
            if v2 == nil {
                body = []*Node{nod(OAS, v1, hv1)}
                break
            }
    
            // for v1, v2 := range ha { body }
            if cheapComputableIndex(n.Type.Elem().Width) {
                // v1, v2 代表index, value, v2还是数组+索引 a[hv1] 
                // v1, v2 = hv1, ha[hv1]
                tmp := nod(OINDEX, ha, hv1)
                tmp.SetBounded(true)
                // Use OAS2 to correctly handle assignments
                // of the form "v1, a[v1] := range".
                a := nod(OAS2, nil, nil)
                a.List.Set2(v1, v2)
                a.Rlist.Set2(hv1, tmp)
                body = []*Node{a}
                break
            }
    
            // TODO(austin): OFORUNTIL is a strange beast, but is
            // necessary for expressing the control flow we need
            // while also making "break" and "continue" work. It
            // would be nice to just lower ORANGE during SSA, but
            // racewalk needs to see many of the operations
            // involved in ORANGE's implementation. If racewalk
            // moves into SSA, consider moving ORANGE into SSA and
            // eliminating OFORUNTIL.
    
            // TODO(austin): OFORUNTIL inhibits bounds-check
            // elimination on the index variable (see #20711).
            // Enhance the prove pass to understand this.
            ifGuard = nod(OIF, nil, nil)
            ifGuard.Left = nod(OLT, hv1, hn)
            translatedLoopOp = OFORUNTIL
    
            hp := temp(types.NewPtr(n.Type.Elem()))
            tmp := nod(OINDEX, ha, nodintconst(0))
            tmp.SetBounded(true)
            init = append(init, nod(OAS, hp, nod(OADDR, tmp, nil)))
    
            // Use OAS2 to correctly handle assignments
            // of the form "v1, a[v1] := range".
            a := nod(OAS2, nil, nil)
            a.List.Set2(v1, v2)
            a.Rlist.Set2(hv1, nod(ODEREF, hp, nil))
            body = append(body, a)
    
            // Advance pointer as part of the late increment.
            //
            // This runs *after* the condition check, so we know
            // advancing the pointer is safe and won't go past the
            // end of the allocation.
            a = nod(OAS, hp, addptr(hp, t.Elem().Width))
            a = typecheck(a, ctxStmt)
            n.List.Set1(a)
    ...
``` 

最终转化成
```go
ha := a
hv1 := 0
hn := len(ha)

v1 := hv1
v2 := nil

for ; hv1 < hn; hv1++ {
    tmp := ha[hv1]
    v1, v2 = hv1, tmp
    ...
}
```

可以看到`range`语句最终还是会转化为`for`循环，转换的过程是通过编译器实现的，通过修改`语法树`。

