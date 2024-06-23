package nego

import "math/rand/v2"

type randomSelector[T any] struct {
	Items         []T
	selectedCount []int
	rnd           *rand.Rand
}

func (s *randomSelector[T]) Reset() {
	for i := range s.selectedCount {
		s.selectedCount[i] = 0
	}
}

func (s *randomSelector[T]) Select() T {
	if len(s.selectedCount) == 0 {
		s.selectedCount = make([]int, len(s.Items))
	}

	var rndMax float64
	priprities := make([]float64, len(s.Items))
	for i := range priprities {
		p := s.priority(i)
		priprities[i] = p
		rndMax += p
	}

	rndValue := s.rnd.Float64() * rndMax
	for i := range s.Items {
		if rndValue < priprities[i] {
			return s.Items[i]

		}
		rndValue -= priprities[i]
	}

	// fallback
	return s.Items[len(s.Items)-1]
}

func (s *randomSelector[T]) priority(idx int) float64 {
	return 1.0 / float64(1+s.selectedCount[idx])
}
