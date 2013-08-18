package ld

import (
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
	actual := getSampleIndexes("#CHROM POS ID REF ALT QUAL FILTER INFO FORMAT HG00096 HG00097 HG00099 HG00100 HG00101 HG00102 HG00103 HG00104 HG00106 HG00108", []string{"HG00096", "HG00099", "HG00171"})
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
