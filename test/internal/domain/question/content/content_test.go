//go:build legacy
// +build legacy

package content_test

import (
	"testing"

	"gitflic.ru/lms/backend/internal/domain/question/content"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Run("успех", func(t *testing.T) {
		tc := []struct{
			name string
			cType content.Type
			cValue string
		} {
			{
				name: "валидный текст",
				cType: content.TypeText,
				cValue: "текст",
			},
			{
				name: "валидное изображение",
				cType: content.TypeImage,
				cValue: "изображение",
			},
			{
				name: ""
			}
		}
	})
}

func TestNew_Success(t *testing.T) {
	//Arrange
	c, err := content.New(content.TypeText, "text")

	//Assert
	assert.NoError(t, err)
	assert.Equal(t, content.TypeText, c.Type())
	assert.Equal(t, "text", c.Value())
}

func TestNew_Fail(t *testing.T) {
	//Arrange
	_, err := content.New(content.TypeText, "")

	//Assert
	assert.Error(t, err)
	assert.ErrorIs(t, err, content.ErrInvalid)
}
