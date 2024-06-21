package nego

import (
	"fmt"
	"math"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/ebitenginejam03/build"
	"github.com/noppikinatta/ebitenginejam03/drawing"
	"github.com/noppikinatta/ebitenginejam03/geom"
)

type NegotiationGameScene struct {
	Negotiation *build.Negotiation
	StagePos    geom.PointF
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

	return &NegotiationGameScene{Negotiation: &n, StagePos: geom.PointF{X: 0, Y: 40}}
}

func (s *NegotiationGameScene) Update() error {
	x, _ := ebiten.CursorPosition()
	x += int(s.StagePos.X)
	s.Negotiation.Update(float64(x))
	return nil
}

func (s *NegotiationGameScene) Draw(screen *ebiten.Image) {
	s.drawApprovedEquips(screen)
	s.drawMoney(screen)
	s.drawVendors(screen)
	s.drawManagers(screen)
	s.drawDecisionMaker(screen)
	s.drawProposals(screen)
	s.drawProposalDelay(screen)
}

func (s *NegotiationGameScene) drawApprovedEquips(screen *ebiten.Image) {
	if len(s.Negotiation.ApprovedEquips) == 0 {
		return
	}

	topLeft := geom.PointF{}
	topLeft.X = s.Negotiation.Size.X
	topLeft.Y = s.StagePos.Y

	rightBottom := geom.PointF{}
	rightBottom.X = float64(screen.Bounds().Max.X)
	rightBottom.Y = float64(screen.Bounds().Max.Y)

	areaHeight := rightBottom.Y - topLeft.Y
	itemHeight := math.Min(64, areaHeight/float64(len(s.Negotiation.ApprovedEquips)))

	for i, e := range s.Negotiation.ApprovedEquips {
		y := itemHeight * float64(i)

		opt := ebiten.DrawImageOptions{}
		opt.GeoM.Translate(topLeft.X, y)

		drawing.DrawText(screen, e.Name, 12, &opt)
	}
}

func (s *NegotiationGameScene) drawMoney(screen *ebiten.Image) {
	drawing.DrawText(screen, fmt.Sprint(s.Negotiation.Money), 18, &ebiten.DrawImageOptions{})
}

func (s *NegotiationGameScene) drawVendors(screen *ebiten.Image) {
	size := geom.PointF{X: 128, Y: 128}
	idxs := []uint16{0, 1, 2, 0, 2, 3}

	for i, v := range s.Negotiation.VendorSelector.Vendors {
		bottomCenter := s.Negotiation.ProposalStartPosition(i)
		bottomCenter = bottomCenter.Add(s.StagePos)

		vertices := make([]ebiten.Vertex, 4)
		vertices[0] = ebiten.Vertex{
			DstX: float32(bottomCenter.X - 0.5*size.X),
			DstY: float32(bottomCenter.Y - size.Y),
		}
		vertices[1] = ebiten.Vertex{
			DstX: float32(bottomCenter.X - 0.5*size.X),
			DstY: float32(bottomCenter.Y),
		}
		vertices[2] = ebiten.Vertex{
			DstX: float32(bottomCenter.X + 0.5*size.X),
			DstY: float32(bottomCenter.Y),
		}
		vertices[3] = ebiten.Vertex{
			DstX: float32(bottomCenter.X + 0.5*size.X),
			DstY: float32(bottomCenter.Y - size.Y),
		}

		topt := ebiten.DrawTrianglesOptions{}
		screen.DrawTriangles(vertices, idxs, drawing.WhitePixel, &topt)

		iopt := ebiten.DrawImageOptions{}
		iopt.GeoM.Translate(float64(vertices[0].DstX), float64(vertices[0].DstY))
		drawing.DrawText(screen, v.Name, 12, &iopt)
	}
}

func (s *NegotiationGameScene) drawManagers(screen *ebiten.Image) {
	size := geom.PointF{X: 128, Y: 128}
	idxs := []uint16{0, 1, 2, 0, 2, 3}

	for i, m := range s.Negotiation.Managers {
		line := s.Negotiation.ManagerHLine(i)
		center := line.Center()
		center = center.Add(s.StagePos)

		vertices := make([]ebiten.Vertex, 4)
		vertices[0] = ebiten.Vertex{
			DstX: float32(center.X - 0.5*size.X),
			DstY: float32(center.Y - 0.5*size.Y),
		}
		vertices[1] = ebiten.Vertex{
			DstX: float32(center.X - 0.5*size.X),
			DstY: float32(center.Y + 0.5*size.Y),
		}
		vertices[2] = ebiten.Vertex{
			DstX: float32(center.X + 0.5*size.X),
			DstY: float32(center.Y + 0.5*size.Y),
		}
		vertices[3] = ebiten.Vertex{
			DstX: float32(center.X + 0.5*size.X),
			DstY: float32(center.Y - 0.5*size.Y),
		}

		topt := ebiten.DrawTrianglesOptions{}
		screen.DrawTriangles(vertices, idxs, drawing.WhitePixel, &topt)

		iopt := ebiten.DrawImageOptions{}
		iopt.GeoM.Translate(float64(vertices[0].DstX), float64(vertices[0].DstY))
		drawing.DrawText(screen, m.Name, 12, &iopt)
	}
}

func (s *NegotiationGameScene) drawDecisionMaker(screen *ebiten.Image) {
	left := s.Negotiation.DecisionMaker.Left
	top, _ := s.Negotiation.DecisionMaker.LinearFn.Y(left)
	width := s.Negotiation.DecisionMaker.Length

	topLeft := geom.PointF{X: left, Y: top - 2}
	bottomRight := geom.PointF{X: left + width, Y: top + 2}
	topLeft = topLeft.Add(s.StagePos)
	bottomRight = bottomRight.Add(s.StagePos)

	idxs := []uint16{0, 1, 2, 0, 2, 3}
	vertices := make([]ebiten.Vertex, 4)
	vertices[0] = ebiten.Vertex{
		DstX: float32(topLeft.X),
		DstY: float32(topLeft.Y),
	}
	vertices[1] = ebiten.Vertex{
		DstX: float32(topLeft.X),
		DstY: float32(bottomRight.Y),
	}
	vertices[2] = ebiten.Vertex{
		DstX: float32(bottomRight.X),
		DstY: float32(bottomRight.Y),
	}
	vertices[3] = ebiten.Vertex{
		DstX: float32(bottomRight.X),
		DstY: float32(topLeft.Y),
	}

	topt := ebiten.DrawTrianglesOptions{}
	screen.DrawTriangles(vertices, idxs, drawing.WhitePixel, &topt)
}

func (s *NegotiationGameScene) drawProposals(screen *ebiten.Image) {
	idxs := []uint16{0, 1, 2, 0, 2, 3}

	for _, p := range s.Negotiation.Proposals {
		hit := p.Hit

		gm := ebiten.GeoM{}
		gm.Rotate(p.Rotate)
		gm.Translate(hit.Center.X, hit.Center.Y)
		gm.Translate(s.StagePos.X, s.StagePos.Y)
		var x, y float64

		vertices := make([]ebiten.Vertex, 4)
		x, y = gm.Apply(-hit.Radius, -hit.Radius)
		vertices[0] = ebiten.Vertex{
			DstX: float32(x),
			DstY: float32(y),
		}
		x, y = gm.Apply(-hit.Radius, hit.Radius)
		vertices[1] = ebiten.Vertex{
			DstX: float32(x),
			DstY: float32(y),
		}
		x, y = gm.Apply(hit.Radius, hit.Radius)
		vertices[2] = ebiten.Vertex{
			DstX: float32(x),
			DstY: float32(y),
		}
		x, y = gm.Apply(hit.Radius, -hit.Radius)
		vertices[3] = ebiten.Vertex{
			DstX: float32(x),
			DstY: float32(y),
		}

		topt := ebiten.DrawTrianglesOptions{}
		screen.DrawTriangles(vertices, idxs, drawing.WhitePixel, &topt)

		iopt := ebiten.DrawImageOptions{}
		iopt.GeoM.Translate(float64(vertices[0].DstX), float64(vertices[0].DstY))
		drawing.DrawText(screen, p.ImageName(), 12, &iopt)
	}
}

func (s *NegotiationGameScene) drawProposalDelay(screen *ebiten.Image) {
	if s.Negotiation.ProposalDelay != nil {
		return
	}

	startPos := s.Negotiation.ProposalDelay.Hit.Center.Add(s.StagePos)
	triangleVerts := make([]ebiten.Vertex, 3)
	triangleVerts[0] = ebiten.Vertex{
		DstX: float32(startPos.X),
		DstY: float32(startPos.Y),
	}
	triangleVerts[0] = ebiten.Vertex{
		DstX: float32(startPos.X + 16),
		DstY: float32(startPos.Y - 32),
	}
	triangleVerts[0] = ebiten.Vertex{
		DstX: float32(startPos.X - 16),
		DstY: float32(startPos.Y - 32),
	}

	topt := ebiten.DrawTrianglesOptions{}
	screen.DrawTriangles(triangleVerts, []uint16{0, 1, 2}, drawing.WhitePixel, &topt)

	baloonVerts := make([]ebiten.Vertex, 4)
	baloonTopLeft := geom.PointF{
		X: s.StagePos.X + 8,
		Y: startPos.Y - 32 - 64,
	}
	baloonBottomRight := geom.PointF{
		X: s.StagePos.X + s.Negotiation.Size.X - 8,
		Y: startPos.Y - 32 - 64,
	}
	baloonVerts[0] = ebiten.Vertex{
		DstX: float32(baloonTopLeft.X),
		DstY: float32(baloonTopLeft.Y),
	}
	baloonVerts[1] = ebiten.Vertex{
		DstX: float32(baloonTopLeft.X),
		DstY: float32(baloonBottomRight.Y),
	}
	baloonVerts[2] = ebiten.Vertex{
		DstX: float32(baloonBottomRight.X),
		DstY: float32(baloonBottomRight.Y),
	}
	baloonVerts[3] = ebiten.Vertex{
		DstX: float32(baloonBottomRight.X),
		DstY: float32(baloonTopLeft.Y),
	}
	idxs := []uint16{0, 1, 2, 0, 2, 3}
	screen.DrawTriangles(baloonVerts, idxs, drawing.WhitePixel, &topt)

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
