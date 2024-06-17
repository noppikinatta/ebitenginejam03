package build

import "math/rand/v2"

type Vendor struct {
	proposals   []*Proposal
	interval    int
	waitingTime int
	rnd         *rand.Rand
}

func NewVendor(proposals []*Proposal, interval int, rnd *rand.Rand) *Vendor {
	return &Vendor{
		proposals:   proposals,
		interval:    interval,
		waitingTime: interval,
		rnd:         rnd,
	}
}

func (v *Vendor) Update() (*ProposalLaunchDelay, bool) {
	if v.waitingTime > 0 {
		v.waitingTime -= 1
		return nil, false
	}

	p := v.randProposal()
	p = p.Clone()

	va := p.HitBox.Velocity.Abs()
	vr := v.randDirection()

	p.HitBox.Velocity = PointFFromPolar(va, vr)

	v.waitingTime = v.interval

	return NewProposalLaunchDelay(p, 120), true
}

func (v *Vendor) randProposal() *Proposal {
	return v.proposals[v.rnd.IntN(len(v.proposals))]
}

func (v *Vendor) randDirection() float64 {
	// -30 ~ -150 degrees
	return -1 * (v.rnd.Float64()*120 + 30)
}
