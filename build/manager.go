package build

type Manager struct {
	Name      string
	Processor ProposalProcessor
}

func (m *Manager) Process(proposal *Proposal) {
	m.Processor.Process(proposal)
}

type ProposalProcessor interface {
	Process(proposal *Proposal)
}

type ProposalProcessorCollection []ProposalProcessor

func (pp ProposalProcessorCollection) Process(proposal *Proposal) {
	for _, p := range pp {
		p.Process(proposal)
	}
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

func (p *ProposalProcessorCustomImageName) Process(proposal *Proposal) {
	proposal.CustomImageName = p.ImageName
}

type ProposalProcessorImprove struct {
}

func (p *ProposalProcessorImprove) Process(proposal *Proposal) {
	proposal.Equip.ImprovedCount++
}

type ProposalProcessorReduceCost struct {
	Multiplier float64
}

func (p *ProposalProcessorReduceCost) Process(proposal *Proposal) {
	proposal.Cost = int(float64(proposal.Cost) * p.Multiplier)
}
