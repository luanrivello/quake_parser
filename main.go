package main

import (
	"os"
	"strings"
)

func main() {

	// Path to the log file
	var path string = getPath()
	println(path)

}

func getPath() string {

	if len(os.Args) > 1 {
		return strings.Join(os.Args[1:], " ")
	} else {
		return "./data/qgames.log"
	}

}