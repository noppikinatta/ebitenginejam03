package build

type Proposal struct {
	Equip           *Equip
	Cost            int
	HitBox          *HitBox
	CustomImageName string
}

func (p *Proposal) EquipName() string {
	return p.Equip.Name
}

func (p *Proposal) Clone() *Proposal {
	copyP := *p
	return &copyP
}

func (p *Proposal) SetPosition(pt PointF) {
	p.HitBox.Position = pt
}

func (p *Proposal) SetVelocity(v PointF) {
	p.HitBox.Velocity = v
}

func (p *Proposal) Update() {
	p.HitBox.Update()
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
