package ld

import (
	"fmt"
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
	reorderQueue := startReorderQueue(&params, chResult)
	workQueue := startWorkers(&params, reorderQueue)
	var head *variantList
	tail := head
	var last *Variant
	for v := range chVariant {
		if last != nil && last.Pos > v.Pos {
			panic(fmt.Sprintf("variants out of order: %v(pos %v) followed by %v(pos %v)", last.Rsid, last.Pos, v.Rsid, v.Pos))
		}
		last = v
		//fmt.Printf("last: %v v: %v\n", last, v)
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
	for tail != nil { //serve up the tail
		workQueue <- tail
		tail = tail.next
	}
	close(workQueue)
}

func startReorderQueue(params *EngineParameters, chResult chan *Result) (chReorder chan *Result) {
	chReorder = make(chan *Result, params.NumWorkers*10) //make reorder queue large enough so workers don't get blocked on output. TODO: expose constant 10
	go reorderResults(chReorder, params.NumWorkers, chResult)
	//go dontReorderResults(chReorder, params.NumWorkers, chResult)
	return chReorder
}

func startWorkers(params *EngineParameters, chResult chan *Result) (chWork chan *variantList) {
	chWork = make(chan *variantList, params.NumWorkers*10) //make work channel a multichannel so that workers always have work. TODO: expose constant 10
	for i := uint16(0); i < params.NumWorkers; i++ {
		go runWorker(chWork, chResult, params)
	}
	return chWork
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

func dontReorderResults(in chan *Result, numWorkers uint16, out chan *Result) {
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

func reorderResults(in chan *Result, numWorkers uint16, out chan *Result) {
	var finishedWorkers uint16
	type resultList struct {
		prev   *resultList
		result *Result
		next   *resultList
	}
	type resultQueue struct {
		tail  *resultList
		head  *resultList
		depth uint32
	}
	queue := resultQueue{nil, nil, 0}
	limit := uint32(numWorkers * 20)

	for r := range in {
		if r == nil {
			finishedWorkers++
			if finishedWorkers >= numWorkers {
				for queue.tail != nil {
					//fmt.Printf("tail%v depth: %v\n", queue.tail, queue.depth)
					out <- queue.tail.result
					queue.tail = queue.tail.next
					queue.depth--
				}
				close(out)
				break
			}
		} else {
			if queue.head == nil {
				queue.head = &resultList{nil, r, nil}
				queue.tail = queue.head
				queue.depth++
			} else {
				current := queue.head
				//for current != queue.tail && current.result.Variant.Pos > r.Variant.Pos {
				for current != nil && current.result.Variant.Pos > r.Variant.Pos {
					current = current.prev
				}
				if current != nil {
					//fmt.Printf("head: %v current: %v\n", queue.head, current)
					//fmt.Printf("head: %v current: %v\n", queue.head.result.Variant.Pos, current.result.Variant.Pos)
					next := current.next
					current.next = &resultList{current, r, next}
					if next != nil {
						next.prev = current.next
					}
				} else {
					queue.tail = &resultList{nil, r, queue.tail}
					queue.tail.next.prev = queue.tail
				}

				if queue.head.next != nil {
					queue.head = queue.head.next
				}
				queue.depth++
				//fmt.Printf("queue(%v): ", queue.depth)
				//for cur := queue.tail; cur != nil; cur = cur.next {
				//fmt.Printf("%v ", cur.result.Variant.Pos)
				//}
				//fmt.Printf("\n")
			}
			if queue.depth >= limit {
				out <- queue.tail.result
				queue.tail = queue.tail.next
				queue.tail.prev = nil
				queue.depth--
			}
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
	if pa == 0 || pb == 0 || pa == bitLength || pb == bitLength { //if any of this is true, r2 formula is undefined. so just return 0. TODO: document this logic.
		return 0
	}
	//fmt.Printf("pa: %v, pb: %v, pAB: %v, bitLength: %v\n", pa, pb, pAB, bitLength)
	return round(calculateR2(int64(pa), int64(pb), int64(pAB), int64(bitLength)), 6)
}

func calculateR2(pa int64, pb int64, aorb int64, size int64) float64 {
	return math.Pow(float64((size-aorb)*size-(size-pa)*(size-pb)), 2) / float64((size-pa)*pa*(size-pb)*pb)
}
