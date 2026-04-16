package dob

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Run("should not return err", func(t *testing.T) {
		now := time.Now()
		dob, err := NewDateOfBirth(now)

		require.NoError(t, err)

		assert.Equal(t, now, dob.Date())
	})

	t.Run("should throw err", func(t *testing.T) {
		now := time.Now().Add(time.Hour)

		_, err := NewDateOfBirth(now)

		assert.ErrorIs(t, err, ErrInvalidDob)
	})
}

func TestAgeAt(t *testing.T) {
	t.Run("valid date", func(t *testing.T) {
		dobTime := time.Date(2000, 2, 29, 0, 0, 0, 0, time.UTC)

		d, err := NewDateOfBirth(dobTime)
		require.NoError(t, err, "failed to setup date of birth")

		age, err := d.AgeAt(time.Date(2026, 4, 16, 0, 0, 0, 0, time.UTC))

		assert.NoError(t, err)
		assert.Equal(t, 26, age)
	})

	t.Run("invalid date", func(t *testing.T) {
		now := time.Now()

		d, err := NewDateOfBirth(now)
		require.NoError(t, err, "failed to setup date of birth")

		at := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

		_, err = d.AgeAt(at)
		assert.ErrorIs(t, err, ErrInvalidAt)
	})
}
