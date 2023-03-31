package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"quake_parser/parser"
)

func main() {

	//* Path to the log file
	var path string = getPath()
	println(path)

	//* Get contents of the log file
	var content string = getContent(path)

	//* Extract data
	parser.Parse(content)

}

func getPath() string {
	var path string

	//* Get path from args or default path
	if len(os.Args) > 1 {
		path = strings.Join(os.Args[1:], " ")
	} else {
		path = "./data/qgames.log"
	}

	//* Convert to absolute path
	fileAbsPath, err := filepath.Abs(path)
	if err != nil {
		fmt.Println("Error getting file path:", err)
	}

	return fileAbsPath

}

func getContent(path string) string {

	content, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("Error reading file:", err)
	}

	return string(content)

}
