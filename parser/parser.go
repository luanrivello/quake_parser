package parser

import (
	"strings"
	"sync"
)

type Match struct {
	id         int
	totalKills int
	players    []string
	killCount  map[string]int
}

func Parse(log string) {
	var waitgroup sync.WaitGroup
	var matchNumber int = 0

	lines := strings.Split(log, "\n")

	for lineNumber, line := range lines {
		line := strings.TrimSpace(line)
		tokens := strings.Split(line, " ")

		if len(tokens) > 2 {
			if tokens[1] == "InitGame:" {
				waitgroup.Add(1)
				matchNumber++
				go extractMatch(lines, lineNumber+1, matchNumber, &waitgroup)
			}
		}
	}

	waitgroup.Wait()
}

func extractMatch(lines []string, lineNumber int, matchNumber int, waitgroup *sync.WaitGroup) Match {
	defer waitgroup.Done()

	// Match data
	var match Match = Match{
		id:         matchNumber,
		totalKills: 0,
		players:    make([]string, 0),
		killCount:  make(map[string]int),
	}

	for lineNumber < len(lines) {
		line := strings.TrimSpace(lines[lineNumber])
		tokens := strings.Split(line, " ")
		if len(tokens) > 1 {
			if tokens[1] == "Kill:" {
				//println(line)
				match.totalKills++
			} else if tokens[1] == "InitGame:" {
				//println(line)
				println("Match", match.id, "TotalKills", match.totalKills)
				return match
			}
		}
		lineNumber++
	}

	return match
}
