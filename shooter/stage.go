package shooter

import "github.com/noppikinatta/ebitenginejam03/geom"

type Stage struct {
	Size          geom.PointF
	firstAngle    float64
	MyShip        *MyShip
	EnemyLauncher *EnemyLauncher
	HitTest       *HitTest
}

func (s *Stage) Update(cursorPos geom.PointF) {
	angle := s.calcAngle(cursorPos)
	s.updateMyShipAngle(angle)
	s.updateMyShip()
	s.updateEnemies()
	s.updateHitTest()
}

func (s *Stage) calcAngle(cursorPos geom.PointF) float64 {
	centerMinus := s.Size.Multiply(-0.5)
	relative := cursorPos.Add(centerMinus)
	angle := relative.Angle()

	if s.firstAngle == 0 {
		// s.firstAngle is sometimes set to 0, but the next frame angle will probably have a non-zero value, so it's probably ok
		s.firstAngle = angle
	}
	return angle - s.firstAngle
}

func (s *Stage) updateMyShipAngle(angle float64) {
	s.MyShip.UpdateAngle(angle)
}

func (s *Stage) updateMyShip() {
	s.MyShip.Update()
}

func (s *Stage) updateEnemies() {
	s.EnemyLauncher.Update()
}

func (s *Stage) updateHitTest() {
	s.HitTest.Update()
}

func (s *Stage) End() bool {
	if s.EnemyLauncher.Annihilated {
		return true
	}
	if !s.MyShip.IsLiving() {
		return true
	}
	return false
}

func (s *Stage) Won() bool {
	return s.End() && s.MyShip.IsLiving()
}
