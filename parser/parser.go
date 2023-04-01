package parser

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
)

// * Match data
type Match struct {
	id         int
	totalKills int
	players    []string
	killCount  map[string]int
}

func Parse(log string) []*Match {
	//* Keep track of matches parallel processes
	var waitgroup sync.WaitGroup
	var matchs []*Match = make([]*Match, 0)
	var matchNumber int = 0

	//* Log lines as array
	var lines []string = strings.Split(log, "\n")

	//* Iterate log lines
	for lineNumber, line := range lines {
		var line string = strings.TrimSpace(line)
		var tokens []string = strings.Split(line, " ")

		//* Find next match
		if len(tokens) > 2 {
			if tokens[1] == "InitGame:" {
				//* New Match
				waitgroup.Add(1)
				matchNumber++
				var newMatch Match = Match{
					id:         matchNumber,
					totalKills: 0,
					players:    make([]string, 0),
					killCount:  make(map[string]int),
				}
				matchs = append(matchs, &newMatch)

				//* Extract the data in parallel processe
				go extractMatchData(&newMatch, lines, lineNumber+1, &waitgroup)
			}
		}
	}

	//* Wait parallel processes to finish
	waitgroup.Wait()
	return matchs
}

func extractMatchData(match *Match, lines []string, lineNumber int, waitgroup *sync.WaitGroup) {
	defer waitgroup.Done()

	//* Iterate log lines of specific match
	for lineNumber < len(lines) {
		var line string = strings.TrimSpace(lines[lineNumber])
		var tokens []string = strings.Split(line, " ")

		if len(tokens) > 1 {
			switch tokens[1] {
			//* Kill log
			case "Kill:":
				registerKill(match, tokens)

			//* Player log
			case "ClientUserinfoChanged:":
				registerPlayer(match, tokens)

			//* Another match is starting
			case "InitGame:":
				return
			}
		}

		lineNumber++
	}
}

func registerKill(match *Match, tokens []string) {
	//* Add total kills
	match.totalKills++

	//* Extract killer name
	regex := regexp.MustCompile(`^.* killed`)
	killer := regex.FindString(strings.Join(tokens[5:], " "))
	killer = killer[0 : len(killer)-7]

	//* If killer was not <world>
	if killer != "<world>" {
		//* Register kill
		match.killCount[killer]++
	} else {
		//* Extract victims name
		regex := regexp.MustCompile(`killed .* by`)
		victim := regex.FindString(strings.Join(tokens[5:], " "))
		victim = victim[7 : len(victim)-3]

		//* Subtract kill from the victim of <world>
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
