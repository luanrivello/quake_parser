package parser

import (
	"quake_parser/parser"
	"reflect"
	"strings"
	"testing"
)

func TestNewLeaderboard(t *testing.T) {
	//* Create test data
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

	//* FUNCTION CALL
	parser.NewLeaderboard(match)

	//* Check leaderboard is ranked correctly
	expectedLeaderboard := map[int]string{
		1: "Player3",
		2: "Player1",
		3: "Player2",
	}

	if !reflect.DeepEqual(match.Leaderboard, expectedLeaderboard) {
		t.Errorf("Expected leaderboard %v, but got %v", expectedLeaderboard, match.Leaderboard)
	}
}

func TestRegisterKill(t *testing.T) {
	t.Run("Test player kill", func(t *testing.T) {
		//* Create test data
		match := &parser.Match{
			TotalKills:  0,
			Players:     []string{"Player1", "Player2", "Player3"},
			KillCount:   map[string]int{},
			Leaderboard: map[int]string{},
			KillMeans:   map[string]int{},
		}

		line := "22:18 Kill: 2 2 7: Player1 killed Player2 by MOD_RAILGUN"
		tokens := strings.Split(line, " ")

		//* FUNCTION CALL
		parser.RegisterKill(match, tokens)

		//* Check if total kills increased
		if match.TotalKills != 1 {
			t.Errorf("registerKill did not increment TotalKills. Got %v, expected %v", match.TotalKills, 1)
		}

		//* Check if killer's kill count increased
		if match.KillCount["Player1"] != 1 {
			t.Errorf("registerKill did not register a kill for the killer. Got %v, expected %v", match.KillCount["Player1"], 1)
		}

		//* Check if victim's kill count remained the same
		if match.KillCount["Player2"] != 0 {
			t.Errorf("registerKill subtracted a kill from <world> instead of adding it to the victim. Got %v, expected %v", match.KillCount["Player2"], 0)
		}

		//* Check if kill means were updated
		if match.KillMeans["MOD_RAILGUN"] != 1 {
			t.Errorf("registerKill did not register the correct kill mean. Got %v, expected %v", match.KillMeans["MOD_RIFLE"], 1)
		}
	})

	t.Run("Test world kill", func(t *testing.T) {
		//* Create test data
		match := &parser.Match{
			TotalKills:  0,
			Players:     []string{"Player1", "Player2", "Player3"},
			KillCount:   map[string]int{},
			Leaderboard: map[int]string{},
			KillMeans:   map[string]int{},
		}

		line := "22:18 Kill: 2 2 7: <world> killed Player2 by MOD_FALLING"
		tokens := strings.Split(line, " ")

		//* FUNCTION CALL
		parser.RegisterKill(match, tokens)

		//* Check if total kills increased
		if match.TotalKills != 1 {
			t.Errorf("registerKill did not increment TotalKills. Got %v, expected %v", match.TotalKills, 1)
		}

		//* Check if victim's kill count decreased
		if match.KillCount["Player2"] != -1 {
			t.Errorf("registerKill didnt subtracted a kill from <world>. Got %v, expected %v", match.KillCount["Player2"], 0)
		}

		//* Check if kill means were updated
		if match.KillMeans["MOD_FALLING"] != 1 {
			t.Errorf("registerKill did not register the correct kill mean. Got %v, expected %v", match.KillMeans["MOD_RIFLE"], 1)
		}
	})

	t.Run("Test suicide kill", func(t *testing.T) {
		//* Create test data
		match := &parser.Match{
			TotalKills:  0,
			Players:     []string{"Player1", "Player2", "Player3"},
			KillCount:   map[string]int{},
			Leaderboard: map[int]string{},
			KillMeans:   map[string]int{},
		}

		line := "22:18 Kill: 2 2 7: Player2 killed Player2 by MOD_ROCKET_SPLASH"
		tokens := strings.Split(line, " ")

		//* FUNCTION CALL
		parser.RegisterKill(match, tokens)

		//* Check if total kills increased
		if match.TotalKills != 1 {
			t.Errorf("registerKill did not increment TotalKills. Got %v, expected %v", match.TotalKills, 1)
		}

		//* Check if killer's kill count increased
		if match.KillCount["Player2"] != 1 {
			t.Errorf("registerKill did not register a kill for the killer. Got %v, expected %v", match.KillCount["Player1"], 1)
		}

		//* Check if kill means were updated
		if match.KillMeans["MOD_SUICIDE"] != 1 {
			t.Errorf("registerKill did not register the correct kill mean. Got %v, expected %v", match.KillMeans["MOD_RIFLE"], 1)
		}
	})
}
