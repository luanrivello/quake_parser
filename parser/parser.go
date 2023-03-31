package parser

import "strings"

type Match struct {
	totalKills int
	players    []string
	killCount  map[string]int
}

func Parse(log string) {
	lines := strings.Split(log, "\n")

	for _, line := range lines {
		println(line)
	}
}
