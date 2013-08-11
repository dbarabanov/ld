package ld

type Variant struct {
	pos   uint32
	rsid  uint64
	minor []uint32
}

type VariantReader interface {
	Read() chan Variant
}

type VariantProperties struct {
	chromosome     string
	populationSize uint16
}

type variantReader struct {
	prop VariantProperties
}

func (r variantReader) Read() chan Variant {
	return nil
}

func NewVcfReader(vcfFilePath string) (v VariantReader, err error) {
	return VariantReader(&variantReader{VariantProperties{"", 0}}), nil
}
