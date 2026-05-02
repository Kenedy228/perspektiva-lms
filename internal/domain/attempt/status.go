package attempt

type Status string

const (
	StatusFinished    Status = "finished"
	StatusExpired     Status = "expired"
	StatusInProgress  Status = "in_progress"
	StatusInterrupted Status = "interrupted"
	StatusCancelled   Status = "cancelled"
)

func (s Status) Title() string {
	switch s {
	case StatusFinished:
		return "завершен"
	case StatusExpired:
		return "просрочен"
	case StatusInProgress:
		return "в процессе"
	case StatusInterrupted:
		return "прерван"
	case StatusCancelled:
		return "отменен"
	default:
		return ""
	}
}

func (s Status) String() string {
	return string(s)
}
