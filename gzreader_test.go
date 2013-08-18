package ld

import (
	"bytes"
	"testing"
)

func TestGzVcfFileReader(t *testing.T) {
	var (
		part                             []byte
		prefix                           bool
		err                              error
		commentLineCount, totalLineCount int
	)

	reader, err := OpenGzVcfFile("sample.vcf.gz")
	if err != nil {
		t.Errorf(err.Error())
	}

	buffer := bytes.NewBuffer(make([]byte, 0))
	for {
		if part, prefix, err = reader.ReadLine(); err != nil {
			break
		}
		buffer.Write(part)
		if !prefix {
			line := buffer.String()
			if line[0] == '#' {
				commentLineCount++
			}
			totalLineCount++
			buffer.Reset()
		}
	}

	ec, et := 1, 11
	if commentLineCount != ec || totalLineCount != et {
		t.Errorf("got commentLineCount, totalLineCount: %v, %v. want: %v, %v",
			commentLineCount, totalLineCount, ec, et)
	}

}
