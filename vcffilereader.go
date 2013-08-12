package ld

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type vcfFileReader struct {
	prop        VariantProperties
	vcfFilePath string
}

func (r vcfFileReader) Read() chan *Variant {
	ch := make(chan *Variant)
	go readVcfFile(r.vcfFilePath, ch)
	return ch
}

func readVcfFile(vcfFilePath string, ch chan *Variant) {
	var (
		file   *os.File
		part   []byte
		prefix bool
		err    error
	)
	if file, err = os.Open(vcfFilePath); err != nil {
		//		fmt.Printf("error reading vcf file: %v", err.Error())
		panic(fmt.Sprintf("error reading vcf file: %v", err.Error()))
		close(ch)
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	buffer := bytes.NewBuffer(make([]byte, 0))
	for {
		if part, prefix, err = reader.ReadLine(); err != nil {
			break
		}
		buffer.Write(part)
		if !prefix {
			//TODO read hgIndexes from panel file
			hgIndexes := []uint16{9, 10, 11, 12, 13, 14, 15, 16, 17, 18}
			v := VcfLineToVariant(buffer.String(), hgIndexes)
			if v != nil {
				ch <- v
			}
			buffer.Reset()
		}
	}
	if err == io.EOF {
		close(ch)
	} else {
		//		fmt.Printf("error reading vcf file: %v", err.Error())
		panic(fmt.Sprintf("error reading vcf file: %v", err.Error()))
		close(ch)
	}
}

func VcfLineToVariant(line string, hgIndexes []uint16) (variant *Variant) {
	if line[0] == '#' {
		return nil
	}
	tokens := strings.Split(line, " ")
	//fmt.Printf("%v", strings.Split(line, " "))
	var (
		pos  int
		rsid uint64
		err  error
	)
	if pos, err = strconv.Atoi(tokens[1]); err != nil {
		panic(err.Error())
	}
	if rsid, err = strconv.ParseUint(tokens[2][2:len(tokens[2])-1], 0, 64); err != nil {
		panic(err.Error())
	}
	alleles := make([]uint8, len(hgIndexes))
	//	for _, token := range tokens[9:] {
	for _, index := range hgIndexes {
		token := tokens[index]
		//fmt.Printf("%v\n", token)
		genotype := strings.Split(token, ":")[0]
		//fmt.Printf("%v\n", genotype)
		//TODO panic if genotype is unphased
		if genotype == "0|0" {
			alleles = append(alleles, 0)
		} else if genotype == "0|1" {
			alleles = append(alleles, 1)
		} else if genotype == "1|0" {
			alleles = append(alleles, 2)
		} else {
			alleles = append(alleles, 3)
		}
	}

	return &Variant{uint32(pos), rsid, nil}
}

func NewVcfReader(vcfFilePath string) (v VariantReader, err error) {
	return VariantReader(&vcfFileReader{VariantProperties{"", 0}, vcfFilePath}), nil
}

type VariantProperties struct {
	chromosome     string
	populationSize uint16
}
