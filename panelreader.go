package ld

import (
	"io/ioutil"
	"strings"
)

func GetSampleIds(panelFilePath string, population string) []string {
	var (
		content []byte
		err     error
	)
	if content, err = ioutil.ReadFile(panelFilePath); err != nil {
		panic(err)
	}
	lines := strings.Split(string(content), "\n")

	populations := make(map[string][]string)
	for _, line := range lines {
		tokens := strings.Split(line, "\t")
		if len(tokens) > 2 {
			populations[tokens[2]] = append(populations[tokens[2]], tokens[0])
		}
	}

	return populations[population]
}
