package parser

import (
	"strings"
	"sync"
)

type Match struct {
	totalKills int
	players    []string
	killCount  map[string]int
}

func Parse(log string) {
	var waitgroup sync.WaitGroup

	lines := strings.Split(log, "\n")

	for i, line := range lines {
		line := strings.TrimSpace(line)
		tokens := strings.Split(line, " ")

		if len(tokens) > 2 {
			if tokens[1] == "InitGame:" {
				waitgroup.Add(1)
				go extractMatch(lines, i, waitgroup)
				break
			}
		}
	}

	waitgroup.Wait()
}

func extractMatch(lines []string, lineNumber int, waitgroup sync.WaitGroup) Match {
	// Match data
	var match Match = Match{
		totalKills: 0,
		players:    make([]string, 0),
		killCount:  make(map[string]int),
	}

	for _, line := range lines {
		line := strings.TrimSpace(line)
		tokens := strings.Split(line, " ")
		if len(tokens) > 2 {
			if tokens[1] == "Kill:" {
				match.totalKills++
				println(line)
			} else if tokens[1] == "ShutdownGame:" {
				println(line)
				println(match.totalKills)
				waitgroup.Done()
				return match
			}
		}
	}

	println(match.totalKills)
	return match
}
