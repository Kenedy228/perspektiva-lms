package title_test

import (
	"strings"
	"testing"

	"gitflic.ru/lms/internal/domain/quiz/title"
	"github.com/stretchr/testify/assert"
)

func TestNew_Success(t *testing.T) {
	tc := []struct {
		name  string
		value string
	}{
		{
			name:  "валидный кейс",
			value: "Заголовок для промежуточного тестирования",
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			//Arrange
			got, err := title.New(tt.value)

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
		value   string
		wantErr error
	}{
		{
			name:    "пустой заголовк",
			value:   "",
			wantErr: title.ErrInvalid,
		},
		{
			name:    "количество символов выходит за пределы допустимого",
			value:   strings.Repeat("a", 1e5),
			wantErr: title.ErrInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//Arrange
			_, err := title.New(tt.value)

			//Assert
			assert.Error(t, err)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
