package main

import (
//	"github.com/redgrittybrick/mbox"
	"bufio"
	"mbox"
	"fmt"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s filename\n", os.Args[0])
		return
	}

	fileName := os.Args[1]
	if !mbox.IsMbox(fileName) {
		fmt.Println(fileName, "is NOT an mbox file.")
		return
	}

	parser := mbox.New()
	msgs, err := parser.IndexFile(fileName)

	if err != nil {
		fmt.Println("Unable to index", fileName, "because", err)
		return
	}

	fmt.Printf("Found %d messages in %s\n", len(msgs), fileName)
	for n, msg := range msgs {
		fmt.Printf("%4d: %-16s  %s -> %s : %s\n",
			n+1, msg.Date(), msg.From(), msg.To(), msg.Subject())
	}

	reader := bufio.NewReader(os.Stdin)
	reply := ""
	for reply != "q" {
		fmt.Print("Enter message number to view? (q to quit): ")
		reply, _ = reader.ReadString('\n')
		reply = reply[:len(reply)-1] // remove trailing CR
		if reply == "q" {
			return
		}
		n, err := strconv.Atoi(reply)
		switch {
		case err != nil :
			fmt.Println(err)
		case n < 1 || n > len(msgs) :
			fmt.Println("Message %d is out of range\n", n)
		default:
			fmt.Println(divider)
			strt, end := msgs[n-1].Start(), msgs[n-1].End()
			raw, err := mbox.GetMsg(fileName, strt, end)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(raw)
			}
			fmt.Println(divider)
		}
	}

}

const divider = "----------------------------------------" +
                "--------------------------------------"
