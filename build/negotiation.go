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

func (n *Negotiation) Update() {
	n.updateVendors()
	n.updateProposalDelays()
}

func (n *Negotiation) updateVendors() {
	for _, v := range n.Vendors {
		d, ok := v.Update()
		if !ok {
			continue
		}

		n.ProposalDelays = append(n.ProposalDelays, d)
	}
}

func (n *Negotiation) updateProposalDelays() {
	n.ProposalDelays = loopUpdaterSlice(
		n.ProposalDelays,
		func(p *Proposal) {
			n.Proposals = append(n.Proposals, p)
		},
	)
}

func (n *Negotiation) End() bool {
	return n.Money == 0
}

type updater[T any] interface {
	Update() (T, bool)
}

func loopUpdaterSlice[T any, S updater[T]](s []S, fn func(T)) []S {
	length := len(s)
	if length == 0 {
		return s
	}

	i := 0

	for i < length {
		item := s[i]
		result, ok := item.Update()
		if ok {
			fn(result)
			if i < (length - 1) {
				s[i] = s[length-1]
			}
			length--
		} else {
			i++
		}
	}

	return s[:length]
}
