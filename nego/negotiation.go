package nego

import (
	"github.com/noppikinatta/ebitenginejam03/geom"
	"github.com/noppikinatta/ebitenginejam03/name"
)

type Negotiation struct {
	Size               geom.PointF
	DecisionMaker      *DecisionMaker
	VendorSelector     *VendorSelector
	LastSelectedVendor *Vendor
	Managers           []*Manager
	Money              int
	ProposalDelay      *Proposal
	Proposals          []*Proposal
	ApprovedEquips     []*Equip
}

func (n *Negotiation) UpdateDecisionMaker(decisionMakerX float64) {
	n.updateDecisionMaker(decisionMakerX)
}

func (n *Negotiation) UpdateOthers() {
	n.updateVendors()
	n.updateProposals()
}

func (n *Negotiation) updateDecisionMaker(decisionMakerX float64) {
	min := n.DecisionMaker.Width / 2
	max := n.Size.X - min
	if decisionMakerX < min {
		decisionMakerX = min
	}
	if decisionMakerX > max {
		decisionMakerX = max
	}
	n.DecisionMaker.Update(decisionMakerX)
}

func (n *Negotiation) updateVendors() {
	vendor, ok := n.VendorSelector.Update()
	if !ok {
		return
	}

	idx := n.VendorSelector.IndexOf(vendor)
	pos := n.ProposalStartPosition(idx)

	n.LastSelectedVendor = vendor
	p := vendor.Propose(pos)
	if n.ProposalDelay != nil {
		n.Proposals = append(n.Proposals, n.ProposalDelay)
	}
	n.ProposalDelay = p
}

func (n *Negotiation) ProposalStartPosition(idx int) geom.PointF {
	width := n.Size.X / float64(len(n.VendorSelector.Vendors()))
	x := (float64(idx) + 0.5) * width

	return geom.PointF{X: x, Y: n.Size.Y}
}

func (n *Negotiation) updateProposals() {
	n.Proposals = convertPartial(
		n.Proposals,
		n.updateProposal,
		func(p *Proposal) {
			n.Money -= p.Cost
			if n.Money < 0 {
				n.Money = 0
			}
			n.ApprovedEquips = append(n.ApprovedEquips, p.Equip)
			n.updateImproved()
		},
	)
}

func (n *Negotiation) updateImproved() {
	for i := range n.ApprovedEquips {
		n.ApprovedEquips[i].ImprovedByNext = false
		n.ApprovedEquips[i].ImprovedByPrev = false
	}

	eqplen := len(n.ApprovedEquips)

	for i, e := range n.ApprovedEquips {
		if e.Name != name.TextKeyEquip6Exhaust {
			continue
		}

		prev := ((i - 1) + eqplen) % eqplen
		if prev != i {
			n.ApprovedEquips[prev].ImprovedByNext = true
		}

		next := (i + 1) % eqplen
		if next != i && prev != next {
			n.ApprovedEquips[next].ImprovedByPrev = true
		}
	}
}

func (n *Negotiation) updateProposal(proposal *Proposal) (*Proposal, bool) {
	oldPos := proposal.Hit.Center
	proposal.Update()
	proposal.BoundBottom(n.Size.Y)
	proposal.BoundLeft(0)
	proposal.BoundRight(n.Size.X)

	if n.DecisionMaker.Hit(proposal.Hit) {
		proposal.BoundTop(0)
	}

	newPos := proposal.Hit.Center

	managerY := n.ManagerY()
	for i, m := range n.Managers {
		mL, mR := n.ManagerXLeftRight(i)
		if oldPos.Y >= managerY && newPos.Y <= managerY {
			pX := (oldPos.X + newPos.X) / 2
			if pX >= mL && pX <= mR {
				m.Process(proposal)
			}
		}
	}

	if proposal.Hit.Center.Y < 0 {
		return proposal, true
	}

	return nil, false
}

func (n *Negotiation) ManagerY() float64 {
	return n.Size.Y * 0.5
}

func (n *Negotiation) ManagerXLeftRight(idx int) (l, r float64) {
	width := n.Size.X / float64(len(n.Managers))
	left := width * float64(idx)

	return left, left + width
}

func (n *Negotiation) End() bool {
	return n.Money == 0
}

func (n *Negotiation) Reset(money int) {
	n.ApprovedEquips = nil
	n.ProposalDelay = nil
	n.Proposals = nil
	n.VendorSelector.Reset()
	n.Money = money
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
