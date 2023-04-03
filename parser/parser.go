package parser

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
)

// * Match data
type Match struct {
	TotalKills  int            `json:"total_kills"`
	Players     []string       `json:"players"`
	KillCount   map[string]int `json:"kills"`
	Leaderboard map[int]string `json:"player_ranking"`
	KillMeans   map[string]int `json:"kills_by_means"`
}

func NewMatch(matchs map[string]*Match, matchNumber int) *Match {
	var newMatch Match = Match{
		TotalKills:  0,
		Players:     make([]string, 0),
		KillCount:   make(map[string]int),
		Leaderboard: make(map[int]string),
		KillMeans:   make(map[string]int),
	}

	fillKillMeans(&newMatch.KillMeans)

	matchName := fmt.Sprintf("game_%02d", matchNumber)
	matchs[matchName] = &newMatch

	return &newMatch
}

func fillKillMeans(means *map[string]int) {
	(*means)["MOD_UNKNOWN"] = 0
	(*means)["MOD_SHOTGUN"] = 0
	(*means)["MOD_GAUNTLET"] = 0
	(*means)["MOD_MACHINEGUN"] = 0
	(*means)["MOD_GRENADE"] = 0
	(*means)["MOD_GRENADE_SPLASH"] = 0
	(*means)["MOD_ROCKET"] = 0
	(*means)["MOD_ROCKET_SPLASH"] = 0
	(*means)["MOD_PLASMA"] = 0
	(*means)["MOD_PLASMA_SPLASH"] = 0
	(*means)["MOD_RAILGUN"] = 0
	(*means)["MOD_LIGHTNING"] = 0
	(*means)["MOD_BFG"] = 0
	(*means)["MOD_BFG_SPLASH"] = 0
	(*means)["MOD_WATER"] = 0
	(*means)["MOD_SLIME"] = 0
	(*means)["MOD_LAVA"] = 0
	(*means)["MOD_CRUSH"] = 0
	(*means)["MOD_TELEFRAG"] = 0
	(*means)["MOD_FALLING"] = 0
	(*means)["MOD_SUICIDE"] = 0
	(*means)["MOD_TARGET_LASER"] = 0
	(*means)["MOD_TRIGGER_HURT"] = 0
	(*means)["MOD_NAIL"] = 0
	(*means)["MOD_CHAINGUN"] = 0
	(*means)["MOD_PROXIMITY_MINE"] = 0
	(*means)["MOD_KAMIKAZE"] = 0
	(*means)["MOD_JUICED"] = 0
	(*means)["MOD_GRAPPL"] = 0
}

func NewLeaderboard(match *Match) {
	//* Create leaderboard
	leaderboard := match.Leaderboard

	//* Fill leaderboard
	for i := 1; i < len(match.Players)+1; i++ {
		leaderboard[i] = match.Players[i-1]
	}

	//* Order by kills
	n := len(leaderboard)
	for i := 1; i <= n; i++ {
		for j := 1; j <= n-i; j++ {
			if match.KillCount[leaderboard[j]] < match.KillCount[leaderboard[j+1]] {
				leaderboard[j], leaderboard[j+1] = leaderboard[j+1], leaderboard[j]
			}
		}
	}
}

func Parse(log string) map[string]*Match {
	//* Keep track of matches parallel processes
	var waitgroup sync.WaitGroup
	var matchs map[string]*Match = make(map[string]*Match, 0)
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
				matchNumber++
				waitgroup.Add(1)
				newMatch := NewMatch(matchs, matchNumber)

				//* Extract the data in parallel processe
				go extractMatchData(newMatch, lines, lineNumber+1, &waitgroup)
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
			//* Kill log line
			case "Kill:":
				registerKill(match, tokens)

			//* Player log line
			case "ClientUserinfoChanged:":
				registerPlayer(match, tokens)

			//* Another match has started
			case "InitGame:":
				NewLeaderboard(match)
				return
			}
		}

		lineNumber++
	}

	//* End of file
	NewLeaderboard(match)
}

func registerKill(match *Match, tokens []string) {
	//* Add total kills
	match.TotalKills++

	//* Extract killer name
	regex := regexp.MustCompile(`^.* killed`)
	killer := regex.FindString(strings.Join(tokens[5:], " "))
	killer = killer[0 : len(killer)-7]

	//* Extract victim name
	regex = regexp.MustCompile(`killed .* by`)
	victim := regex.FindString(strings.Join(tokens[5:], " "))
	victim = victim[7 : len(victim)-3]

	//* Extract kill mean
	regex = regexp.MustCompile(`by .*`)
	killMean := regex.FindString(strings.Join(tokens[5:], " "))
	killMean = killMean[3:]

	if killer != "<world>" {
		//* Register kill
		match.KillCount[killer]++
	} else {
		//* Subtract kill from the victim of <world>
		match.KillCount[victim]--
	}

	//* Check if it was suicide
	if killer == victim {
		match.KillMeans["MOD_SUICIDE"]++
	} else {
		match.KillMeans[killMean]++
	}
}

func registerPlayer(match *Match, tokens []string) {
	//* Extract Player Name
	regex := regexp.MustCompile(`[^\\n](\w*|\w* )*`)
	player := regex.FindString(strings.Join(tokens[3:], " "))

	if len(player) > 1 {
		//* Register new player
		if contains(match.Players, player) {
			return
		} else {
			match.Players = append(match.Players, player)
			match.KillCount[player] = 0
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
