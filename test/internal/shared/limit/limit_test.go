package limit_test

import (
	"testing"
	"time"

	"gitflic.ru/lms/internal/domain/shared/limit"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Run("should return error, if limit is negative", func(t *testing.T) {
		newLimitBuilder().withValue(-1).build(t, limit.ErrInvalidLimit)
	})

	t.Run("should return error, if limit is greater than 1 year", func(t *testing.T) {
		newLimitBuilder().withValue(365*24*60*60+1).build(t, limit.ErrInvalidLimit)
	})

	t.Run("should create infinite limit if value is 0", func(t *testing.T) {
		limit := newLimitBuilder().withValue(0).build(t, nil)

		assert.Equal(t, limit.Value(), 0)
		assert.True(t, limit.IsInfinite())
	})

	t.Run("should create finite limit if value is not 0", func(t *testing.T) {
		limit := newLimitBuilder().withValue(10).build(t, nil)

		assert.Equal(t, limit.Value(), 10)
		assert.False(t, limit.IsInfinite())
	})
}

func TestDuration(t *testing.T) {
	t.Run("should return time.Duration value", func(t *testing.T) {
		limit := newLimitBuilder().withValue(100).build(t, nil)

		assert.Equal(t, time.Duration(time.Second*100), limit.Duration())
	})
}
