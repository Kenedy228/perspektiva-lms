package title_test

import (
	"strings"
	"testing"

	"gitflic.ru/lms/internal/domain/question/content"
	"gitflic.ru/lms/internal/domain/question/title"
	"github.com/stretchr/testify/assert"
)

func TestNew_Success(t *testing.T) {
	tc := []struct {
		name  string
		cType content.Type
		value string
	}{
		{
			name:  "валидный кейс",
			cType: content.TypeText,
			value: "Что изображено на картинке?",
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			//Arrange
			got, err := title.New(makeContent(tt.cType, tt.value))

			//Assert
			assert.NoError(t, err)
			assert.Equal(t, got.Value(), tt.value)
		})
	}
}

func TestNew_Fail(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		cType   content.Type
		value   string
		wantErr error
	}{
		{
			name:    "тип контента не текст",
			cType:   content.TypeImage,
			value:   "image.png",
			wantErr: title.ErrInvalid,
		},
		{
			name:    "значение контента выходит за пределы количества символов",
			cType:   content.TypeText,
			value:   strings.Repeat("a", 1e5),
			wantErr: title.ErrInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//Arrange
			_, err := title.New(makeContent(tt.cType, tt.value))

			//Assert
			assert.Error(t, err)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
