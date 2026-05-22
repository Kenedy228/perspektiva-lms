//go:build legacy
// +build legacy

package enrollment_test

import (
	"testing"
	"testing/synctest"
	"time"

	"gitflic.ru/lms/backend/internal/domain/enrollment"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewEnrollment_Success(t *testing.T) {
	t.Run("basic creation", func(t *testing.T) {
		//Arrange
		e := newEnrollmentBuilder().build(t, nil)
		gotFrom := e.ActivatedAt()
		gotTo := e.DeactivatedAt()

		//Assert
		assert.NotEqual(t, uuid.Nil, e.ID())
		assert.NotEqual(t, uuid.Nil, e.CourseID())
		assert.NotEqual(t, uuid.Nil, e.CourseVersionID())
		assert.NotEqual(t, uuid.Nil, e.AccountID())

		assert.False(t, e.CreatedAt().IsZero())
		assert.False(t, e.UpdatedAt().IsZero())
		assert.Equal(t, e.CreatedAt(), e.UpdatedAt())

		assert.Equal(t, 0, gotFrom.Hour())
		assert.Equal(t, 0, gotFrom.Minute())
		assert.Equal(t, 0, gotFrom.Second())
		assert.Equal(t, 0, gotFrom.Nanosecond())

		assert.Equal(t, 0, gotTo.Hour())
		assert.Equal(t, 0, gotTo.Minute())
		assert.Equal(t, 0, gotTo.Second())
		assert.Equal(t, 0, gotTo.Nanosecond())
	})
}

func TestEnrollment_StatusAndIsActive(t *testing.T) {
	today := time.Now()
	start := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 0, 10)

	e := newEnrollmentBuilder().
		withWindow(start, end).
		build(t, nil)

	t.Run("inactive before activation", func(t *testing.T) {
		//Act
		before := start.AddDate(0, 0, -1)
		status := e.Status(before)

		//Assert
		assert.Equal(t, enrollment.StatusInactive, status)
		assert.False(t, e.IsActive(before))
	})

	t.Run("active during window", func(t *testing.T) {
		//Act
		mid := start.AddDate(0, 0, 5)
		status := e.Status(mid)

		//Assert
		assert.Equal(t, enrollment.StatusActive, status)
		assert.True(t, e.IsActive(mid))
	})

	t.Run("expired after deactivation", func(t *testing.T) {
		//Act
		after := end.AddDate(0, 0, 1)
		status := e.Status(after)

		//Assert
		assert.Equal(t, enrollment.StatusExpired, status)
		assert.False(t, e.IsActive(after))
	})
}

func TestEnrollment_ChangeActivationWindow(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		today := time.Now()
		start := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, time.UTC)
		end := start.AddDate(0, 0, 10)

		e := newEnrollmentBuilder().
			withWindow(start, end).
			build(t, nil)

		initialUpdatedAt := e.UpdatedAt()

		time.Sleep(time.Second * 10)

		newStart := start.AddDate(0, 0, 5)
		newEnd := newStart.AddDate(0, 0, 5)

		err := e.ChangeActivationWindow(newStart, newEnd)
		require.NoError(t, err)

		assert.True(t, e.UpdatedAt().After(initialUpdatedAt))

		before := start.AddDate(0, 0, 1)
		assert.Equal(t, enrollment.StatusInactive, e.Status(before))
	})
}
