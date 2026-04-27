package timelimit_test

import (
	"gitflic.ru/lms/internal/domain/shared/timelimit"
)

type limitBuilder struct {
	seconds int
}

func newLimitBuilder() *limitBuilder {
	return &limitBuilder{
		seconds: -1,
	}
}

func (b *limitBuilder) withSeconds(seconds int) *limitBuilder {
	b.seconds = seconds
	return b
}

func (b *limitBuilder) build() (timelimit.Limit, error) {
	limit, err := timelimit.New(b.seconds)
	return limit, err
}
