package nego

import (
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/ebitenginejam03/build"
	"github.com/noppikinatta/ebitenginejam03/geom"
)

type NegotiationGameScene struct {
	Negotiation *build.Negotiation
}

func NewNegotiationGameScene() *NegotiationGameScene {
	n := build.Negotiation{
		Size:          geom.PointF{X: 600, Y: 600},
		DecisionMaker: build.NewDecisionMaker(0, 100),
		VendorSelector: build.NewVendorSelector(
			createVendors(),
			180,
			rand.New(rand.NewPCG(123, 456))),
		Managers: createManagers(),
		Money:    10000,
	}

	return &NegotiationGameScene{Negotiation: &n}
}

func (s *NegotiationGameScene) Update() error {
	x, _ := ebiten.CursorPosition()
	s.Negotiation.Update(float64(x))
	return nil
}

func (s *NegotiationGameScene) Draw(screen *ebiten.Image) {

}

func (s *NegotiationGameScene) End() bool {
	return s.Negotiation.End()
}

func (s *NegotiationGameScene) Reset() {
	s.Negotiation.Reset(10000)
}

func createVendors() []*build.Vendor {
	pm := createProposals()

	vv := make([]*build.Vendor, 0)
	vv = append(vv, build.NewVendor("samurai-avionics", selectProposals(pm, "laser-cannon", "space-missile", "harakiri-system"), rand.New(rand.NewPCG(314, 1592))))
	vv = append(vv, build.NewVendor("salamis-industry", selectProposals(pm, "barrier", "armor-plate", "weak-point"), rand.New(rand.NewPCG(6535, 8979))))
	vv = append(vv, build.NewVendor("cultual-victory-co", selectProposals(pm, "opera-house"), rand.New(rand.NewPCG(3238, 4626))))

	return vv
}

func selectProposals(m map[string]*build.Proposal, names ...string) []*build.Proposal {
	pp := make([]*build.Proposal, 0, len(names))
	for _, n := range names {
		pp = append(pp, m[n])
	}

	return pp
}

func createProposals() map[string]*build.Proposal {
	hit := geom.Circle{Radius: 32}

	m := map[string]*build.Proposal{
		"laser-cannon":    {Equip: &build.Equip{Name: "laser-cannon"}, Cost: 1000, Hit: hit},
		"space-missile":   {Equip: &build.Equip{Name: "space-missile"}, Cost: 800, Hit: hit},
		"harakiri-system": {Equip: &build.Equip{Name: "harakiri-system"}, Cost: 1200, Hit: hit},
		"barrier":         {Equip: &build.Equip{Name: "barrier"}, Cost: 1000, Hit: hit},
		"armor-plate":     {Equip: &build.Equip{Name: "armor-plate"}, Cost: 800, Hit: hit},
		"weak-point":      {Equip: &build.Equip{Name: "weak-point"}, Cost: 500, Hit: hit},
		"opera-house":     {Equip: &build.Equip{Name: "opera-house"}, Cost: 5000, Hit: hit},
	}

	return m
}

func createManagers() []*build.Manager {
	mm := make([]*build.Manager, 0)
	mm = append(mm, build.NewManager(
		"mach-sonic",
		&build.ProposalProcessorAccelerate{Value: 1.5},
		&build.ProposalProcessorStopRotate{},
		&build.ProposalProcessorCustomImageName{ImageName: ""}))
	mm = append(mm, build.NewManager(
		"birdie-pat",
		&build.ProposalProcessorReduceCost{Multiplier: 0.8},
		&build.ProposalProcessorCustomImageName{ImageName: "golf"}))
	mm = append(mm, build.NewManager(
		"long-winded",
		&build.ProposalProcessorAccelerate{Value: 0.8},
		&build.ProposalProcessorRotate{Value: 1},
		&build.ProposalProcessorImprove{}))

	return mm
}
