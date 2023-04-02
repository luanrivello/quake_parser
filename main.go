package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"quake_parser/parser"
)

func main() {
	println("---------- Quake Log Parser ----------")

	//* Get path to the log file
	path, err := getPath()
	if err != nil {
		fmt.Println("Error getting file path:", err)
		return
	}

	//* Get contents of the log file
	println("Reading log file at", path)
	content, err := getContent(path)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	//* Extract data
	println("Extracting log data...")
	matchs := parser.Parse(content)

	//* Write report
	println("Writing report to ./report/report.json")
	parser.Write(matchs)

	println("-------------- Finished --------------")
}

func getPath() (string, error) {
	var path string

	//* Get path from args or default path
	if len(os.Args) > 1 {
		path = strings.Join(os.Args[1:], " ")
	} else {
		path = "./data/qgames.log"
	}

	return path, nil
}

func getContent(path string) (string, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(content), nil
}
