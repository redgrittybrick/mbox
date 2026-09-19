package mbox

import (
	"fmt"
	"os"
)

func check(err error, context string) {
	if err != nil {
		fmt.Println(err.Error() + "[" + context + "]")
		os.Exit(9)
	}
}
