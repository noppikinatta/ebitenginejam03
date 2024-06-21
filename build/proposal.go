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

func (p *Proposal) ImageName() string {
	if len(p.CustomImageName) > 0 {
		return p.CustomImageName
	}
	return p.Equip.Name
}

func (p *Proposal) Clone() *Proposal {
	copyP := &Proposal{}
	copyE := &Equip{}
	*copyP = *p
	*copyE = *p.Equip
	copyP.Equip = copyE

	return copyP
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

func (p *Proposal) BoundBottom(y float64) {
	if p.Hit.Bottom() > y {
		p.Hit.Center.Y -= (p.Hit.Bottom() - y)
	}
	if p.Velocity.Y > 0 {
		p.Velocity.Y *= -1
	}
}

func (p *Proposal) BoundLeft(x float64) {
	if p.Hit.Left() < x {
		p.Hit.Center.X += (x - p.Hit.Left())
	}
	if p.Velocity.X < 0 {
		p.Velocity.X *= -1
	}
}

func (p *Proposal) BoundRight(x float64) {
	if p.Hit.Right() > x {
		p.Hit.Center.X -= (p.Hit.Right() - x)
	}
	if p.Velocity.X > 0 {
		p.Velocity.X *= -1
	}
}

func (p *Proposal) MultiplyVelocity(v float64) {
	a := p.Velocity.Abs()
	r := p.Velocity.Direction360()

	a *= v

	p.Velocity = geom.PointFFromPolar(a, r)
}

func (p *Proposal) AddRotateVelocity(v float64) {
	p.RotateVelocity += v
}
