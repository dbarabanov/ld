package ld

import (
	//"fmt"
	"strconv"
	"testing"
)

var MockVariantDataset1 = []*Variant{&Variant{9411243, 19161214, []uint32{0, 0}},
	&Variant{9411245, 18169135, []uint32{0, 0}},
	&Variant{9411254, 18612630, []uint32{0, 0}},
	&Variant{9411618, 14122189, []uint32{256, 0}},
	&Variant{9412099, 14613416, []uint32{0, 0}},
	&Variant{9412126, 14943728, []uint32{268435456, 0}},
	&Variant{9412339, 19085147, []uint32{0, 0}},
	&Variant{9412503, 7122088, []uint32{4160747375, 0}},
	&Variant{9412603, 0, []uint32{0, 0}},
	&Variant{9412604, 0, []uint32{0, 0}},
}

var MockResultDataset1 = []*Result{&Result{&Variant{9411243, 19161214, []uint32{0, 0}}, nil},
	&Result{&Variant{9411245, 18169135, []uint32{0, 0}}, nil},
	&Result{&Variant{9411254, 18612630, []uint32{0, 0}}, nil},
	&Result{&Variant{9411618, 14122189, []uint32{256, 0}}, nil},
	&Result{&Variant{9412099, 14613416, []uint32{0, 0}}, nil},
	&Result{&Variant{9412126, 14943728, []uint32{268435456, 0}}, nil},
	&Result{&Variant{9412339, 19085147, []uint32{0, 0}}, nil},
	&Result{&Variant{9412503, 7122088, []uint32{4160747375, 0}}, nil},
	&Result{&Variant{9412603, 0, []uint32{0, 0}}, nil},
	&Result{&Variant{9412604, 0, []uint32{0, 0}}, nil},
}

func TestRunEngine(t *testing.T) {
	var (
		engine Engine
		err    error
		actual []*Result
	)
	if engine, err = CreateEngine(EngineParameters{10, 17, 0, 2}); err != nil {
		panic(err)
	}
	chVariant := mockVariantChannel(MockVariantDataset1)
	results := engine.Run(chVariant)
	for result := range results {
		//fmt.Printf("result: %v\n", result)
		actual = append(actual, result)
	}

	expected := MockResultDataset1
	if _, err := EqualResults(actual, expected); err != nil {
		t.Errorf(err.Error())
	}
}

func mockVariantChannel(variants []*Variant) chan *Variant {
	chVariant := make(chan *Variant)
	go pumpChannel(variants, chVariant)
	return chVariant
}

func pumpChannel(variants []*Variant, chVariant chan *Variant) {
	for _, v := range variants {
		chVariant <- v
	}
	close(chVariant)
}

func TestComputeR2(t *testing.T) {
	a := []uint32{toInt("1001"), toInt("1")}
	b := []uint32{toInt("1010"), toInt("10")}
	//fmt.Println(toString(a))
	//fmt.Println(toString(b))
	c := make([]uint32, len(a), len(b))
	for i := range a {
		c[i] = a[i] & b[i]
	}
	//fmt.Println(toString(c))
}

func toInt(bitString string) uint32 {
	if i, err := strconv.ParseUint(bitString, 2, 32); err != nil {
		panic(err)
	} else {
		return uint32(i)
	}
	return 0
}

func toString(a []uint32) (bitStrings []string) {
	bitStrings = make([]string, len(a), len(a))

	for i := range a {
		bitStrings[i] = strconv.FormatUint(uint64(a[i]), 2)
	}
	return bitStrings
}

func TestBitCountSingle(t *testing.T) {
	expected := uint16(6)
	actual := bitCountSingle(toInt("01011100000000000000001010"))
	if actual != expected {
		t.Errorf("got %v. want %v", actual, expected)
	}
}

func TestBitCount(t *testing.T) {
	expected := uint16(5)
	actual := bitCount([]uint32{toInt("10"), toInt("101"), toInt("1000000000000000000001")})
	if actual != expected {
		t.Errorf("got %v. want %v", actual, expected)
	}
}

func TestUnion(t *testing.T) {
	expected := []uint32{toInt("10"), toInt("101")}
	actual := union([]uint32{toInt("0"), toInt("100")}, []uint32{toInt("10"), toInt("1")})
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
