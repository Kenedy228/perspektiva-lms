package attempt

type Status string

const (
	// NOTE: успешно выполнен
	StatusSubmitted Status = "submitted"
	// NOTE: не успел пройти вовремя
	StatusExpired Status = "expired"
	// NOTE: проходит
	StatusInProgress Status = "in_progress"
)
