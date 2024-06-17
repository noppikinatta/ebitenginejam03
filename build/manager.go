package build

type Manager struct {
	Name      string
	Position  *HorizontalHitLine
	Processor ProposalProcessor
}

func (m *Manager) Process(proposal *Proposal) {
	m.Processor.Process(proposal)
}

type ProposalProcessor interface {
	Process(proposal *Proposal)
}

type ProposalProcessorAccelerate struct {
	Value float64
}

func (p *ProposalProcessorAccelerate) Process(proposal *Proposal) {
	proposal.HitBox.MultiplyVelocity(p.Value)
}

type ProposalProcessorRotate struct {
	Value float64
}

func (p *ProposalProcessorRotate) Process(proposal *Proposal) {
	proposal.HitBox.AddRotateVelocity(p.Value)
}

type ProposalProcessorCustomImageName struct {
	ImageName string
}

func (p *ProposalProcessorCustomImageName) Procell(proposal *Proposal) {
	proposal.CustomImageName = p.ImageName
}
