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
