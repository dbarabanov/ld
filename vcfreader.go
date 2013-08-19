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
				if line[1] != '#' {
					sampleIndexes = getSampleIndexes(line, sampleIds)
				}
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

func VcfLineToVariant(line string, sampleIndexes []uint16) (variant *Variant) {
	if sampleIndexes == nil || len(sampleIndexes) == 0 {
		panic(fmt.Sprintf("sampleIndexes not initialized: %v", sampleIndexes))
	}
	tokens := strings.Split(line, "\t")
	if len(tokens) < int(MaxInt(sampleIndexes)) {
		panic(fmt.Sprintf("too few tokens(%v) in line: %v", len(tokens), line))
	}

	var (
		pos  int
		rsid uint64
		err  error
	)
	if pos, err = strconv.Atoi(tokens[1]); err != nil {
		panic(fmt.Sprintf("bad pos in line: %v", line[0:min(30, len(line))]))
	}
	rsidString := tokens[2]
    if rsidString == "."{
    rsid = 0} else {
	if len(rsidString) < 4 || rsidString[0:2] != "rs" {
		panic(fmt.Sprintf("bad rsid in line: %v", line[0:min(30, len(line))]))
	}
	if rsid, err = strconv.ParseUint(rsidString[2:len(tokens[2])-1], 0, 64); err != nil {
		panic(fmt.Sprintf("bad rsid in line: %v", line[0:min(30, len(line))]))
	}
    }

	genotypes := make([]string, len(sampleIndexes))
	r, _ := regexp.Compile(`^[0,1]\|[0,1]$`) //"0|0", "0|1","1|0", "1|1" 
	for i, index := range sampleIndexes {
		token := tokens[index]
		//fmt.Printf("%v\n", token)
		genotype := strings.Split(token, ":")[0]
		//fmt.Printf("%v\n", genotype)
		if !r.MatchString(genotype) {
			panic(fmt.Sprintf("bad genotype: %v in line: %v", genotype, line))
		}
		genotypes[i] = genotype
	}
	//fmt.Printf("pos: %v, rsid: %v, genotypes: %v\n", pos, rsid, genotypes)

	return &Variant{uint32(pos), rsid, PackGenotypes(genotypes)}
}

func CreateVariantChannel(reader *bufio.Reader, sampleIds []string) chan *Variant {
	ch := make(chan *Variant)
	go readVariants(reader, sampleIds, ch)
	return ch
}
