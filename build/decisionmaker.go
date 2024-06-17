package build

type DecisionMaker struct {
	HitLine *HorizontalHitLine
}

func (d *DecisionMaker) Update(xCenter float64) {
	halfLen := d.HitLine.Length / 2
	newX := xCenter - halfLen
	d.HitLine.Position.X = newX
}
