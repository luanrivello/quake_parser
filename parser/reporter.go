package parser

import (
	"encoding/json"
	"os"
)

func Write(matchs []*Match) {
	createJsonReport(matchs)
}

func createJsonReport(matchs []*Match) {
	jsonData, err := json.MarshalIndent(matchs, " ", " ")
	if err != nil {
		println("Error marshalling to JSON:", err)
		return
	}

	file, err := os.Create("report/report.json")
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

func createReport(matchs []*Match) {
	for _, match := range matchs {
		println("-------------------------- Match", match.Id, "Report --------------------------")
		println("TotalKills:", match.TotalKills)
		for player, kills := range match.KillCount {
			println(player, kills)
		}
	}
}
