package ld

import ("testing"
"bytes")

func TestVcfFileReader(t *testing.T) {
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

	hgIds := []string{"HG00096", "HG00099", "HG00108"}

	reader, err := OpenVcfFile("sample.vcf")
	if err != nil {
		t.Errorf(err.Error())
	}
	chVariant := CreateVariantChannel(reader, hgIds)
	var actual []*Variant

	//TODO: add timeout here
	for v := range chVariant {
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
