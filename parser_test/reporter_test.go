package parser_test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"quake_parser/parser"
)

func TestWrite(t *testing.T) {
	//* Create test data
	match := &parser.Match{
		TotalKills: 10,
		Players:    []string{"Player1", "Player2"},
		KillCount: map[string]int{
			"Player1": 5,
			"Player2": 3,
		},
		Leaderboard: map[int]string{
			1: "Player1",
			2: "Player2",
		},
		KillMeans: map[string]int{
			"MOD_SHOTGUN": 5,
			"MOD_ROCKET":  3,
		},
	}

	//* Create a test matches map
	matches := make(map[string]*parser.Match)
	matches["test_match"] = match

	//* FUNCTION CALL
	parser.Write(matches)

	//* Read the contents of the temporary file
	tempPath := "./report/report.json"
	fileData, err := ioutil.ReadFile(tempPath)
	if err != nil {
		t.Fatalf("Error reading temp file: %v", err)
	}

	//* Unmarshal the JSON data
	var parsedData map[string]*parser.Match
	err = json.Unmarshal(fileData, &parsedData)
	if err != nil {
		t.Fatalf("Error unmarshalling JSON data: %v", err)
	}

	//* Compare the parsed data to the original matches
	if !reflect.DeepEqual(parsedData, matches) {
		t.Fatalf("Parsed data does not match original matches:\n%v\n%v", parsedData, matches)
	}

	//* Remove test report
	os.Remove(tempPath)
}
