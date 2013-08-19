package ld

import (
	"bufio"
	"compress/gzip"
	"fmt"
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

const COMPRESSED_GENOTYPE_BIT_LENGTH = 2

const GPI = 32 / COMPRESSED_GENOTYPE_BIT_LENGTH //16. genotypes per unsigned 32 bit integer. every genotype is compressed to 2 bits.

func PackGenotypes(genotypes []string) (compressed []uint32) {
	num_words := len(genotypes)/GPI + len(genotypes)%GPI
	compressed = make([]uint32, num_words, num_words)
	for i, genotype := range genotypes {
		word := compressed[i/GPI]
		posInWord := GPI - i%GPI - 1
		offset := uint32((posInWord) * COMPRESSED_GENOTYPE_BIT_LENGTH)
		word = word | (genotypeToBits(genotype) << offset)
		compressed[i/GPI] = word
	}
	return compressed
}

func UnpackGenotypes(compressed []uint32, length int) (genotypes []string, err error) {
	if len(compressed)*GPI < length {
		return nil, fmt.Errorf("len(compressed) too short")
	}
	genotypes = make([]string, length, length)
	for i := 0; i < length; i++ {
		word := compressed[i/GPI]
		posInWord := GPI - i%GPI - 1
		offset := uint32((posInWord) * COMPRESSED_GENOTYPE_BIT_LENGTH)
		bits := ((3 << offset) & word) >> offset //geting value of 2 bits of interest
		genotypes[i] = bitsToGenotype(bits)
	}
	return genotypes, nil
}

func genotypeToBits(genotype string) uint32 {
	var retval uint32
	if genotype == "0|0" {
		retval = 0
	} else if genotype == "0|1" {
		retval = 1
	} else if genotype == "1|0" {
		retval = 2
	} else if genotype == "1|1" {
		retval = 3
	} else {
		panic("bad genotype: " + genotype)
	}
	return retval
}

func bitsToGenotype(bits uint32) string {
	var retval string
	if bits == 0 {
		retval = "0|0"
	} else if bits == 1 {
		retval = "0|1"
	} else if bits == 2 {
		retval = "1|0"
	} else if bits == 3 {
		retval = "1|1"
	} else {
		panic(fmt.Sprintf("bad genotype bits: %v", bits))
	}
	return retval
}
