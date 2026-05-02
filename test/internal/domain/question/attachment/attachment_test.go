package attachment_test

import (
	"testing"

	"gitflic.ru/lms/internal/domain/question/attachment"
	"gitflic.ru/lms/internal/domain/question/content"
	"github.com/stretchr/testify/assert"
)

func TestNew_Success(t *testing.T) {
	tc := []struct {
		name  string
		cType content.Type
		value string
	}{
		{
			name:  "валидный контент типа изображение",
			cType: content.TypeImage,
			value: "image.png",
		},
		{
			name:  "валидный контент типа аудио",
			cType: content.TypeAudio,
			value: "audio.mp3",
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			//Arrange
			got, err := attachment.New(makeContent(tt.cType, tt.value))

			//Assert
			assert.NoError(t, err)
			assert.Equal(t, got.Content().Type(), tt.cType)
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
			name:    "текстовый контент недопустим в качестве вложения",
			cType:   content.TypeText,
			value:   "текст вложения",
			wantErr: attachment.ErrInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//Arrange
			_, err := attachment.New(makeContent(tt.cType, tt.value))

			//Assert
			assert.Error(t, err)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
