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
				go extractMatch(lines, i, &waitgroup)
			}
		}
	}

	waitgroup.Wait()
}

func extractMatch(lines []string, lineNumber int, waitgroup *sync.WaitGroup) Match {
	defer waitgroup.Done()

	// Match data
	var match Match = Match{
		totalKills: 0,
		players:    make([]string, 0),
		killCount:  make(map[string]int),
	}

	for {
		line := strings.TrimSpace(lines[lineNumber])
		tokens := strings.Split(line, " ")
		if len(tokens) > 1 {
			if tokens[1] == "Kill:" {
				match.totalKills++
			} else if tokens[1] == "ShutdownGame:" {
				println(match.totalKills)
				return match
			}
		}
		lineNumber++
	}

}
