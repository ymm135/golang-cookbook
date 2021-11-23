package main

import "fmt"

func main() {
	m := map[string]string{
		"name":    "ccmouse",
		"course":  "golang",
		"site":    "imooc",
		"quality": "notbad",
	}
	// map[]
	m2 := make(map[string]int) // m2 == empty map
	// map[]
	var m3 map[string]int // m3 == nil

	fmt.Println("m, m2, m3:")
	fmt.Println(m, m2, m3)

	name2Everything := make(map[string]interface{})
	name2Everything["xiaoming"] = 2
	name2Everything["xiaohong"] = "girl"
	name2Everything["age"] = 100
	name2Everything["1"] = 100
	name2Everything["2"] = 100
	name2Everything["3"] = 100
	name2Everything["4"] = 100
	name2Everything["5"] = 100
	name2Everything["6"] = 100
	name2Everything["7"] = 100
	name2Everything["8"] = 100

	fmt.Println(name2Everything)
}
