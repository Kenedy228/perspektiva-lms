package attempt

type Status string

const (
	StatusFinished    Status = "finished"
	StatusExpired     Status = "expired"
	StatusInProgress  Status = "in_progress"
	StatusInterrupted Status = "interrupted"
	StatusCancelled   Status = "cancelled"
)
