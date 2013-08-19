package ld

import (
	"bytes"
	"math"
	"testing"
)

func TestVcfFileReader(t *testing.T) {
	expected := []*Variant{&Variant{9411243, 19161214, []uint32{0, 0}},
		&Variant{9411245, 18169135, []uint32{0, 0}},
		&Variant{9411254, 18612630, []uint32{0, 0}},
		&Variant{9411618, 14122189, []uint32{256, 0}},
		&Variant{9412099, 14613416, []uint32{0, 0}},
		&Variant{9412126, 14943728, []uint32{uint32(math.Pow(2., 28.)), 0}},
		&Variant{9412339, 19085147, []uint32{0, 0}},
		&Variant{9412503, 7122088, []uint32{4160747375, 0}},
		&Variant{9412603, 0, []uint32{0, 0}},
		&Variant{9412604, 0, []uint32{0, 0}},
	}
	sampleIds := []string{"HG00096", "HG00101", "HG00102", "HG00103", "HG00104", "HG00106", "HG00108", "HG00109", "HG00110", "HG00111", "HG00112", "HG00113", "HG00114", "HG00116", "HG00117", "HG00118", "HG00119"}

	//sampleIds := []string{"HG00096", "HG00099", "HG00108"}
	//sampleIds := getSampleIds

	reader, err := OpenVcfFile("test.vcf")
	if err != nil {
		t.Errorf(err.Error())
	}
	chVariant := CreateVariantChannel(reader, sampleIds)
	var actual []*Variant

	//TODO: add timeout here
	for v := range chVariant {
		actual = append(actual, v)
	}

	if len(actual) != len(expected) {
		t.Errorf("length mismatch. got %v, want %v", actual, expected)
	} else {
		for i, v := range actual {
			if !Equal(v, expected[i]) {
				t.Errorf("at position %v got %v, want %v", i, actual[i], expected[i])
				break
			}
		}
	}
}

func TestGzVcfFileReader(t *testing.T) {
	var (
		part                             []byte
		prefix                           bool
		err                              error
		commentLineCount, totalLineCount int
	)

	reader, err := OpenGzVcfFile("test.vcf.gz")
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
