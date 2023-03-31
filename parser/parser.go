package parser

import (
	"fmt"
	"regexp"
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

	//* Iterate log lines
	for lineNumber, line := range lines {
		line := strings.TrimSpace(line)
		tokens := strings.Split(line, " ")

		//* Find next match
		if len(tokens) > 2 {
			if tokens[1] == "InitGame:" {
				waitgroup.Add(1)
				matchNumber++
				go extractMatch(lines, lineNumber+1, matchNumber, &waitgroup)
				if matchNumber == 2 {
					break
				}
			}
		}
	}

	//* Wait processes
	waitgroup.Wait()
}

func extractMatch(lines []string, lineNumber int, matchNumber int, waitgroup *sync.WaitGroup) Match {
	defer waitgroup.Done()

	//* Match data
	var match Match = Match{
		id:         matchNumber,
		totalKills: 0,
		players:    make([]string, 0),
		killCount:  make(map[string]int),
	}

	//* Log lines form specific match
	for lineNumber < len(lines) {
		line := strings.TrimSpace(lines[lineNumber])
		tokens := strings.Split(line, " ")
		if len(tokens) > 1 {
			//* Kill data
			if tokens[1] == "Kill:" {
				match.totalKills++

				//* Player name
			} else if tokens[1] == "ClientUserinfoChanged:" {
				//fmt.Println(tokens[3])
				regex := regexp.MustCompile(`\\(\w+|\w+ )+\\`) 
				player := regex.FindString(strings.Join(tokens[3:], " "))
				if len(player) > 1 {
					player := strings.ReplaceAll(player, "\\", "")
					println(player)
				} else {
					fmt.Println("No match found")
				}
				// extract name
				// check if exists

				//* End of match
			} else if tokens[1] == "InitGame:" {
				println("Match", match.id, "TotalKills", match.totalKills)
				return match
			}
		}
		lineNumber++
	}

	println("Match", match.id, "TotalKills", match.totalKills)
	return match
}
