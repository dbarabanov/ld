package ld

import (
	//"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestFileWriter(t *testing.T) {
	var (
		writer  ResultWriter
		err     error
		content []byte
	)
	expected := []string{"POS1\tPOS2\tR^2", "9412603\t9412604\t0.484848", ""}
	filepath := "filewritertest.tmp"
	if writer, err = NewFileWriter(filepath); err != nil {
		panic(err)
	}
	results := mockResultChannel(MockResultDataset1)

	writer.WriteResults(results)
	if content, err = ioutil.ReadFile(filepath); err != nil {
		panic(err)
	}
	lines := strings.Split(string(content), "\n")
	if len(lines) != len(expected) {
		t.Errorf("number of lines writen to %v doesn't match. want %v, got %v\n%v", filepath, len(expected), len(lines), lines)
	} else {
		for i, line := range lines {

			if line != expected[i] {
				t.Errorf("line mismatch at position %v. want %v, got %v\n", i, expected[i], line)
			}
		}
	}
	os.Remove(filepath)
	//for result := range results {
	////fmt.Printf("result: %v\n", result)
	//actual = append(actual, result)
	//}

	//expected := MockResultDataset1
	//if _, err := EqualResults(actual, expected); err != nil {
	//t.Errorf(err.Error())
	//}
}

func mockResultChannel(results []*Result) chan *Result {
	chResult := make(chan *Result)
	go pumpResults(results, chResult)
	return chResult
}

func pumpResults(results []*Result, chResult chan *Result) {
	for _, r := range results {
		chResult <- r
	}
	close(chResult)
}
