package parser

import (
	"encoding/json"
	"os"
)

func Write(matchs map[string]*Match) {
	writeGroupedInformation(matchs)
	writePlayerRanking(matchs)
}

func writeGroupedInformation(matchs map[string]*Match) {
	writeJsonToFile(matchs)
}

type Rank struct {
	Leaderboard map[int]string `json:"player_ranking"`
}

func newRank(match Match) Rank {
	//* Create rank
	var result Rank = Rank{
		Leaderboard: make(map[int]string),
	}
	leaderboard := result.Leaderboard

	//* Fill ranking
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

	return result
}

func writePlayerRanking(matchs map[string]*Match) {
	var ranks map[string]Rank = make(map[string]Rank)

	//* Generate a rank for each match
	for matchName, match := range matchs {
		ranks[matchName] = newRank(*match)
	}

	writeJsonToFile(ranks)
}

func writeJsonToFile(data interface{}) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		println("Error marshalling to JSON:", err)
		return
	}

	file, err := os.Create("report/player_ranking.json")
	if err != nil {
		println("Error creating JSON file:", err)
		return
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		println("Error writing JSON to file:", err)
		return
	}
}

//* Print report
//func createReport(matchs map[string]*Match) {
//	for index, match := range matchs {
//		println("-------------------------- " + index + "Report --------------------------")
//		println("TotalKills:", match.TotalKills)
//		for player, kills := range match.KillCount {
//			println(player, kills)
//		}
//	}
//}
