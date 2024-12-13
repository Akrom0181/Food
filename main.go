package main

import (
	"fmt"
)

func main() {

	for i := 1; i <= 100; i++ {
		if i%3 == 0 && i%5 == 0 {
		fmt.Println(i, "😒")
	}
}

	var a int = 5
	var b int = 10

	c := a + b
	fmt.Println(c)
}
