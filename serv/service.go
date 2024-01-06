package serv

type Service struct {
	Status Status
	Exit   chan bool
}

type Status int

const (
	RUNNING Status = iota
	ERR
	STOPPED
	STARTING
)

func (s Status) String() string {
	switch s {
	case RUNNING:
		return "RUNNING"
	case ERR:
		return "ERR"
	case STOPPED:
		return "STOPPED"
	case STARTING:
		return "STARTING"
	default:
		return "INVALID"
	}
}
