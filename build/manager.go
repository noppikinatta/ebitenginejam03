package build

type Manager struct {
	Name       string
	Processors []ProposalProcessor
}

func NewManager(name string, processors ...ProposalProcessor) *Manager {
	return &Manager{
		Name:       name,
		Processors: processors,
	}
}

func (m *Manager) Process(proposal *Proposal) {
	for _, p := range m.Processors {
		p.Process(proposal)
	}
}

type ProposalProcessor interface {
	Process(proposal *Proposal)
}

type ProposalProcessorAccelerate struct {
	Value float64
}

func (p *ProposalProcessorAccelerate) Process(proposal *Proposal) {
	proposal.MultiplyVelocity(p.Value)
}

type ProposalProcessorRotate struct {
	Value float64
}

func (p *ProposalProcessorRotate) Process(proposal *Proposal) {
	proposal.AddRotateVelocity(p.Value)
}

type ProposalProcessorStopRotate struct {
}

func (p *ProposalProcessorStopRotate) Process(proposal *Proposal) {
	proposal.RotateVelocity = 0
	proposal.Rotate = 0
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
