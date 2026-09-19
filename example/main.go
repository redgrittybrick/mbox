package main

import (
//	"github.com/redgrittybrick/mbox"
	"mbox"
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

		parser := mbox.New().LogLevel(0)
		msgs, err := parser.IndexFile(fileName)

		if err != nil {
			fmt.Println("Unable to index", fileName, "because", err)
			continue
		}
		fmt.Printf("Found %d messages in %s\n", len(msgs), fileName)
		for n, msg := range msgs {
			fmt.Printf("%4d: %-16s  %s -> %s : %s\n",
				n+1, msg.Date(), msg.From(), msg.To(), msg.Subject())
		}
	}
	fmt.Printf("Examined %d files\n", len(os.Args)-1)
}
