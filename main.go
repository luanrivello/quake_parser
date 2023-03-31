package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func main() {

	// Path to the log file
	var path string = getPath()
	println(path)

	var content string = getContent(path)

	// Print the file contents
	fmt.Println(string(content))

}

func getPath() string {

	if len(os.Args) > 1 {
		return strings.Join(os.Args[1:], " ")
	} else {
		return "./data/qgames.log"
	}

}

func getContent(path string) string {

	fileAbsPath, err := filepath.Abs(path)
	if err != nil {
		fmt.Println("Error getting file path:", err)
	}

	// Get file contents
	content, err := ioutil.ReadFile(fileAbsPath)
	if err != nil {
		fmt.Println("Error reading file:", err)
	}

	return string(content)

}
