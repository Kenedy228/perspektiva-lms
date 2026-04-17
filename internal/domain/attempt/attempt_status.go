package attempt

type Status string

const (
	StatusSubmitted  Status = "submitted"
	StatusExpired    Status = "expired"
	StatusInProgress Status = "in_progress"
	StatusFailed     Status = "failed"
)
