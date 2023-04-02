package main

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestMain(t *testing.T) {
	os.Args = []string{"program"}
	os.Remove("./report/report.json")

	//* Run main function
	main()

	//* Verify report file was written
	if _, err := os.Stat("./report/report.json"); os.IsNotExist(err) {
		t.Fatal("Report file was not created")
	}
}

func TestGetPath(t *testing.T) {
	t.Run("Test with no args", func(t *testing.T) {
		expectedPath := "/data/qgames.log"
		os.Args = []string{"program"}

		actualPath, err := getPath()

		actualPath = actualPath[len(actualPath)-16:]
		if err != nil {
			t.Errorf("getPath() error = %v", err)
		} else if actualPath != expectedPath {
			t.Errorf("Expected '%s', but got '%s'", expectedPath, actualPath)
		}
	})

	t.Run("Test with args", func(t *testing.T) {
		expectedPath := "testpath.log"
		os.Args = []string{"program", expectedPath}

		actualPath, err := getPath()

		actualPath = actualPath[len(actualPath)-12:]
		if err != nil {
			t.Errorf("getPath() error = %v", err)
		} else if actualPath != expectedPath {
			t.Errorf("Expected '%s', but got '%s'", expectedPath, actualPath)
		}
	})
}

func TestGetContent(t *testing.T) {
	t.Run("Test with valid path", func(t *testing.T) {
		//* Mock input file
		err := ioutil.WriteFile("./data/test.log", []byte("test log file"), 0644)
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove("./data/test.log")

		expectedContent := "test log file"
		actualContent, err := getContent("./data/test.log")
		if err != nil {
			t.Errorf("getContent() error = %v", err)
		} else if actualContent != expectedContent {
			t.Errorf("Expected '%s', but got '%s'", expectedContent, actualContent)
		}
	})

	t.Run("Test with invalid path", func(t *testing.T) {
		content, err := getContent("./invalid/path.log")
		if err == nil {
			t.Errorf("Expected an invalid path error, but got '%s'", content)
		}
	})
}
