package nego

import (
	"fmt"
	"image/color"
	"math"
	"math/rand/v2"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/ebitenginejam03/drawing"
	"github.com/noppikinatta/ebitenginejam03/geom"
	"github.com/noppikinatta/ebitenginejam03/name"
	"github.com/noppikinatta/ebitenginejam03/nego"
)

type negotiationGameScene struct {
	Negotiation *nego.Negotiation
	StagePos    geom.PointF
}

func newNegotiationGameScene() *negotiationGameScene {
	n := nego.Negotiation{
		Size:          geom.PointF{X: 600, Y: 600},
		DecisionMaker: nego.NewDecisionMaker(0, 100),
		VendorSelector: nego.NewVendorSelector(
			createVendors(),
			180,
			rand.New(rndSrc())),
		Managers: createManagers(),
		Money:    10000,
	}

	return &negotiationGameScene{Negotiation: &n, StagePos: geom.PointF{X: 0, Y: 40}}
}

func (s *negotiationGameScene) Update() error {
	x, _ := ebiten.CursorPosition()
	x += int(s.StagePos.X)
	s.Negotiation.Update(float64(x))
	return nil
}

func (s *negotiationGameScene) Draw(screen *ebiten.Image) {
	s.drawBackground(screen)
	s.drawApprovedEquips(screen)
	s.drawMoney(screen)
	s.drawVendors(screen)
	s.drawManagers(screen)
	s.drawDecisionMaker(screen)
	s.drawProposals(screen)
	s.drawProposalDelay(screen)
}

func (s *negotiationGameScene) drawRect(screen *ebiten.Image, topLeft, bottomRight geom.PointF, colorVert ebiten.Vertex) {
	idxs := []uint16{0, 1, 2, 0, 2, 3}

	v := colorVert
	vertices := make([]ebiten.Vertex, 4)
	v.DstX = float32(topLeft.X)
	v.DstY = float32(topLeft.Y)
	vertices[0] = v
	v.DstX = float32(topLeft.X)
	v.DstY = float32(bottomRight.Y)
	vertices[1] = v
	v.DstX = float32(bottomRight.X)
	v.DstY = float32(bottomRight.Y)
	vertices[2] = v
	v.DstX = float32(bottomRight.X)
	v.DstY = float32(topLeft.Y)
	vertices[3] = v

	opt := ebiten.DrawTrianglesOptions{}
	screen.DrawTriangles(vertices, idxs, drawing.WhitePixel, &opt)
}

func (s *negotiationGameScene) drawBackground(screen *ebiten.Image) {
	screen.Fill(color.Gray{Y: 96})

	v := ebiten.Vertex{
		ColorR: 0,
		ColorG: 0,
		ColorB: 0,
		ColorA: 0.25,
	}

	topLeft := geom.PointF{
		X: s.StagePos.X,
		Y: s.StagePos.Y,
	}

	bottomRight := geom.PointF{
		X: s.StagePos.X + s.Negotiation.Size.X,
		Y: s.StagePos.Y + s.Negotiation.Size.Y,
	}

	s.drawRect(screen, topLeft, bottomRight, v)
}

func (s *negotiationGameScene) drawApprovedEquips(screen *ebiten.Image) {
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
		opt.ColorScale.Scale(1, 0.5, 0, 1)
		opt.GeoM.Translate(topLeft.X, y)

		drawing.DrawText(screen, e.Name, 12, &opt)
	}
}

func (s *negotiationGameScene) drawMoney(screen *ebiten.Image) {
	drawing.DrawText(screen, fmt.Sprint(s.Negotiation.Money), 18, &ebiten.DrawImageOptions{})
}

func (s *negotiationGameScene) drawVendors(screen *ebiten.Image) {
	size := geom.PointF{X: 128, Y: 128}

	vert := ebiten.Vertex{
		ColorR: 1,
		ColorG: 1,
		ColorB: 0,
		ColorA: 0.25,
	}

	for i, v := range s.Negotiation.VendorSelector.Vendors() {
		bottomCenter := s.Negotiation.ProposalStartPosition(i)
		bottomCenter = bottomCenter.Add(s.StagePos)
		topLeft := geom.PointF{
			X: bottomCenter.X - 0.5*size.X,
			Y: bottomCenter.Y - size.Y,
		}

		bottomRight := geom.PointF{
			X: bottomCenter.X + 0.5*size.X,
			Y: bottomCenter.Y,
		}

		s.drawRect(screen, topLeft, bottomRight, vert)

		iopt := ebiten.DrawImageOptions{}
		iopt.GeoM.Translate(topLeft.X, topLeft.Y)
		drawing.DrawText(screen, v.Name, 12, &iopt)
	}
}

func (s *negotiationGameScene) drawManagers(screen *ebiten.Image) {
	size := geom.PointF{X: 128, Y: 128}
	idxs := []uint16{0, 1, 2, 0, 2, 3}

	mY := s.Negotiation.ManagerY()
	for i, m := range s.Negotiation.Managers {
		mL, mR := s.Negotiation.ManagerXLeftRight(i)
		center := geom.PointF{X: (mL + mR) / 2, Y: mY}
		center = center.Add(s.StagePos)

		v := ebiten.Vertex{
			ColorR: 0.4,
			ColorG: 0.7,
			ColorB: 0.5,
			ColorA: 1,
		}

		vertices := make([]ebiten.Vertex, 4)
		v.DstX = float32(center.X - 0.5*size.X)
		v.DstY = float32(center.Y - 0.5*size.Y)
		vertices[0] = v
		v.DstX = float32(center.X - 0.5*size.X)
		v.DstY = float32(center.Y + 0.5*size.Y)
		vertices[1] = v
		v.DstX = float32(center.X + 0.5*size.X)
		v.DstY = float32(center.Y + 0.5*size.Y)
		vertices[2] = v
		v.DstX = float32(center.X + 0.5*size.X)
		v.DstY = float32(center.Y - 0.5*size.Y)
		vertices[3] = v

		topt := ebiten.DrawTrianglesOptions{
			Address: ebiten.AddressRepeat,
		}
		screen.DrawTriangles(vertices, idxs, drawing.WhitePixel, &topt)

		iopt := ebiten.DrawImageOptions{}
		iopt.GeoM.Translate(float64(vertices[0].DstX), float64(vertices[0].DstY))
		drawing.DrawText(screen, m.Name, 12, &iopt)
	}
}

func (s *negotiationGameScene) drawDecisionMaker(screen *ebiten.Image) {
	left := s.Negotiation.DecisionMaker.LeftX
	top, _ := s.Negotiation.DecisionMaker.LinearFn.Y(left)
	width := s.Negotiation.DecisionMaker.Width

	topLeft := geom.PointF{X: left, Y: top - 2}
	bottomRight := geom.PointF{X: left + width, Y: top + 2}
	topLeft = topLeft.Add(s.StagePos)
	bottomRight = bottomRight.Add(s.StagePos)

	v := ebiten.Vertex{
		ColorR: 0.4,
		ColorG: 0.7,
		ColorB: 0.5,
		ColorA: 1,
	}

	s.drawRect(screen, topLeft, bottomRight, v)
}

func (s *negotiationGameScene) drawProposals(screen *ebiten.Image) {
	idxs := []uint16{0, 1, 2, 0, 2, 3}

	for _, p := range s.Negotiation.Proposals {
		hit := p.Hit

		gm := ebiten.GeoM{}
		gm.Rotate(p.Rotate)
		gm.Translate(hit.Center.X, hit.Center.Y)
		gm.Translate(s.StagePos.X, s.StagePos.Y)
		var x, y float64
		v := ebiten.Vertex{
			ColorR: 0.5,
			ColorG: 0.5,
			ColorB: 0.7,
			ColorA: 1,
		}

		vertices := make([]ebiten.Vertex, 4)
		x, y = gm.Apply(-hit.Radius, -hit.Radius)
		v.DstX = float32(x)
		v.DstY = float32(y)
		vertices[0] = v
		x, y = gm.Apply(-hit.Radius, hit.Radius)
		v.DstX = float32(x)
		v.DstY = float32(y)
		vertices[1] = v
		x, y = gm.Apply(hit.Radius, hit.Radius)
		v.DstX = float32(x)
		v.DstY = float32(y)
		vertices[2] = v
		x, y = gm.Apply(hit.Radius, -hit.Radius)
		v.DstX = float32(x)
		v.DstY = float32(y)
		vertices[3] = v

		topt := ebiten.DrawTrianglesOptions{}
		screen.DrawTriangles(vertices, idxs, drawing.WhitePixel, &topt)

		iopt := ebiten.DrawImageOptions{}
		iopt.GeoM.Translate(float64(vertices[0].DstX), float64(vertices[0].DstY))
		drawing.DrawText(screen, p.ImageName(), 12, &iopt)
	}
}

func (s *negotiationGameScene) drawProposalDelay(screen *ebiten.Image) {
	if s.Negotiation.ProposalDelay == nil {
		return
	}
	v := ebiten.Vertex{
		ColorR: 0,
		ColorG: 0,
		ColorB: 1,
		ColorA: 0.5,
	}

	startPos := s.Negotiation.ProposalDelay.Hit.Center.Add(s.StagePos)
	triangleVerts := make([]ebiten.Vertex, 3)
	v.DstX = float32(startPos.X)
	v.DstY = float32(startPos.Y)
	triangleVerts[0] = v
	v.DstX = float32(startPos.X + 16)
	v.DstY = float32(startPos.Y - 32)
	triangleVerts[1] = v
	v.DstX = float32(startPos.X - 16)
	v.DstY = float32(startPos.Y - 32)
	triangleVerts[2] = v

	topt := ebiten.DrawTrianglesOptions{}
	screen.DrawTriangles(triangleVerts, []uint16{0, 1, 2}, drawing.WhitePixel, &topt)

	baloonTopLeft := geom.PointF{
		X: s.StagePos.X + 8,
		Y: startPos.Y - 32 - 64,
	}
	baloonBottomRight := geom.PointF{
		X: s.StagePos.X + s.Negotiation.Size.X - 8,
		Y: startPos.Y - 32,
	}
	s.drawRect(screen, baloonTopLeft, baloonBottomRight, v)

	iopt := ebiten.DrawImageOptions{}
	iopt.GeoM.Translate(baloonTopLeft.X, baloonTopLeft.Y)
	drawing.DrawText(screen, s.Negotiation.ProposalDelay.ImageName(), 12, &iopt)
}

func (s *negotiationGameScene) End() bool {
	return s.Negotiation.End()
}

func (s *negotiationGameScene) Reset() {
	s.Negotiation.Reset(10000)
}

func createVendors() []*nego.Vendor {
	pm := createProposals()

	vv := make([]*nego.Vendor, 0)
	vv = append(vv, nego.NewVendor(name.VendorSamuraiAvionics, selectProposals(pm, name.EquipLaserCannon, name.EquipSpaceMissile, name.EquipHarakiriSystem), rand.New(rndSrc())))
	vv = append(vv, nego.NewVendor(name.VendorSalamisIndustry, selectProposals(pm, name.EquipBarrier, name.EquipArmorPlate, name.EquipThermalExhaustPort), rand.New(rndSrc())))
	vv = append(vv, nego.NewVendor(name.VendorCultualVictoryCo, selectProposals(pm, name.EquipStonehenge, name.EquipSushiBar, name.EquipOperaHouse), rand.New(rndSrc())))

	return vv
}

func selectProposals(m map[string]*nego.Proposal, names ...string) []*nego.Proposal {
	pp := make([]*nego.Proposal, 0, len(names))
	for _, n := range names {
		pp = append(pp, m[n])
	}

	return pp
}

func createProposals() map[string]*nego.Proposal {
	hit := geom.Circle{Radius: 32}
	velocity := geom.PointF{Y: -1}

	m := make(map[string]*nego.Proposal)

	addEquip := func(name string, cost int) {
		m[name] = &nego.Proposal{
			Equip:    &nego.Equip{Name: name},
			Cost:     cost,
			Hit:      hit,
			Velocity: velocity,
		}
	}

	addEquip(name.EquipLaserCannon, 1000)
	addEquip(name.EquipSpaceMissile, 1000)
	addEquip(name.EquipHarakiriSystem, 1500)
	addEquip(name.EquipBarrier, 1000)
	addEquip(name.EquipArmorPlate, 1000)
	addEquip(name.EquipThermalExhaustPort, 500)
	addEquip(name.EquipStonehenge, 1000)
	addEquip(name.EquipSushiBar, 1000)
	addEquip(name.EquipOperaHouse, 3000)

	return m
}

func createManagers() []*nego.Manager {
	mm := make([]*nego.Manager, 0)
	mm = append(mm, nego.NewManager(
		name.ManagerMachSonic,
		&nego.ProposalProcessorAccelerate{Value: 2},
		&nego.ProposalProcessorStopRotate{},
		&nego.ProposalProcessorCustomImageName{ImageName: ""}))
	mm = append(mm, nego.NewManager(
		name.ManagerBirdiePat,
		&nego.ProposalProcessorReduceCost{Multiplier: 0.8},
		&nego.ProposalProcessorCustomImageName{ImageName: name.EquipImageGolf}))
	mm = append(mm, nego.NewManager(
		name.ManagerLongWinded,
		&nego.ProposalProcessorAccelerate{Value: 0.8},
		&nego.ProposalProcessorRotate{Value: 1},
		&nego.ProposalProcessorImprove{}))

	return mm
}

var rndCount byte
var chacha8Base [32]byte = [32]byte{
	3, 14, 15, 92, 65, 35, 89, 79,
	32, 38, 64, 26, 43, 38, 32, 79,
	50, 28, 84, 19, 71, 69, 39, 93,
	75, 10, 58, 20, 97, 49, 44, 59,
}

func rndSrc() rand.Source {
	r := byte(time.Now().UnixNano() % 256)
	r += rndCount
	rndCount++
	c8src := [32]byte{}
	for i := range c8src {
		c8src[i] = chacha8Base[i] + r
	}

	return rand.NewChaCha8(c8src)
}
