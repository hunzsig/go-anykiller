package main

import (
	"fmt"
	"go-anykiller/killer"
	"go-anykiller/lib"
	"os"
	"time"
)

func main() {
	search := ""
	if len(os.Args) == 2 {
		search = os.Args[1]
	} else {
		lib.Panic("error")
	}
	if search != "" {
		t := time.Now()
		fmt.Println("[Search]", search)
		killer.Mikanani(search)
		fmt.Println("[Finish] " + time.Since(t).String())
	}
}
