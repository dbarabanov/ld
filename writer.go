package ld

type Writer interface {
	Write(Variant, []Score) error
}

type writer struct{}

func (writer) Write(Variant, []Score) error {
	return nil
}

func NewFileWriter(filepath string) (Writer, error) {
	return Writer(new(writer)), nil
}
