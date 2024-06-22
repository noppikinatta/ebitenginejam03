package build

import (
	"math"
	"math/rand/v2"

	"github.com/noppikinatta/ebitenginejam03/geom"
)

type Vendor struct {
	Name     string
	selector *randomSelector[*Proposal]
	rnd      *rand.Rand
}

func NewVendor(name string, proposals []*Proposal, rnd *rand.Rand) *Vendor {
	return &Vendor{
		Name:     name,
		selector: &randomSelector[*Proposal]{Items: proposals, rnd: rnd},
		rnd:      rnd,
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
	return v.selector.Select()
}

func (v *Vendor) randDirection() float64 {
	// -30 ~ -150 degrees
	return -1 * (v.rnd.Float64()*120 + 30) * math.Pi / 180
}

type VendorSelector struct {
	selector     *randomSelector[*Vendor]
	interval     int
	framesToWait int
	rnd          *rand.Rand
}

func NewVendorSelector(vendors []*Vendor, interval int, rnd *rand.Rand) *VendorSelector {
	return &VendorSelector{
		selector: &randomSelector[*Vendor]{Items: vendors, rnd: rnd},
		interval: interval,
		rnd:      rnd,
	}
}

func (s *VendorSelector) Reset() {
	s.framesToWait = 0
	s.selector.Reset()
}

func (s *VendorSelector) Vendors() []*Vendor {
	return s.selector.Items
}

func (s *VendorSelector) IndexOf(vendor *Vendor) int {
	for i, v := range s.Vendors() {
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

	return s.selector.Select(), true
}
