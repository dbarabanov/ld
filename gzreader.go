package ld

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
)

type gzVcfFileReader struct {
	prop          VariantProperties
	gzVcfFilePath string
}

func (r gzVcfFileReader) Read() chan *Variant {
	ch := make(chan *Variant)
	go readGzVcfFile(r.gzVcfFilePath, ch)
	return ch
}

func readGzVcfFile(gzVcfFilePath string, ch chan *Variant) {
	var (
		file   *os.File
		part   []byte
		prefix bool
		err    error
	)
	if file, err = os.Open(gzVcfFilePath); err != nil {
		panic(fmt.Sprintf("error reading gz vcf file: %v", err.Error()))
		close(ch)
		return
	}
	defer file.Close()

	//	reader := bufio.NewReader(file)
	fileGzip, err := gzip.NewReader(file)
	if err != nil {
		panic(fmt.Sprintf("The file %v is not in gzip format.\n", gzVcfFilePath))
		close(ch)
		return
	}
	reader := bufio.NewReader(fileGzip)

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

func NewGzVcfReader(gzVcfFilePath string) (v VariantReader, err error) {
	return VariantReader(&gzVcfFileReader{VariantProperties{"", 0}, gzVcfFilePath}), nil
}
