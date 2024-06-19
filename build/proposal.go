package build

import "github.com/noppikinatta/ebitenginejam03/geom"

type Proposal struct {
	Equip           *Equip
	Cost            int
	Hit             geom.Circle
	Velocity        geom.PointF
	Rotate          float64
	RotateVelocity  float64
	CustomImageName string
}

func (p *Proposal) EquipName() string {
	return p.Equip.Name
}

func (p *Proposal) Clone() *Proposal {
	copyP := *p
	return &copyP
}

func (p *Proposal) Update() {
	p.Hit.Center.X += p.Velocity.X
	p.Hit.Center.Y += p.Velocity.Y
	p.Rotate += p.RotateVelocity
}

func (p *Proposal) BoundTop(y float64) {
	if p.Hit.Top() < y {
		p.Hit.Center.Y += (y - p.Hit.Top())
	}
	if p.Velocity.Y < 0 {
		p.Velocity.Y *= -1
	}
}

func (b *Proposal) BoundBottom(y float64) {
	if b.Hit.Bottom() > y {
		b.Hit.Center.Y -= (b.Hit.Bottom() - y)
	}
	if b.Velocity.Y > 0 {
		b.Velocity.Y *= -1
	}
}

func (b *Proposal) BoundLeft(x float64) {
	if b.Hit.Left() < x {
		b.Hit.Center.X += (x - b.Hit.Left())
	}
	if b.Velocity.X < 0 {
		b.Velocity.X *= -1
	}
}

func (b *Proposal) BoundRight(x float64) {
	if b.Hit.Right() > x {
		b.Hit.Center.X -= (b.Hit.Right() - x)
	}
	if b.Velocity.X > 0 {
		b.Velocity.X *= -1
	}
}

func (b *Proposal) MultiplyVelocity(v float64) {
	a := b.Velocity.Abs()
	r := b.Velocity.Direction360()

	a *= v

	b.Velocity = geom.PointFFromPolar(a, r)
}

func (b *Proposal) AddRotateVelocity(v float64) {
	b.RotateVelocity += v
}

type ProposalLaunchDelay struct {
	proposal    *Proposal
	waitingTime int
}

func NewProposalLaunchDelay(proposal *Proposal, waitingTime int) *ProposalLaunchDelay {
	return &ProposalLaunchDelay{
		proposal:    proposal,
		waitingTime: waitingTime,
	}
}

func (d *ProposalLaunchDelay) EquipName() string {
	return d.proposal.Equip.Name
}

func (d *ProposalLaunchDelay) Update() (*Proposal, bool) {
	if d.waitingTime > 0 {
		d.waitingTime -= 1
		return nil, false
	}

	return d.proposal, true
}
