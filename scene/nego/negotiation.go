package nego

import (
	"fmt"
	"image/color"
	"math"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/noppikinatta/ebitenginejam03/build"
	"github.com/noppikinatta/ebitenginejam03/drawing"
	"github.com/noppikinatta/ebitenginejam03/geom"
	"github.com/noppikinatta/ebitenginejam03/lang"
	"github.com/noppikinatta/ebitenginejam03/name"
	"github.com/noppikinatta/ebitenginejam03/nego"
	"github.com/noppikinatta/ebitenginejam03/random"
)

type negotiationGameScene struct {
	Negotiation *nego.Negotiation
	StagePos    geom.PointF
	equipDesc   *build.EquipDescriptor
	moneyDrawer *drawing.GaugeDrawer
}

func newNegotiationGameScene() *negotiationGameScene {
	n := nego.Negotiation{
		Size:          geom.PointF{X: 600, Y: 600},
		DecisionMaker: nego.NewDecisionMaker(0, 100),
		VendorSelector: nego.NewVendorSelector(
			createVendors(),
			180,
			rand.New(random.Source())),
		Managers: createManagers(),
		Money:    10000,
	}

	cmax := ebiten.ColorScale{}
	cmax.Scale(0.75, 0.75, 0.75, 0.75)

	cmin := ebiten.ColorScale{}
	cmin.SetG(0.25)
	cmin.SetB(0.25)

	return &negotiationGameScene{
		Negotiation: &n,
		StagePos:    geom.PointF{X: 0, Y: 40},
		equipDesc:   build.NewEquipDescriptor(),
		moneyDrawer: &drawing.GaugeDrawer{
			Current:       n.Money,
			Max:           n.Money,
			TopLeft:       geom.PointF{X: 40, Y: 0},
			BottomRight:   geom.PointF{X: 600, Y: 40},
			TextOffset:    geom.PointF{X: 6, Y: 6},
			FontSize:      18,
			ColorScaleMax: cmax,
			ColorScaleMin: cmin,
		},
	}
}

func (s *negotiationGameScene) Update() error {
	x, _ := ebiten.CursorPosition()
	x -= int(s.StagePos.X)
	s.Negotiation.Update(float64(x))
	s.moneyDrawer.Current = s.Negotiation.Money
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
	opt := ebiten.DrawImageOptions{}
	opt.GeoM.Translate(604, 8)
	drawing.DrawTextByKey(screen, name.TextKeyNegotiationTitle3, 14, &opt)

	if len(s.Negotiation.ApprovedEquips) == 0 {
		return
	}

	const (
		equipNameLeft float64 = 68
		padding       float64 = 2
		itemheightMax float64 = 80
	)

	topLeft := geom.PointF{}
	topLeft.X = s.Negotiation.Size.X + padding
	topLeft.Y = s.StagePos.Y

	rightBottom := geom.PointF{}
	rightBottom.X = float64(screen.Bounds().Max.X)
	rightBottom.Y = float64(screen.Bounds().Max.Y)

	areaHeight := rightBottom.Y - topLeft.Y
	itemHeight := math.Min(itemheightMax, areaHeight/float64(len(s.Negotiation.ApprovedEquips)))

	for i, e := range s.Negotiation.ApprovedEquips {
		y := itemHeight * float64(i)

		opt := ebiten.DrawImageOptions{}
		opt.GeoM.Translate(topLeft.X+equipNameLeft, topLeft.Y+y)

		name := lang.Text(e.Name)
		if e.ImprovedCount > 0 {
			name += fmt.Sprintf("+%d", e.ImprovedCount)
		}

		drawing.DrawText(screen, name, 12, &opt)
	}
}

func (s *negotiationGameScene) drawMoney(screen *ebiten.Image) {
	s.moneyDrawer.Draw(screen)
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
		iopt.GeoM.Translate(topLeft.X, bottomRight.Y-20)
		drawing.DrawTextByKey(screen, v.Name, 12, &iopt)
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
		iopt.GeoM.Translate(float64(vertices[1].DstX), float64(vertices[1].DstY)-20)
		drawing.DrawTextByKey(screen, m.Name, 12, &iopt)
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

	const (
		triangleSize  float64 = 32
		baloonPadding float64 = 8
		baloonHeight  float64 = 80
		equipNameTop  float64 = 4
		equipNameLeft float64 = 80
		equipDescTop  float32 = 32
	)

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
	v.DstX = float32(startPos.X + triangleSize/2)
	v.DstY = float32(startPos.Y - triangleSize)
	triangleVerts[1] = v
	v.DstX = float32(startPos.X - triangleSize/2)
	v.DstY = float32(startPos.Y - triangleSize)
	triangleVerts[2] = v

	topt := ebiten.DrawTrianglesOptions{}
	screen.DrawTriangles(triangleVerts, []uint16{0, 1, 2}, drawing.WhitePixel, &topt)

	baloonTopLeft := geom.PointF{
		X: s.StagePos.X + baloonPadding,
		Y: startPos.Y - triangleSize - baloonHeight,
	}
	baloonBottomRight := geom.PointF{
		X: s.StagePos.X + s.Negotiation.Size.X - baloonPadding,
		Y: startPos.Y - triangleSize,
	}
	s.drawRect(screen, baloonTopLeft, baloonBottomRight, v)

	iopt := ebiten.DrawImageOptions{}
	iopt.GeoM.Translate(baloonTopLeft.X+equipNameLeft, baloonTopLeft.Y+equipNameTop)
	eqpName := s.Negotiation.ProposalDelay.Equip.Name
	drawing.DrawTextByKey(screen, eqpName, 14, &iopt)

	iopt.GeoM.Translate(0, float64(equipDescTop))
	eqpDesc := name.DescKey(eqpName)
	tmplData := s.equipDesc.TemplateData(eqpName, 0)
	drawing.DrawTextTemplate(screen, eqpDesc, tmplData, 12, &iopt)
}

func (s *negotiationGameScene) End() bool {
	return s.Negotiation.End()
}

func (s *negotiationGameScene) Reset() {
	s.Negotiation.Reset(10000)
}

func (s *negotiationGameScene) Result() []*nego.Equip {
	return s.Negotiation.ApprovedEquips
}

func createVendors() []*nego.Vendor {
	pm := createProposals()

	vv := make([]*nego.Vendor, 0)
	vv = append(vv, nego.NewVendor(
		name.TextKeyVendor1,
		selectProposals(pm, name.TextKeyEquip1Laser, name.TextKeyEquip2Missile, name.TextKeyEquip3Harakiri),
		rand.New(random.Source())))
	vv = append(vv, nego.NewVendor(
		name.TextKeyVendor2,
		selectProposals(pm, name.TextKeyEquip4Barrier, name.TextKeyEquip5Armor, name.TextKeyEquip6Exhaust),
		rand.New(random.Source())))
	vv = append(vv, nego.NewVendor(
		name.TextKeyVendor3,
		selectProposals(pm, name.TextKeyEquip7Stonehenge, name.TextKeyEquip8Sushibar, name.TextKeyEquip9Operahouse),
		rand.New(random.Source())))

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

	addEquip(name.TextKeyEquip1Laser, 1000)
	addEquip(name.TextKeyEquip2Missile, 1000)
	addEquip(name.TextKeyEquip3Harakiri, 1500)
	addEquip(name.TextKeyEquip4Barrier, 1000)
	addEquip(name.TextKeyEquip5Armor, 1000)
	addEquip(name.TextKeyEquip6Exhaust, 500)
	addEquip(name.TextKeyEquip7Stonehenge, 1000)
	addEquip(name.TextKeyEquip8Sushibar, 1000)
	addEquip(name.TextKeyEquip9Operahouse, 3000)

	return m
}

func createManagers() []*nego.Manager {
	mm := make([]*nego.Manager, 0)
	mm = append(mm, nego.NewManager(
		name.TextKeyManager1,
		&nego.ProposalProcessorAccelerate{Value: 2},
		&nego.ProposalProcessorStopRotate{},
		&nego.ProposalProcessorCustomImageName{ImageName: ""}))
	mm = append(mm, nego.NewManager(
		name.TextKeyManager2,
		&nego.ProposalProcessorReduceCost{Multiplier: 0.8},
		&nego.ProposalProcessorCustomImageName{ImageName: name.EquipImageGolf}))
	mm = append(mm, nego.NewManager(
		name.TextKeyManager3,
		&nego.ProposalProcessorRotate{Value: 0.25},
		&nego.ProposalProcessorImprove{}))

	return mm
}
