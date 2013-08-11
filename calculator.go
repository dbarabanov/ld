package ld

type Score struct {
	pos uint32
	rs  uint64
	r2  float64
}

type CalculatorProperties struct {
	windowSize     uint32
	populationSize uint16
	r2Threshold    float64
}

type Calculator interface {
	Calculate(chan Variant) chan []Score
}

type calculator struct {
	prop CalculatorProperties
}

func (c calculator) Calculate(chan Variant) chan []Score {
	return nil
}

func NewCalculator(prop CalculatorProperties) (calc Calculator, err error) {
	return Calculator(&calculator{prop}), nil
}
