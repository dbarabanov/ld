package ld

import (
	"bufio"
	"compress/gzip"
	"io/ioutil"
	"os"
	"strings"
)

func OpenVcfFile(vcfFilePath string) (*bufio.Reader, error) {
	var (
		file *os.File
		err  error
	)
	if file, err = os.Open(vcfFilePath); err != nil {
		return nil, err
	}
	//TODO close the file after all goroutines are done with it
	//defer file.Close()
	return bufio.NewReader(file), nil
}

func OpenGzVcfFile(gzVcfFilePath string) (*bufio.Reader, error) {
	var (
		file     *os.File
		fileGzip *gzip.Reader
		err      error
	)
	if file, err = os.Open(gzVcfFilePath); err != nil {
		return nil, err
	}
	//TODO close the file after all goroutines are done with it
	//defer file.Close()
	if fileGzip, err = gzip.NewReader(file); err != nil {
		return nil, err
	}
	return bufio.NewReader(fileGzip), nil
}

func MaxInt(slice []uint16) uint16 {
	var max uint16
	for _, s := range slice {
		if s > max {
			max = s
		}
	}
	return max
}

func Equal(a, b *Variant) bool {
	if a.pos != b.pos || a.rsid != b.rsid || len(a.minor) != len(b.minor) {
		return false
	}
	for i := range a.minor {
		if a.minor[i] != b.minor[i] {
			return false
		}
	}
	return true
}

func getSampleIndexes(line string, sampleIds []string) (sampleIndexes []uint16) {
	//convert slice to map for efficiency
	samples := make(map[string]bool)
	for _, sample := range sampleIds {
		samples[sample] = true
	}

	tokens := strings.Split(line, " ")
	for i, token := range tokens {
		if samples[token] {
			sampleIndexes = append(sampleIndexes, uint16(i))
		}
	}

	return sampleIndexes
}

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

func CompressGenotypes(genotypes []string) []uint32 {
return nil
}
