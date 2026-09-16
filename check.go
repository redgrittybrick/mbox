package mbox

import (
	"log"
)

func check(err error, context string) {
	if err != nil {
		log.Fatal(err.Error() + "[" + context + "]")
	}
}
