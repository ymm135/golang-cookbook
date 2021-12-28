package main

import "fmt"

func main() {
	s := "HelloWorld!"
	sm := modifyString(s)
	fmt.Println("source(", &s, "):", s, ",midify(", &s, "):", sm)
}

func modifyString(s string) string {
	bs := []byte(s)
	bs[0] = 77
	return string(bs)
}
