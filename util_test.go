package ld

import (
	"fmt"
	"math"
	"testing"
)

func TestPanelFileReader(t *testing.T) {
	expected := []string{"HG00096", "HG00099", "HG00171"}
	actual := GetSampleIds("test.panel", "EUR")
	if len(actual) != len(expected) {
		t.Errorf("got %v. want: %v",
			actual, expected)
	} else {
		for i, a := range actual {
			if a != expected[i] {
				t.Errorf("in line %v got %v. want: %v",
					i, actual, expected)
			}
		}
	}
}

func TestGetSampleIndexes(t *testing.T) {
	expected := []uint16{9, 11}
	actual := getSampleIndexes("#CHROM\tPOS\tID\tREF\tALT\tQUAL\tFILTER\tINFO\tFORMAT\tHG00096\tHG00097\tHG00099\tHG00100\tHG00101\tHG00102\tHG00103\tHG00104\tHG00106\tHG0010", []string{"HG00096", "HG00099", "HG00171"}) //17 out of 20 in the test file. just enough to make2 variant uint32-s
	if len(actual) != len(expected) {
		t.Errorf("got %v. want: %v",
			actual, expected)
	} else {
		for i, a := range actual {
			if a != expected[i] {
				t.Errorf("at position %v got %v. want: %v",
					i, actual, expected)
				break
			}
		}
	}
}

func TestPackGenotypes(t *testing.T) {
	expected := []uint32{7, uint32(math.Pow(2., 31.))}
	actual := PackGenotypes([]string{"0|0", "0|0", "0|0", "0|0", "0|0", "0|0", "0|0", "0|0", "0|0", "0|0", "0|0", "0|0", "0|0", "0|0", "0|1", "1|1", "1|0"})
	if len(actual) != len(expected) {
		t.Errorf("wrong length. got %v. want: %v",
			actual, expected)
	} else {
		for i, a := range actual {
			if a != expected[i] {
				t.Errorf("got %v. want: %v",
					actual, expected)
				break
			}
		}
	}
}

func TestUnPackGenotypes(t *testing.T) {
	actual, err := UnpackGenotypes([]uint32{7, 2}, 100)
	//expected_error := errors.New("len(compressed) too short")
	expected_error := fmt.Errorf("len(compressed) too short")
	//if err != expected_error {
	if err == nil || err.Error() != expected_error.Error() {
		t.Errorf("got error %v. want error: %v", err, expected_error)
	}
	actual, err = UnpackGenotypes([]uint32{7, uint32(math.Pow(2., 31.))}, 17)
	expected := []string{"0|0", "0|0", "0|0", "0|0", "0|0", "0|0", "0|0", "0|0", "0|0", "0|0", "0|0", "0|0", "0|0", "0|0", "0|1", "1|1", "1|0"}
	if len(actual) != len(expected) {
		t.Errorf("wrong length. got %v. want: %v",
			actual, expected)
	} else {
		for i, a := range actual {
			if a != expected[i] {
				t.Errorf("got %v. want: %v",
					actual, expected)
				break
			}
		}
	}
}
func TestGenotypeToBits(t *testing.T) {
	actual := genotypeToBits("0|0")
	var expected uint32 = 0
	if actual != expected {
		t.Errorf("got %v. want %v", actual, expected)
	}
}

func TestBitsToGenotype(t *testing.T) {
	actual := bitsToGenotype(1)
	expected := ("0|1")
	if actual != expected {
		t.Errorf("got %v. want %v", actual, expected)
	}
}
