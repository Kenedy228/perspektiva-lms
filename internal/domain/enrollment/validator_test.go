package enrollment

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestEnrollment_Status(t *testing.T) {
	start := time.Date(2025, 5, 10, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 0, 10)

	e := &Enrollment{
		id:              uuid.New(),
		courseID:        uuid.New(),
		courseVersionID: uuid.New(),
		accountID:       uuid.New(),
		activatedAt:     start,
		deactivatedAt:   end,
	}

	t.Run("inactive before activation", func(t *testing.T) {
		//Act
		before := start.AddDate(0, 0, -1)

		//Assert
		assert.Equal(t, StatusInactive, e.Status(before))
		assert.False(t, e.IsActive(before))
	})

	t.Run("active on start date", func(t *testing.T) {
		//Assert
		assert.Equal(t, StatusActive, e.Status(start))
		assert.True(t, e.IsActive(start))
	})

	t.Run("active between", func(t *testing.T) {
		//Act
		mid := start.AddDate(0, 0, 5)

		//Assert
		assert.Equal(t, StatusActive, e.Status(mid))
		assert.True(t, e.IsActive(mid))
	})

	t.Run("active on end date", func(t *testing.T) {
		//Assert
		assert.Equal(t, StatusActive, e.Status(end))
		assert.True(t, e.IsActive(end))
	})

	t.Run("expired after end date", func(t *testing.T) {
		//Act
		after := end.AddDate(0, 0, 1)

		//Assert
		assert.Equal(t, StatusExpired, e.Status(after))
		assert.False(t, e.IsActive(after))
	})
}
