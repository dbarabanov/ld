package ld

import (
//"fmt"
//"math"
//"strconv"
)

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

type engine struct {
	params EngineParameters
}

type variantList struct {
	index   uint32
	variant *Variant
	next    *variantList
}

func (e engine) Run(chVariant chan *Variant) chan *Result {
	chResult := make(chan *Result, 5)
	go runEngine(chVariant, chResult, e.params)
	return chResult
}

func runEngine(chVariant chan *Variant, chResult chan *Result, params EngineParameters) {
	//var workQueue chan *variantList
	//TODO expose queue size somewhere (instead of hardcoding "2")
	//workQueue := make(chan *variantList, params.NumWorkers*2) //make queue a little larger so that workers always have work.
	reorderQueue := startReorderQueue(params.NumWorkers, chResult)
	workQueue := startWorkers(params.NumWorkers, reorderQueue)
	var head *variantList
	tail := head
	for v := range chVariant {
		if head == nil {
			head = &variantList{0, v, nil}
			tail = head
		} else {
			head.next = &variantList{head.index + 1, v, nil}
			head = head.next
		}
		if head.variant.Pos-tail.variant.Pos >= params.WindowSize {
			workQueue <- tail
			tail = tail.next
		}
	}
	for {
		if tail == nil {
			break
		}
		workQueue <- tail
		tail = tail.next
	}
	close(workQueue)
}

func startReorderQueue(numWorkers uint16, chResult chan *Result) (reorderQueue chan *Result) {
	reorderQueue = make(chan *Result, numWorkers*10)
	go reorderResults(reorderQueue, numWorkers, chResult)
	return reorderQueue
}

func startWorkers(numWorkers uint16, chResult chan *Result) (workQueue chan *variantList) {
	//TODO expose queue size somewhere (instead of hardcoding "2")
	workQueue = make(chan *variantList, numWorkers*2) //make queue a little larger so that workers always have work.
	var i uint16
	for i = 0; i < numWorkers; i++ {
		go runWorker(workQueue, chResult)
	}
	return workQueue
}

func runWorker(in chan *variantList, out chan *Result) {
	for vl := range in {
		out <- &Result{vl.variant, nil}
	}
	out <- nil
}

func reorderResults(in chan *Result, numWorkers uint16, out chan *Result) {
	var finishedWorkers uint16
	for r := range in {
		if r == nil {
			finishedWorkers++
			if finishedWorkers >= numWorkers {
				close(out)
				break
			}
		} else {
			out <- r
		}
	}
}

func CreateEngine(params EngineParameters) (e Engine, err error) {
	initializeBitCountMap()
	return Engine(&engine{params}), nil
}

var bitCountMap map[uint16]uint16

func initializeBitCountMap() {
	bitCountMap = make(map[uint16]uint16)
	for pow := uint8(0); pow < 16; pow++ {
		for i := uint16(0); i < 1<<pow; i++ {
			bitCountMap[i+1<<pow] = bitCountMap[i] + 1
		}
	}
}

func bitCountSingle(a uint32) uint16 {
	return bitCountMap[uint16(a)] + bitCountMap[uint16(a>>16)]
}

func bitCount(genotypes []uint32) (bitCount uint16) {
	for _, i := range genotypes {
		bitCount += bitCountSingle(i)
	}
	return bitCount
}

func union(a []uint32, b []uint32) (c []uint32) {
	c = make([]uint32, len(a), len(a))
	for i := range a {
		c[i] = a[i] | b[i]
	}
	return c
}

func ComputeR2(a []uint32, b []uint32, bitLength uint16) (r2 float64) {
	pAB := bitCount(union(a, b))
	pa, pb := bitCount(a), bitCount(b)
	return calculateR2(pa, pb, pAB, bitLength)
}

func calculateR2(pa uint16, pb uint16, pAB uint16, size uint16) float64 {
	return 0
}
