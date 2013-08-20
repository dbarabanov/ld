package ld

import (
	//"fmt"
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
	if engine, err = CreateEngine(EngineParameters{10, 17, 0}); err != nil {
		panic(err)
	}
	chVariant := mockVariantChannel(MockVariantDataset1)
	results := engine.Run(chVariant)
	for result := range results {
		//fmt.Printf("result: %v\n", result)
		actual = append(actual, result)
	}

	expected := MockResultDataset1
	if !EqualResults(actual, expected) {
		t.Errorf("got %v. want %v", actual, expected)
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
