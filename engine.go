package ld

type Score struct {
	Pos  uint32
	Rsid uint64
	R2   float64
}

type Result struct {
	Variant *Variant
	Scores  []Score
}

type EngineParameters struct {
	WindowSize     uint32
	PopulationSize uint16
	R2Threshold    float64
}

type Engine interface {
	Run(chan *Variant) chan *Result
}

type engine struct {
	params EngineParameters
}

func (e engine) Run(chVariant chan *Variant) chan *Result {
	chResult := make(chan *Result)
	go runEngine(chVariant, chResult)
	return chResult
}

func runEngine(chVariant chan *Variant, chResult chan *Result) {
	for v := range chVariant {
		chResult <- &Result{v, nil}
	}
    close(chResult)
}

func CreateEngine(params EngineParameters) (e Engine, err error) {
	return Engine(&engine{params}), nil
}
