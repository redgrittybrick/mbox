package main

import (
	"github.com/redgrittybrick/mbox"
	"fmt"
)

func main() {
	fmt.Println("start")
	if mbox.IsMbox("main.go") {
		fmt.Println("OK") 
	} else {
		fmt.Println("FAIL") 
	}
	fmt.Println("end")
}
