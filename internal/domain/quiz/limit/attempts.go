package limit

type Attempts struct {
	count int
}

func NewAttempts(count int) (Attempts, error) {
	if err := validateAttemptsCount(count); err != nil {
		return Attempts{}, err
	}

	return Attempts{
		count: count,
	}, nil
}

func (a Attempts) Count() int {
	return a.count
}

func (a Attempts) IsInfinite() bool {
	return a.count == 0
}
