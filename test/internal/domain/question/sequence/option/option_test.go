package option_test

import (
	"strings"
	"testing"

	"gitflic.ru/lms/internal/domain/question/content"
	"gitflic.ru/lms/internal/domain/question/sequence/option"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNew_Success(t *testing.T) {
	tc := []struct {
		name  string
		cType content.Type
		value string
	}{
		{
			name:  "текстовый контент, удовлетворяющий инвариантам",
			cType: content.TypeText,
			value: "текст варианта ответа",
		},
		{
			name:  "контент в виде изображения",
			cType: content.TypeImage,
			value: "изображение варианта ответа",
		},
		{
			name:  "контент в виде аудио",
			cType: content.TypeAudio,
			value: "аудио варианта ответа",
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			//Arrange
			opt, err := newOptionBuilder().
				withContent(tt.cType, tt.value).
				build()

			//Assert
			assert.NoError(t, err)
			assert.NotEqual(t, uuid.Nil, opt.ID())
			assert.Equal(t, tt.value, opt.Content().Value())
			assert.Equal(t, tt.cType, opt.Content().Type())
		})
	}
}

func TestNew_Fail(t *testing.T) {
	tc := []struct {
		name    string
		cType   content.Type
		value   string
		wantErr error
	}{
		{
			name:    "текстовый контент, с превышением количества символов",
			cType:   content.TypeText,
			value:   strings.Repeat("a", 1e5),
			wantErr: option.ErrInvalid,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			//Arrange
			_, err := newOptionBuilder().
				withContent(tt.cType, tt.value).
				build()

			//Assert
			assert.Error(t, err)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
