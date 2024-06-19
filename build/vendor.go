package build

import (
	"math/rand/v2"

	"github.com/noppikinatta/ebitenginejam03/geom"
)

type Vendor struct {
	proposals []*Proposal
	rnd       *rand.Rand
}

func NewVendor(proposals []*Proposal, rnd *rand.Rand) *Vendor {
	return &Vendor{
		proposals: proposals,
		rnd:       rnd,
	}
}

func (v *Vendor) Propose(pos geom.PointF) *ProposalLaunchDelay {
	p := v.randProposal()
	p = p.Clone()
	p.Hit.Center = pos

	va := p.Velocity.Abs()
	vr := v.randDirection()

	p.Velocity = geom.PointFFromPolar(va, vr)

	return NewProposalLaunchDelay(p, 120)
}

func (v *Vendor) randProposal() *Proposal {
	return v.proposals[v.rnd.IntN(len(v.proposals))]
}

func (v *Vendor) randDirection() float64 {
	// -30 ~ -150 degrees
	return -1 * (v.rnd.Float64()*120 + 30)
}

type VendorSelector struct {
	vendors       []*Vendor
	selectedCount []int
	interval      int
	framesToWait  int
	rnd           *rand.Rand
}

func NewVendorSelector(vendors []*Vendor, interval int, rnd *rand.Rand) *VendorSelector {
	return &VendorSelector{
		vendors:       vendors,
		selectedCount: make([]int, len(vendors)),
		interval:      interval,
		framesToWait:  interval,
		rnd:           rnd,
	}
}

func (s *VendorSelector) Reset() {
	for i := range s.selectedCount {
		s.selectedCount[i] = 0
	}
	s.framesToWait = s.interval
}

func (s *VendorSelector) Length() int {
	return len(s.vendors)
}

func (s *VendorSelector) IndexOf(vendor *Vendor) int {
	for i, v := range s.vendors {
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
	s.framesToWait = int(float64(s.interval)*(s.rnd.Float64()*0.4) + 0.8)

	rndMax := s.rndMax()
	rndVal := s.rnd.IntN(rndMax)
	max := 0
	for i, v := range s.vendors {
		max += (1 + s.selectedCount[i])
		if rndVal < max {
			s.selectedCount[i]++
			return v, true
		}
	}

	return nil, false
}

func (s *VendorSelector) rndMax() int {
	base := len(s.vendors)
	for _, v := range s.selectedCount {
		base += v
	}
	return base
}
