package nego

import (
	"fmt"
	"image/color"
	"math"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/noppikinatta/ebitenginejam03/build"
	"github.com/noppikinatta/ebitenginejam03/data"
	"github.com/noppikinatta/ebitenginejam03/drawing"
	"github.com/noppikinatta/ebitenginejam03/geom"
	"github.com/noppikinatta/ebitenginejam03/lang"
	"github.com/noppikinatta/ebitenginejam03/name"
	"github.com/noppikinatta/ebitenginejam03/nego"
	"github.com/noppikinatta/ebitenginejam03/random"
	"github.com/noppikinatta/ebitenginejam03/scene"
)

type negotiationGameScene struct {
	Negotiation *nego.Negotiation
	StagePos    geom.PointF
	equipDesc   *build.EquipDescriptor
	moneyDrawer *drawing.GaugeDrawer
	preprocess  scene.Scene
	postprocess scene.Scene
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
	cmax.Scale(0.75, 0.75, 0.75, 1)

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
			TopLeft:       geom.PointF{X: 0, Y: 0},
			BottomRight:   geom.PointF{X: 600, Y: 40},
			TextOffset:    geom.PointF{X: 46, Y: 6},
			FontSize:      18,
			ColorScaleMax: cmax,
			ColorScaleMin: cmin,
		},
		preprocess: scene.NewContainer(
			scene.NewFadeIn(15),
			scene.NewShowImageScene(0, &ruleDescriptionDrawer{}),
			&scene.ScrollText{TextKey: name.TextKeyNegotiationTitle1},
		),
		postprocess: scene.NewContainer(
			&scene.ScrollText{TextKey: name.TextKeyNegotiationTitle2},
			scene.NewFadeOut(15),
		),
	}
}

func (s *negotiationGameScene) Update() error {
	x, _ := ebiten.CursorPosition()
	x -= int(s.StagePos.X)
	s.Negotiation.UpdateDecisionMaker(float64(x))
	if !s.preprocess.End() {
		return s.preprocess.Update()
	}
	if s.Negotiation.End() {
		return s.postprocess.Update()
	}
	s.Negotiation.UpdateOthers()
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
	if !s.preprocess.End() {
		s.preprocess.Draw(screen)
	}
	if s.Negotiation.End() && !s.postprocess.End() {
		s.postprocess.Draw(screen)
	}
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
	screen.Fill(color.Gray{Y: 24})

	v := ebiten.Vertex{
		ColorR: 0,
		ColorG: 0,
		ColorB: 0,
		ColorA: 1,
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
		opt.GeoM.Translate(topLeft.X, topLeft.Y+y)

		img := drawing.Image(name.ImgKey(e.Name))
		screen.DrawImage(img, &opt)
		opt.GeoM.Translate(equipNameLeft, 8)

		name := lang.Text(e.Name)
		if e.ImprovedCount > 0 {
			name += fmt.Sprintf("+%d", e.ImprovedCount)
		}
		drawing.DrawText(screen, name, 12, &opt)
	}
}

func (s *negotiationGameScene) drawMoney(screen *ebiten.Image) {
	s.moneyDrawer.Draw(screen)
	coinImg := drawing.Image(name.ImgKeyCoin)
	opt := ebiten.DrawImageOptions{}
	opt.GeoM.Translate(4, 4)
	screen.DrawImage(coinImg, &opt)
}

func (s *negotiationGameScene) drawVendors(screen *ebiten.Image) {
	for i, v := range s.Negotiation.VendorSelector.Vendors() {
		img := drawing.Image(name.ImgKey(v.Name))
		imgSize := geom.PointFFromPoint(img.Bounds().Size())

		bottomCenter := s.Negotiation.ProposalStartPosition(i)
		bottomCenter = bottomCenter.Add(s.StagePos)

		opt := ebiten.DrawImageOptions{}
		opt.GeoM.Translate(bottomCenter.X, bottomCenter.Y-24)
		opt.GeoM.Translate(-imgSize.X*0.5, -imgSize.Y)
		opt.ColorScale = data.ColorTheme(v.Name)
		opt.ColorScale.ScaleAlpha(0.5)

		screen.DrawImage(img, &opt)

		opt.GeoM.Translate(0, imgSize.Y)
		opt.ColorScale = ebiten.ColorScale{}
		drawing.DrawTextByKey(screen, v.Name, 12, &opt)
	}
}

func (s *negotiationGameScene) drawManagers(screen *ebiten.Image) {
	mY := s.Negotiation.ManagerY()
	for i, m := range s.Negotiation.Managers {
		img := drawing.Image(name.ImgKey(m.Name))
		imgSize := geom.PointFFromPoint(img.Bounds().Size())

		mL, mR := s.Negotiation.ManagerXLeftRight(i)
		center := geom.PointF{X: (mL + mR) / 2, Y: mY}
		center = center.Add(s.StagePos)

		opt := ebiten.DrawImageOptions{}
		opt.GeoM.Translate(center.X-imgSize.X*0.5, center.Y-imgSize.Y*0.5)
		opt.ColorScale = data.ColorTheme(m.Name)
		opt.ColorScale.ScaleAlpha(0.5)
		screen.DrawImage(img, &opt)

		opt.ColorScale = ebiten.ColorScale{}
		opt.GeoM.Translate(0, imgSize.Y)
		drawing.DrawTextByKey(screen, m.Name, 12, &opt)
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
	for _, p := range s.Negotiation.Proposals {
		img := drawing.Image(p.ImageName())
		imgSize := geom.PointFFromPoint(img.Bounds().Size())

		hit := p.Hit

		gm := ebiten.GeoM{}
		gm.Translate(-imgSize.X*0.5, -imgSize.Y*0.5)
		gm.Rotate(p.Rotate)
		gm.Translate(hit.Center.X, hit.Center.Y)
		gm.Translate(s.StagePos.X, s.StagePos.Y)

		opt := ebiten.DrawImageOptions{}
		opt.GeoM = gm

		screen.DrawImage(img, &opt)
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
		ColorB: 0,
		ColorA: 0.75,
	}

	if s.Negotiation.LastSelectedVendor != nil {
		vendor := s.Negotiation.LastSelectedVendor
		cs := data.ColorTheme(vendor.Name)
		v.ColorR = cs.R() * v.ColorA
		v.ColorG = cs.G() * v.ColorA
		v.ColorB = cs.B() * v.ColorA
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
	iopt.GeoM.Translate(baloonTopLeft.X, baloonTopLeft.Y)
	img := drawing.Image(s.Negotiation.ProposalDelay.ImageName())
	screen.DrawImage(img, &iopt)

	iopt.GeoM.Translate(equipNameLeft, equipNameTop)
	eqpName := s.Negotiation.ProposalDelay.Equip.Name
	drawing.DrawTextByKey(screen, eqpName, 14, &iopt)

	coinImg := drawing.Image(name.ImgKeyCoin)
	coinImgSize := geom.PointFFromPoint(coinImg.Bounds().Size())
	coinGm := ebiten.GeoM{}
	coinGm.Translate(-coinImgSize.X*0.5, -coinImgSize.Y*0.5)
	coinGm.Scale(0.5, 0.5)
	coinGm.Translate(0, coinImgSize.Y*0.25)

	copt := iopt
	copt.GeoM.Translate(200, 4)
	coinGm.Concat(copt.GeoM)
	copt.GeoM = coinGm
	screen.DrawImage(coinImg, &copt)

	copt = iopt
	copt.GeoM.Translate(200+coinImgSize.X*0.5, 4)
	drawing.DrawText(screen, fmt.Sprint(s.Negotiation.ProposalDelay.Cost), 12, &copt)

	iopt.GeoM.Translate(0, float64(equipDescTop))
	eqpDesc := name.DescKey(eqpName)
	tmplData := s.equipDesc.TemplateData(eqpName, 0)
	drawing.DrawTextTemplate(screen, eqpDesc, tmplData, 12, &iopt)
}

func (s *negotiationGameScene) End() bool {
	return s.postprocess.End()
}

func (s *negotiationGameScene) Reset() {
	s.preprocess.Reset()
	s.Negotiation.Reset(10000)
	s.postprocess.Reset()
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
		&nego.ProposalProcessorCustomImageName{ImageName: name.ImgKeyGolf}))
	mm = append(mm, nego.NewManager(
		name.TextKeyManager3,
		&nego.ProposalProcessorRotate{Value: 0.25},
		&nego.ProposalProcessorImprove{}))

	return mm
}

type ruleDescriptionDrawer struct {
}

func (d *ruleDescriptionDrawer) Draw(screen *ebiten.Image) {
	size := screen.Bounds().Size()
	clr := color.RGBA{A: 128}
	vector.DrawFilledRect(screen, 0, 0, float32(size.X), float32(size.Y), clr, false)

	opt := ebiten.DrawImageOptions{}
	opt.GeoM.Translate(24, 240)
	drawing.DrawTextByKey(screen, name.TextKeyNegotiationDesc1, 18, &opt)
}
