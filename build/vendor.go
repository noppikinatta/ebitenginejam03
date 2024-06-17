package build

type Vendor struct {
	Proposals   []*Proposal
	WaitFrame   int
	CurrentWait int
}
