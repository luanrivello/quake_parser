package parser

import (
	"quake_parser/parser"
	"reflect"
	"testing"
)

func TestNewLeaderboard(t *testing.T) {
	match := &parser.Match{
		TotalKills: 8,
		Players:    []string{"Player1", "Player2", "Player3"},
		KillCount: map[string]int{
			"Player1": 2,
			"Player2": 1,
			"Player3": 5,
		},
		Leaderboard: map[int]string{},
		KillMeans:   map[string]int{},
	}

	parser.NewLeaderboard(match)

	expectedLeaderboard := map[int]string{
		1: "Player3",
		2: "Player1",
		3: "Player2",
	}

	if !reflect.DeepEqual(match.Leaderboard, expectedLeaderboard) {
		t.Errorf("Expected leaderboard %v, but got %v", expectedLeaderboard, match.Leaderboard)
	}
}
