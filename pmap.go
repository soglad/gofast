package gofast

//PMap is the precence map of a TemplateUnit.
type PMap struct {
	bitMap []byte
	offset int
}

//HasNextPresenceBit checks if next bit is set to true in bitmap.
func (pmap *PMap) HasNextPresenceBit() bool {
	return false
}
