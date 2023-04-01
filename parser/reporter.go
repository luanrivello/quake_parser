package parser

import (
	"encoding/json"
	"os"
)

func Write(matchs []*Match) {
	createJsonReport(matchs)
}

func createJsonReport(matchs []*Match) {
	jsonData, err := json.Marshal(matchs[0])
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
		println("-------------------------- Match", match.id, "Report --------------------------")
		println("TotalKills:", match.totalKills)
		for player, kills := range match.killCount {
			println(player, kills)
		}
	}
}
