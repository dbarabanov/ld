package ld

import (
	"compress/gzip"
	"os"
    "bufio"
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
