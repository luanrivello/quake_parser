package main

import (
	"os"
	"strings"
)

func main() {

	// Path to the log file
	var path string
	if len(os.Args) > 1 {
		path = strings.Join(os.Args[1:], " ")
	} else {
		path = "./data/qgames.log"
	}

	println(path)

}
