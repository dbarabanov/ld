package ld

import (
	"fmt"
	"os"
)

type fileWriter struct {
	file *os.File
}

func (w fileWriter) WriteResults(ch chan *Result) {
	for r := range ch {
		pos1 := r.Variant.Pos
		for _, s := range r.Scores {
			w.file.WriteString(fmt.Sprintf("%v\t%v\t%v\n", pos1, s.Pos, s.R2))
		}
	}
}

func NewFileWriter(filepath string) (ResultWriter, error) {
	fo, err := os.Create(filepath)
	if err != nil {
		return nil, err
	}
	fo.WriteString("POS1\tPOS2\tR^2\n")
	return ResultWriter(&fileWriter{fo}), nil
}
