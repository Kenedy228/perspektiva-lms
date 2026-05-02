package variant_test

import (
	"strings"
	"testing"

	"gitflic.ru/lms/internal/domain/question/content"
	"gitflic.ru/lms/internal/domain/question/short/variant"
	"github.com/stretchr/testify/assert"
)

func TestNew_Success(t *testing.T) {
	tc := []struct {
		name  string
		cType content.Type
		value string
	}{
		{
			name:  "корректный тип и количество символов",
			cType: content.TypeText,
			value: "значение опции",
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			//Arrange
			v, err := makeVariant(makeContent(tt.cType, tt.value))

			//Assert
			assert.NoError(t, err)
			assert.Equal(t, tt.value, v.Content().Value())
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
			name:    "любой тип, кроме текстового - невалидный",
			cType:   content.TypeImage,
			value:   "значение опции",
			wantErr: variant.ErrInvalid,
		},
		{
			name:    "количество символов превышает допустимое количество",
			cType:   content.TypeText,
			value:   strings.Repeat("a", 1e5),
			wantErr: variant.ErrInvalid,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			//Arrange
			_, err := makeVariant(makeContent(tt.cType, tt.value))

			//Assert
			assert.Error(t, err)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
