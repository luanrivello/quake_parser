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

func Parse(log string) []*Match {
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
			}
		}
	}

	//* Wait processes
	waitgroup.Wait()

	return matchs
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
				registerKill(match, tokens)

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

func registerKill(match *Match, tokens []string) {
	match.totalKills++

	regex := regexp.MustCompile(`^.* killed`)
	killer := regex.FindString(strings.Join(tokens[5:], " "))
	killer = killer[0 : len(killer)-7]

	if killer != "<world>" {
		match.killCount[killer]++
	} else {
		regex := regexp.MustCompile(`killed .* by`)
		victim := regex.FindString(strings.Join(tokens[5:], " "))
		victim = victim[7:len(victim)-3]
		match.killCount[victim]--
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
			match.killCount[player] = 0
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
