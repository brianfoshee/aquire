package main

//[1 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]

import (
	"fmt"
)

func main() {
	buf := make([]byte, 20)
	buf[0] = 1
	fmt.Println(buf)

	if buf[0] == 1 {
		fmt.Println("Success")
	} else {
		fmt.Println("Error")
	}
}
