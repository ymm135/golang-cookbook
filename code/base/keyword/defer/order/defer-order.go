package main

func main() {
	defer println("A")
	defer println("B")
	defer println("C")
}
