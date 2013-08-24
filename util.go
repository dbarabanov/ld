package ld

import (
	"bufio"
	"compress/gzip"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
    "math"
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
	if a.Pos != b.Pos || a.Rsid != b.Rsid || len(a.Genotypes) != len(b.Genotypes) {
		return false
	}
	for i := range a.Genotypes {
		if a.Genotypes[i] != b.Genotypes[i] {
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

	tokens := strings.Split(line, "\t")
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
	num_words := len(genotypes) / GPI
	if len(genotypes)%GPI != 0 {
		num_words++
	}
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

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func EqualResult(a, b *Result) (areEqual bool, err error) {
	if _, err := EqualVariant(a.Variant, b.Variant); err != nil {
		return false, err
	}
	if _, err := EqualScores(a.Scores, b.Scores); err != nil {
		return false, err
	}

	return true, nil
}

func EqualVariant(a, b *Variant) (areEqual bool, err error) {
	//if a.Pos != b.Pos || a.Rsid != b.Rsid || len(a.Genotypes) != len(b.Genotypes) {
	if a.Pos != b.Pos {
		return false, errors.New(fmt.Sprintf("position mismatch: %v!=%v\n", a.Pos))
	} else if a.Rsid != b.Rsid {
		return false, errors.New(fmt.Sprintf("rsid mismatch: %v!=%v\n", a.Rsid, b.Rsid))
	} else if len(a.Genotypes) != len(b.Genotypes) {
		return false, errors.New(fmt.Sprintf("len(genotypes) mismatch: %v!=%v\n", len(a.Genotypes), len(b.Genotypes)))
	}

	for i := range a.Genotypes {
		if a.Genotypes[i] != b.Genotypes[i] {
			return false, errors.New(fmt.Sprintf("genotype mismatch at position %v: %v!=%v\n", i, a.Genotypes[i], b.Genotypes[i]))
		}
	}
	return true, nil
}

func EqualScores(a, b []Score) (areEqual bool, err error) {
	if len(a) != len(b) {
		return false, errors.New(fmt.Sprintf("Scores lengths not equal. %v!=%v\n", len(a), len(b)))
	}
	for i, s := range a {
		//if _, err := EqualScore(b[i], s); err != nil {
		if b[i] != s {
			return false, errors.New(fmt.Sprintf("different scores at position %v: %v!=%v\n", i, b[i], s))
		}
	}
	return true, nil
}

//func EqualScore(a,b Score) (areEqual bool, err error) {
//}

func EqualResults(a, b []*Result) (areEqual bool, err error) {
	if len(a) != len(b) {
		return false, errors.New(fmt.Sprintf("Results lengths not equal. %v!=%v\n", len(a), len(b)))
	}
	for i, r := range a {
		if _, err := EqualResult(b[i], r); err != nil {
			return false, errors.New(fmt.Sprintf("different results at position %v: %v\n", i, err))
		}
	}
	return true, nil
}

func (r Result) String() string {
	return fmt.Sprintf("Result{%v, %v}", r.Variant, r.Scores)
}

func (v Variant) String() string {
	return fmt.Sprintf("Variant{%v, %v, %v}", v.Pos, v.Rsid, v.Genotypes)
}

func round(x float64, prec int) float64 {
	var rounder float64
	pow := math.Pow(10, float64(prec))
	intermed := x * pow
	_, frac := math.Modf(intermed)

	if frac >= 0.5 {
		rounder = math.Ceil(intermed)
	} else {
		rounder = math.Floor(intermed)
	}

	return rounder / pow
}
