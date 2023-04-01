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
	var matchs []*Match = make([]*Match, 0)
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

				//* New Match
				var newMatch Match = Match{
					id:         matchNumber,
					totalKills: 0,
					players:    make([]string, 0),
					killCount:  make(map[string]int),
				}
				matchs = append(matchs, &newMatch)

				go extractMatchData(&newMatch, lines, lineNumber+1, &waitgroup)

				if matchNumber == 3 {
					break
				}
			}
		}
	}

	//* Wait processes
	waitgroup.Wait()

	//* Create report
	createReport(matchs)
}

func createReport(matchs []*Match) {
	for _, match := range matchs {
		println("-------------------------------- Match", match.id, "Report --------------------------------")
		println(strings.Join(match.players[:], ";"))
		println("TotalKills", match.totalKills)
	}
}

func extractMatchData(match *Match, lines []string, lineNumber int, waitgroup *sync.WaitGroup) {
	defer waitgroup.Done()

	//* Log lines form specific match
	for lineNumber < len(lines) {
		line := strings.TrimSpace(lines[lineNumber])
		tokens := strings.Split(line, " ")
		if len(tokens) > 1 {
			switch tokens[1] {
			//* Kill data
			case "Kill:":
				match.totalKills++

			//* Player name
			case "ClientUserinfoChanged:":
				registerPlayer(match, tokens)

			//* End of match
			case "InitGame:":
				return
			}
		}

		lineNumber++
	}

}

func registerPlayer(match *Match, tokens []string) {
	//* Extract Player Name
	regex := regexp.MustCompile(`[^\\n](\w*|\w* )*`)
	player := regex.FindString(strings.Join(tokens[3:], " "))

	if len(player) > 1 {
		//* Register new player
		if contains(match.players, player) {
			return
		} else {
			match.players = append(match.players, player)
		}
	} else {
		fmt.Println("No match found")
	}
}

func contains(array []string, find string) bool {
	for _, aux := range array {
		if aux == find {
			return true
		}
	}
	return false
}
