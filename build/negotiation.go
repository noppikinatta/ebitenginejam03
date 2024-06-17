package build

type Negotiation struct {
	Size           PointF
	DecisionMaker  *DecisionMaker
	Vendors        []*Vendor
	Managers       []*Manager
	Money          int
	ProposalDelays []*ProposalLaunchDelay
	Proposals      []*Proposal
	ApprovedEquips []*Equip
}

func (n *Negotiation) Update(decisionMakerX float64) {
	n.updateDecisionMaker(decisionMakerX)
	n.updateVendors()
	n.updateProposalDelays()
	n.updateProposals()
}

func (n *Negotiation) updateDecisionMaker(decisionMakerX float64) {
	n.DecisionMaker.Update(decisionMakerX)
}

func (n *Negotiation) updateVendors() {
	// vendor position calced from negotiation area <- use to decide proposal start position
	for _, v := range n.Vendors {
		d, ok := v.Update()
		if !ok {
			continue
		}

		n.ProposalDelays = append(n.ProposalDelays, d)
	}
}

func (n *Negotiation) updateProposalDelays() {
	n.ProposalDelays = convertPartial(
		n.ProposalDelays,
		func(d *ProposalLaunchDelay) (*Proposal, bool) {
			return d.Update()
		},
		func(p *Proposal) {
			n.Proposals = append(n.Proposals, p)
		},
	)
}

func (n *Negotiation) updateProposals() {
	n.Proposals = convertPartial(
		n.Proposals,
		n.updateProposal,
		func(e *Equip) {
			n.ApprovedEquips = append(n.ApprovedEquips, e)
		},
	)
}

func (n *Negotiation) updateProposal(proposal *Proposal) (*Equip, bool) {
	// proposal old pos and new pos
	// update proposal pos
	// bound
	// determine cross manager // manager position calced from negotiation
	// process if crossed
	// if y < 0, go to approved
	//  - remove from proposals
	//  - add to approves

}

func (n *Negotiation) End() bool {
	return n.Money == 0
}

func convertPartial[T1, T2 any](t1Slice []T1, updateT1Fn func(T1) (T2, bool), processT2Fn func(T2)) []T1 {
	length := len(t1Slice)
	if length == 0 {
		return t1Slice
	}

	i := 0

	for i < length {
		item := t1Slice[i]
		result, ok := updateT1Fn(item)
		if ok {
			processT2Fn(result)
			if i < (length - 1) {
				t1Slice[i] = t1Slice[length-1]
			}
			length--
		} else {
			i++
		}
	}

	return t1Slice[:length]
}
