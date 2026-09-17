package main

import (
	"github.com/redgrittybrick/mbox"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s filename...\n", os.Args[0]) 
	}
	for _, fileName := range os.Args[1:] {
		if !mbox.IsMbox(fileName) {
			fmt.Println(fileName, "is NOT an mbox file.") 
			continue
		}
		msgs, err := mbox.IndexFile(fileName)
		if err != nil {
			fmt.Println("Unable to index", fileName, "because", err)
			continue
		}
		for n, msg := range msgs {
			fmt.Printf("%4d %s\n", n, msg.Subject) 
		}
	}
	fmt.Printf("%d files examined\n", len(os.Args)-1)
}
