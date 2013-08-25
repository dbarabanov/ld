package ld

type Variant struct {
	Pos       uint32
	Rsid      uint64
	Genotypes []uint32
}

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
	NumWorkers     uint16
}

type Engine interface {
	Run(chan *Variant) chan *Result
}

type ResultWriter interface {
WriteResults(chan *Result)
}
