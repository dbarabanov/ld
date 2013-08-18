package ld

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

func readVariants(reader *bufio.Reader, sampleIds []string, ch chan *Variant) {
	var (
		part          []byte
		prefix        bool
		err           error
		sampleIndexes []uint16
	)

	buffer := bytes.NewBuffer(make([]byte, 0))
	for {
		if part, prefix, err = reader.ReadLine(); err != nil {
			break
		}
		buffer.Write(part)
		if !prefix {
			line := buffer.String()
			if line[0] == '#' {
				sampleIndexes = getSampleIndexes(line, sampleIds)
			} else {
				ch <- VcfLineToVariant(line, sampleIndexes)
			}
			buffer.Reset()
		}
	}
	if err == io.EOF {
		close(ch)
	} else {
		panic(err)
	}
}

func getSampleIndexes(line string, SampleIds []string) (sampleIndexes []uint16) {
	return []uint16{9, 10, 11, 12, 13, 14, 15, 16}
}

func VcfLineToVariant(line string, sampleIndexes []uint16) (variant *Variant) {
	if sampleIndexes == nil || len(sampleIndexes) == 0 {
		panic("sampleIndexes not initialized")
	}
	tokens := strings.Split(line, " ")
	if len(tokens) < int(MaxInt(sampleIndexes))+3 {
		panic(fmt.Sprintf("too few tokens in line: %v", line))
	}

	var (
		pos  int
		rsid uint64
		err  error
	)
	if pos, err = strconv.Atoi(tokens[1]); err != nil {
		panic(fmt.Sprintf("bad line: %v", line))
	}
	rsidString := tokens[2]
	if len(rsidString) < 4 || rsidString[0:2] != "rs" {
		panic(fmt.Sprintf("bad rsid in line: %v", line))
	}
	if rsid, err = strconv.ParseUint(rsidString[2:len(tokens[2])-1], 0, 64); err != nil {
		panic(fmt.Sprintf("bad rsid in line: %v", line))
	}

	alleles := make([]uint8, len(sampleIndexes))
	r, _ := regexp.Compile(`^[0,1]\|[0,1]$`) //"0|0", "0|1","1|0", "1|1" 
	for _, index := range sampleIndexes {
		token := tokens[index]
		//fmt.Printf("%v\n", token)
		genotype := strings.Split(token, ":")[0]
		//fmt.Printf("%v\n", genotype)
		if !r.MatchString(genotype) {
			panic(fmt.Sprintf("bad genotype: %v" + genotype))
		}
		if genotype == "0|0" {
			alleles = append(alleles, 0)
		} else if genotype == "0|1" {
			alleles = append(alleles, 1)
		} else if genotype == "1|0" {
			alleles = append(alleles, 2)
		} else if genotype == "1|1" {
			alleles = append(alleles, 3)
		}
	}

	return &Variant{uint32(pos), rsid, nil}
}

func CreateVariantChannel(reader *bufio.Reader, sampleIds []string) chan *Variant {
	ch := make(chan *Variant)
	go readVariants(reader, sampleIds, ch)
	return ch
}
