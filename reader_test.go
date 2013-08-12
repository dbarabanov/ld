package ld

import "testing"

func TestVcfFileReader(t *testing.T) {
	//expected := []*Variant{&Variant{0, 0, nil}, &Variant{1, 1, nil}}
	expected := []*Variant{&Variant{9411243, 19161214, nil},
		&Variant{9411245, 18169135, nil},
		&Variant{9411254, 18612630, nil},
		&Variant{9411618, 14122189, nil},
		&Variant{9412099, 14613416, nil},
		&Variant{9412126, 14943728, nil},
		&Variant{9412339, 19085147, nil},
		&Variant{9412503, 7122088, nil},
		&Variant{9412603, 14130669, nil},
		&Variant{9412604, 14504037, nil},
	}
	reader, err := NewVcfReader("sample.vcf")
	if err != nil {
		t.Errorf(err.Error())
	}
	var actual []*Variant
	//TODO: add timeout here
	for v := range reader.Read() {
		actual = append(actual, v)
	}

	if len(actual) != len(expected) {
		t.Errorf("got %v, want %v", actual, expected)
	} else {
		for i, v := range actual {
			if !Equal(v, expected[i]) {
				t.Errorf("at position %v got %v, want %v", i, actual[i], expected[i])
				break
			}
		}
	}
}
