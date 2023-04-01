package parser

func Write(matchs []*Match) {
	createReport(matchs)
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
