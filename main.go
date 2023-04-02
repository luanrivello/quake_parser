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

	println("---------- Quake Log Parser ----------")
	//* Path to the log file
	var path string = getPath()

	//* Get contents of the log file
	var content string = getContent(path)

	//* Extract data
	println("Extracting log data...")
	var matchs map[string]*parser.Match = parser.Parse(content)

	//* Write report
	println("Writing report to ./report/report.json")
	parser.Write(matchs)

	println("-------------- Finished --------------")
}

func getPath() string {
	var path string

	//* Get path from args or default path
	if len(os.Args) > 1 {
		path = strings.Join(os.Args[1:], " ")
	} else {
		path = "./data/qgames.log"
	}

	println("Reading log file at", path)

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
