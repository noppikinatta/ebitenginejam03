package build

import (
	"math/rand/v2"

	"github.com/noppikinatta/ebitenginejam03/geom"
)

type Vendor struct {
	Name      string
	proposals []*Proposal
	rnd       *rand.Rand
}

func NewVendor(name string, proposals []*Proposal, rnd *rand.Rand) *Vendor {
	return &Vendor{
		Name:      name,
		proposals: proposals,
		rnd:       rnd,
	}
}

func (v *Vendor) Propose(pos geom.PointF) *Proposal {
	p := v.randProposal()
	p = p.Clone()
	p.Hit.Center = pos

	va := p.Velocity.Abs()
	vr := v.randDirection()

	p.Velocity = geom.PointFFromPolar(va, vr)

	return p
}

func (v *Vendor) randProposal() *Proposal {
	return v.proposals[v.rnd.IntN(len(v.proposals))]
}

func (v *Vendor) randDirection() float64 {
	// -30 ~ -150 degrees
	return -1 * (v.rnd.Float64()*120 + 30)
}

type VendorSelector struct {
	Vendors       []*Vendor
	selectedCount []int
	interval      int
	framesToWait  int
	rnd           *rand.Rand
}

func NewVendorSelector(vendors []*Vendor, interval int, rnd *rand.Rand) *VendorSelector {
	return &VendorSelector{
		Vendors:       vendors,
		selectedCount: make([]int, len(vendors)),
		interval:      interval,
		rnd:           rnd,
	}
}

func (s *VendorSelector) Reset() {
	for i := range s.selectedCount {
		s.selectedCount[i] = 0
	}
}

func (s *VendorSelector) Length() int {
	return len(s.Vendors)
}

func (s *VendorSelector) IndexOf(vendor *Vendor) int {
	for i, v := range s.Vendors {
		if v == vendor {
			return i
		}
	}

	return -1
}

func (s *VendorSelector) Update() (*Vendor, bool) {
	if s.framesToWait > 0 {
		s.framesToWait--
		return nil, false
	}
	s.framesToWait = int(float64(s.interval) * (s.rnd.Float64()*0.4 + 0.8))

	rndMax := s.rndMax()
	rndVal := s.rnd.IntN(rndMax)
	max := 0
	for i, v := range s.Vendors {
		max += (1 + s.selectedCount[i])
		if rndVal < max {
			s.selectedCount[i]++
			return v, true
		}
	}

	return nil, false
}

func (s *VendorSelector) rndMax() int {
	base := len(s.Vendors)
	for _, v := range s.selectedCount {
		base += v
	}
	return base
}
