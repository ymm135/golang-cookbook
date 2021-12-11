package main

import (
	"fmt"
	"strconv"
)

func main() {
	m := map[string]string{
		"name":    "ccmouse",
		"course":  "golang",
		"site":    "imooc",
		"quality": "notbad",
	}

	for i := 0; i < 21; i++ {
		m[strconv.Itoa(i)] = strconv.Itoa(i)
	}

	delete(m, "name")
	fmt.Println("Hello Go")
}
