package nego

type Equip struct {
	Name           string
	ImprovedCount  int
	ImprovedByNext bool
	ImprovedByPrev bool
}

func (e Equip) CalcedImprovedCount() int {
	c := e.ImprovedCount
	if e.ImprovedByNext {
		c++
	}
	if e.ImprovedByPrev {
		c++
	}
	return c
}
