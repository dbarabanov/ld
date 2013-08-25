package ld

import (
	//"fmt"
	"math"
)

type engine struct {
	params EngineParameters
}

type variantList struct {
	index   uint32
	variant *Variant
	next    *variantList
}

func (e engine) Run(chVariant chan *Variant) chan *Result {
	chResult := make(chan *Result, 10) //TODO expose 10, depth of result queue
	go runEngine(chVariant, chResult, e.params)
	return chResult
}

func runEngine(chVariant chan *Variant, chResult chan *Result, params EngineParameters) {
	//var workQueue chan *variantList
	//TODO expose queue size somewhere (instead of hardcoding "2")
	//workQueue := make(chan *variantList, params.NumWorkers*2) //make queue a little larger so that workers always have work.
	reorderQueue := startReorderQueue(&params, chResult)
	workQueue := startWorkers(&params, reorderQueue)
	var head *variantList
	tail := head
	for v := range chVariant {
		//fmt.Printf("v: %v\n", v)
		if head == nil {
			head = &variantList{0, v, nil}
			tail = head
		} else {
			head.next = &variantList{head.index + 1, v, nil}
			head = head.next
		}
		//fmt.Printf("head: %v, tail: %v\n", head, tail)
		if head.variant.Pos-tail.variant.Pos > params.WindowSize {
			workQueue <- tail
			tail = tail.next
		}
	}
	for tail != nil { //surve up the tail
		workQueue <- tail
		tail = tail.next
	}
	close(workQueue)
}

func startReorderQueue(params *EngineParameters, chResult chan *Result) (reorderQueue chan *Result) {
	reorderQueue = make(chan *Result, params.NumWorkers*10)
	go reorderResults(reorderQueue, params.NumWorkers, chResult)
	return reorderQueue
}

func startWorkers(params *EngineParameters, chResult chan *Result) (workQueue chan *variantList) {
	//TODO expose queue size somewhere (instead of hardcoding "5")
	workQueue = make(chan *variantList, params.NumWorkers*5) //make queue a little larger so that workers always have work.
	var i uint16
	for i = 0; i < params.NumWorkers; i++ {
		go runWorker(workQueue, chResult, params)
	}
	return workQueue
}

func runWorker(in chan *variantList, out chan *Result, params *EngineParameters) {
	for vlist := range in {
		v := vlist.variant
		next := vlist.next
		var scores []Score
		for next != nil && next.variant.Pos-v.Pos <= params.WindowSize {
			r2 := ComputeR2(v.Genotypes, next.variant.Genotypes, params.PopulationSize*2)
			if r2 > params.R2Threshold {
				scores = append(scores, Score{next.variant.Pos, next.variant.Rsid, r2})
			}
			next = next.next
		}
		out <- &Result{v, scores}
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
	if pa == 0 || pb == 0 || pa == bitLength || pb == bitLength {
		return -1
	}
	//fmt.Printf("pa: %v, pb: %v, pAB: %v, bitLength: %v\n", pa, pb, pAB, bitLength)
	return round(calculateR2(int64(pa), int64(pb), int64(pAB), int64(bitLength)), 6)
}

func calculateR2(pa int64, pb int64, aorb int64, size int64) float64 {
	return math.Pow(float64((size-aorb)*size-(size-pa)*(size-pb)), 2) / float64((size-pa)*pa*(size-pb)*pb)
}
