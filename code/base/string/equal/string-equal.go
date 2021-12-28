package main

import "fmt"

func main() {
	s1 := "Hello World"
	s2 := s1
	s1 = "Hello Go"
	fmt.Println("s1:", s1, ",s2:", s2)
}
