package ld

type Variant struct {
	pos   uint32
	rsid  uint64
	minor []uint32
}

func Equal(a, b *Variant) bool {
	if a.pos != b.pos || a.rsid != b.rsid || len(a.minor) != len(b.minor) {
		return false
	}
	for i := range a.minor {
		if a.minor[i] != b.minor[i] {
			return false
		}
	}
	return true
}
