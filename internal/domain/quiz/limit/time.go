package limit

import "time"

type Time struct {
	seconds int
}

func NewTime(seconds int) (Time, error) {
	if err := validateSeconds(seconds); err != nil {
		return Time{}, err
	}

	return Time{
		seconds: seconds,
	}, nil
}

func (t Time) IsInfinite() bool {
	return t.seconds == 0
}

func (t Time) Seconds() int {
	return t.seconds
}

func (t Time) TryDuration() (time.Duration, bool) {
	if t.seconds == 0 {
		return 0, false
	}

	return time.Second * time.Duration(t.seconds), true
}
