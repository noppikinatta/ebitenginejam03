package shooter

type State int

const (
	StateReady State = iota
	StateOnStage
	StateDead
)
