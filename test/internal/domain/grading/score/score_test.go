package score_test

import (
	"testing"

	"gitflic.ru/lms/internal/domain/grading/score"
	"github.com/stretchr/testify/assert"
)

func TestNew_Success(t *testing.T) {
	tc := []struct {
		name  string
		value float64
	}{
		{
			name:  "значение больше 0",
			value: 100,
		},
		{
			name:  "значение 0",
			value: 0,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			//Arrange
			s, err := score.New(tt.value)

			//Assert
			assert.NoError(t, err)
			assert.Equal(t, tt.value, s.Value())
		})
	}
}

func TestNew_Fail(t *testing.T) {
	tc := []struct {
		name  string
		value float64
	}{
		{
			name:  "значение меньше 0",
			value: -1,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			//Arrange
			_, err := score.New(tt.value)

			//Assert
			assert.Error(t, err)
			assert.ErrorIs(t, err, score.ErrInvalid)
		})
	}
}
