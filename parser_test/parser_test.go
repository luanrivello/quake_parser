package parser

import (
	"quake_parser/parser"
	"reflect"
	"strings"
	"sync"
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

func TestParse(t *testing.T) {
	//* Create test data
	log := `0:00 InitGame: map: map1
    0:10 ClientUserinfoChanged: 2 n\player1\t\0\model\sarge\hmodel\sarge\c1\4\c2\5\hc\100\w\0\l\0
    0:20 Kill: 2 3 7: player1 killed player2 by MOD_ROCKET
    0:30 ClientUserinfoChanged: 3 n\player2\t\0\model\sarge\hmodel\sarge\c1\4\c2\5\hc\100\w\0\l\0
    0:40 Kill: 3 2 7: player2 killed player1 by MOD_ROCKET
    0:50 ShutdownGame:
    0:51 InitGame: map: map2
    0:60 ClientUserinfoChanged: 2 n\player3\t\0\model\sarge\hmodel\sarge\c1\4\c2\5\hc\100\w\0\l\0
    0:70 Kill: 2 2 7: player3 killed player3 by MOD_ROCKET
    0:80 ShutdownGame:`

	expectedMatches := make(map[string]*parser.Match)
	parser.NewMatch(expectedMatches, 1)
	parser.NewMatch(expectedMatches, 2)

	//* Match 1
	expectedMatches["game_01"].TotalKills = 2
	expectedMatches["game_01"].Players = []string{"player1", "player2"}
	expectedMatches["game_01"].KillCount =
		map[string]int{
			"player1": 1,
			"player2": 1,
		}
	expectedMatches["game_01"].Leaderboard =
		map[int]string{
			1: "player1",
			2: "player2",
		}
	expectedMatches["game_01"].KillMeans["MOD_ROCKET"] = 2

	//* Match 2
	expectedMatches["game_02"].TotalKills = 1
	expectedMatches["game_02"].Players = []string{"player3"}
	expectedMatches["game_02"].KillCount =
		map[string]int{
			"player3": 1,
		}
	expectedMatches["game_02"].Leaderboard =
		map[int]string{
			1: "player3",
		}
	expectedMatches["game_02"].KillMeans["MOD_SUICIDE"] = 1

	//* FUNCTION CALL
	resultMatches := parser.Parse(log)

	//* Check that the number of matches is correct
	if len(resultMatches) != len(expectedMatches) {
		t.Errorf("Unexpected number of matches. Got %v, expected %v", len(resultMatches), len(expectedMatches))
	}

	//* Check that each match has the expected data
	for matchID, expectedMatch := range expectedMatches {
		resultMatch, ok := resultMatches[matchID]
		if !ok {
			t.Errorf("Match with ID %v not found in result map", matchID)
		}

		if !reflect.DeepEqual(expectedMatch, resultMatch) {
			t.Errorf("Unexpected match data for match %v. Got %v, expected %v", matchID, resultMatch, expectedMatch)
		}
	}
}

func TestExtractMatchData(t *testing.T) {
	//* Create test data
	match := &parser.Match{
		Players:     make([]string, 0),
		KillCount:   make(map[string]int),
		Leaderboard: make(map[int]string),
		KillMeans:   make(map[string]int),
	}
	match.KillMeans["MOD_RAILGUN"] = 0

	lines := []string{
		" 0:00 InitGame: started game",
		" 0:01 ClientUserinfoChanged: 1 n\\Player1\\t\\0\\model\\uriel/uriel\\hmodel\\uriel/uriel\\g_redteam\\gib\\g_blueteam\\uriel",
		" 0:01 ClientUserinfoChanged: 2 n\\Player2\\t\\0\\model\\uriel/uriel\\hmodel\\uriel/uriel\\g_redteam\\gib\\g_blueteam\\uriel",
		" 0:02 Kill: 2 1 1: Player2 killed Player1 by MOD_RAILGUN",
		" 0:03 ClientUserinfoChanged: 3 n\\Player3\\t\\0\\model\\sarge/sarge\\hmodel\\sarge/sarge\\g_redteam\\gib\\g_blueteam\\sarge",
		" 0:04 Kill: 3 1 1: Player3 killed Player1 by MOD_RAILGUN",
		" 0:05 ClientUserinfoChanged: 4 n\\Player4\\t\\0\\model\\sarge/sarge\\hmodel\\sarge/sarge\\g_redteam\\gib\\g_blueteam\\sarge",
		" 0:06 Kill: 4 2 1: Player4 killed Player2 by MOD_RAILGUN",
		" 0:07 ClientUserinfoChanged: 5 n\\Player5\\t\\0\\model\\sarge/sarge\\hmodel\\sarge/sarge\\g_redteam\\gib\\g_blueteam\\sarge",
		" 0:08 Kill: 4 5 1: Player4 killed Player5 by MOD_RAILGUN",
	}
	waitgroup := &sync.WaitGroup{}
	waitgroup.Add(1)

	//* FUNCTION CALL
	parser.ExtractMatchData(match, lines, 1, waitgroup)

	//* Compare results
	if match.TotalKills != 4 {
		t.Errorf("Expected TotalKills to be 4, but got %d", match.TotalKills)
	}
	expectedPlayers := []string{"Player1", "Player2", "Player3", "Player4", "Player5"}
	if !reflect.DeepEqual(match.Players, expectedPlayers) {
		t.Errorf("Expected Players to be %v, but got %v", expectedPlayers, match.Players)
	}
	expectedKills := map[string]int{"Player1": 0, "Player2": 1, "Player3": 1, "Player4": 2, "Player5": 0}
	if !reflect.DeepEqual(match.KillCount, expectedKills) {
		t.Errorf("Expected KillCount to be %v, but got %v", expectedKills, match.KillCount)
	}
	expectedLeaderboard := map[int]string{1: "Player4", 2: "Player2", 3: "Player3", 4: "Player1", 5: "Player5"}
	if !reflect.DeepEqual(match.Leaderboard, expectedLeaderboard) {
		t.Errorf("Expected Leaderboard to be %v, but got %v", expectedLeaderboard, match.Leaderboard)
	}
	expectedKillMeans := map[string]int{"MOD_RAILGUN": 4}
	if !reflect.DeepEqual(match.KillMeans, expectedKillMeans) {
		t.Errorf("Expected KillMeans to be %v, but got %v", expectedKillMeans, match.KillMeans)
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
			KillMeans:   map[string]int{"MOD_RAILGUN": 0},
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
			KillMeans:   map[string]int{"MOD_FALLING": 0},
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
	
	t.Run("Test unknown kill", func(t *testing.T) {
		//* Create test data
		match := &parser.Match{
			TotalKills:  0,
			Players:     []string{"Player1", "Player2", "Player3"},
			KillCount:   map[string]int{},
			Leaderboard: map[int]string{},
			KillMeans:   map[string]int{},
		}

		line := "22:18 Kill: 2 2 7: Player2 killed Player1 by MOD_ABCDEF"
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
		if match.KillMeans["MOD_UNKNOWN"] != 1 {
			t.Errorf("registerKill did not register the correct kill mean. Got %v, expected %v", match.KillMeans["MOD_RIFLE"], 1)
		}
	})
	
}

func TestRegisterPlayer(t *testing.T) {
	//* Create test data
	match := &parser.Match{
		Players:     []string{},
		KillCount:   map[string]int{},
		Leaderboard: map[int]string{},
		KillMeans:   map[string]int{},
	}

	line := "23:04 ClientUserinfoChanged: 2 n\\TestPlayer\\t\\0\\model\\sarge\\hmodel\\sarge\\g_redteam\\none\\g_blueteam\\red"
	tokens := strings.Split(line, " ")

	t.Run("Test register player", func(t *testing.T) {
		//* FUNCTION CALL
		parser.RegisterPlayer(match, tokens)

		//* Assert player has been registered
		expectedPlayers := []string{"TestPlayer"}
		if !reflect.DeepEqual(match.Players, expectedPlayers) {
			t.Errorf("Expected players: %v, but got: %v", expectedPlayers, match.Players)
		}

		//* Assert player kill count is 0
		expectedKillCount := 0
		if match.KillCount["TestPlayer"] != expectedKillCount {
			t.Errorf("Expected kill count for player TestPlayer: %d, but got: %d", expectedKillCount, match.KillCount["TestPlayer"])
		}
	})

	t.Run("Test register existing player", func(t *testing.T) {
		//* FUNCTION CALL
		parser.RegisterPlayer(match, tokens)

		//* Assert player has not been registered again
		if len(match.Players) != 1 {
			t.Errorf("Expected only one player to be registered, but got: %d", len(match.Players))
		}
	})
}

func TestContains(t *testing.T) {
	array := []string{"foo", "bar", "baz"}

	t.Run("Test with existing value", func(t *testing.T) {
		if !parser.Contains(array, "foo") {
			t.Errorf("contains failed: expected true, got false")
		}
	})

	t.Run("Test with non-existing value", func(t *testing.T) {
		if parser.Contains(array, "qux") {
			t.Errorf("contains failed: expected false, got true")
		}
	})
}
