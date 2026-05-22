//go:build legacy
// +build legacy

package blank_test

import (
	"testing"

	"gitflic.ru/lms/backend/internal/domain/question/typed/blank"
	"gitflic.ru/lms/backend/internal/domain/question/content"
	"github.com/stretchr/testify/assert"
)

func TestNew_Success(t *testing.T) {
	//Arrange
	b, err := blank.New("{{placeholder}}", makeContentSlice(5))

	//Assert
	assert.NoError(t, err)
	assert.Equal(t, 5, len(b.Variants()))
	assert.Equal(t, "{{placeholder}}", b.Placeholder())
}

func TestNew_Fail(t *testing.T) {
	t.Run("плейсхолдер должен соответствовать маске", func(t *testing.T) {
		//Arrange
		_, err := blank.New("placeholder", makeContentSlice(5))

		//Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, blank.ErrInvalid)
	})

	t.Run("количество бланков меньше нужного количества", func(t *testing.T) {
		//Arrange
		_, err := blank.New("{{placeholder}}", nil)

		//Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, blank.ErrInvalid)
	})

	t.Run("количество бланков больше нужного количества", func(t *testing.T) {
		//Arrange
		_, err := blank.New("{{placeholder}}", makeContentSlice(1e3))

		//Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, blank.ErrInvalid)
	})

	t.Run("бланк нетекстового формата", func(t *testing.T) {
		//Arrange
		_, err := blank.New("{{placeholder}}", []content.Content{makeContent(content.TypeAudio, "audio")})

		//Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, blank.ErrInvalid)
	})
}
