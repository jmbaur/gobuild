package builder

type Status int

const (
	StatusQueued = iota
	StatusRunning
	StatusFailed
	StatusFinished
)

func (s Status) String() string {
	return []string{"Queued", "Running", "Failed", "Finished"}[s]
}
