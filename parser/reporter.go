package parser

import (
	"encoding/json"
	"os"
)

func Write(matchs map[string]*Match) {
	createJsonReport(matchs)
	createJsonRank(matchs)
}

func createJsonReport(matchs map[string]*Match) {
	jsonData, err := json.MarshalIndent(matchs, "", "  ")
	if err != nil {
		println("Error marshalling to JSON:", err)
		return
	}

	file, err := os.Create("report/grouped_information.json")
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

type rank struct {
	leaderboard map[int]string
}

func newRank(match Match) rank {
	var aux int = 1
	var result rank = rank{
		leaderboard: make(map[int]string),
	}

	for player, _ := range match.KillCount {
		result.leaderboard[aux] = player
	}

	return result
}

func createJsonRank(matchs map[string]*Match) {
	var ranks map[string]rank = make(map[string]rank)

	for matchName, match := range matchs {
		ranks[matchName] = newRank(*match)
	}

	jsonData, err := json.MarshalIndent(ranks, "", "  ")
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
