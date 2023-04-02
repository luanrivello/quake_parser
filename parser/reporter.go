package parser

import (
	"encoding/json"
	"os"
)

func Write(matchs map[string]*Match) {
	writeGroupedInformation(matchs)
}

func writeGroupedInformation(matchs map[string]*Match) {
	writeJsonToFile(matchs, "report")
}

func writeJsonToFile(data interface{}, fileName string) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		println("Error marshalling to JSON:", err)
		return
	}

	file, err := os.Create("report/" + fileName + ".json")
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
//func printReport(matchs map[string]*Match) {
//	for index, match := range matchs {
//		println("-------------------------- " + index + "Report --------------------------")
//		println("TotalKills:", match.TotalKills)
//		for player, kills := range match.KillCount {
//			println(player, kills)
//		}
//	}
//}
