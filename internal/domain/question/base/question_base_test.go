package base

import (
	"testing"

	"gitflic.ru/lms/internal/domain/question"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name        string
		text        question.QText
		description question.QDescription
		imageID     uuid.UUID
	}{
		{
			name:        "nil image",
			text:        createText("text"),
			description: createDescription("description"),
			imageID:     uuid.Nil,
		},
		{
			name:        "non-nil image",
			text:        createText("text"),
			description: createDescription("description"),
			imageID:     uuid.New(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := createParams(tt.text, tt.description, tt.imageID)
			baseObj, err := New(params)

			require.Nil(t, err)
			assert.NotEqual(t, uuid.Nil, baseObj.ID())
			assert.Equal(t, tt.text, baseObj.Text())
			assert.Equal(t, tt.description, baseObj.Description())

			if tt.imageID != uuid.Nil {
				assert.True(t, baseObj.HasImage())
				assert.Equal(t, tt.imageID, baseObj.ImageID())
			} else {
				assert.False(t, baseObj.HasImage())
				assert.Equal(t, tt.imageID, uuid.Nil)
			}
		})
	}
}

func TestUpdateText(t *testing.T) {
	tests := []struct {
		name         string
		old          question.QText
		new          question.QText
		shouldUpdate bool
	}{
		{
			name:         "different text",
			old:          createText("old"),
			new:          createText("new"),
			shouldUpdate: true,
		},
		{
			name:         "same text",
			old:          createText("old"),
			new:          createText("old"),
			shouldUpdate: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := createParams(tt.old, createDescription("description"), uuid.Nil)

			baseObj, err := New(params)
			require.NoError(t, err, "expected no errors on setup")

			oldUpdatedAt := baseObj.UpdatedAt()
			baseObj.UpdateText(tt.new)

			assert.Equal(t, tt.new, baseObj.Text())

			if tt.shouldUpdate {
				assert.True(t, oldUpdatedAt.Before(baseObj.UpdatedAt()))
			} else {
				assert.Equal(t, oldUpdatedAt, baseObj.UpdatedAt())
			}
		})
	}
}

func TestUpdateImage(t *testing.T) {
	tests := []struct {
		name         string
		old          uuid.UUID
		new          uuid.UUID
		hasImage     bool
		shouldUpdate bool
	}{
		{
			name:         "same image nil",
			old:          uuid.Nil,
			new:          uuid.Nil,
			hasImage:     false,
			shouldUpdate: false,
		},
		{
			name:         "different image not nil",
			old:          uuid.Nil,
			new:          uuid.New(),
			hasImage:     true,
			shouldUpdate: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := createParams(createText("text"), createDescription("description"), tt.old)

			baseObj, err := New(params)
			require.NoError(t, err, "expected no errors on setup")

			oldUpdatedAt := baseObj.UpdatedAt()
			baseObj.UpdateImage(tt.new)

			assert.Equal(t, tt.new, baseObj.ImageID())

			if tt.shouldUpdate {
				assert.True(t, oldUpdatedAt.Before(baseObj.UpdatedAt()))
			} else {
				assert.Equal(t, oldUpdatedAt, baseObj.UpdatedAt())
			}
		})
	}
}

func TestRemoveImage(t *testing.T) {
	tests := []struct {
		name         string
		imageID      uuid.UUID
		shouldRemove bool
	}{
		{
			name:         "image is set",
			imageID:      uuid.New(),
			shouldRemove: true,
		},
		{
			name:         "image is not set",
			imageID:      uuid.Nil,
			shouldRemove: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := createParams(createText("text"), createDescription("description"), tt.imageID)

			baseObj, err := New(params)
			require.NoError(t, err, "expected no errors on setup")

			oldUpdatedAt := baseObj.UpdatedAt()
			baseObj.RemoveImage()

			if tt.shouldRemove {
				assert.Equal(t, uuid.Nil, baseObj.ImageID())
				assert.True(t, oldUpdatedAt.Before(baseObj.UpdatedAt()))
			} else {
				assert.Equal(t, oldUpdatedAt, baseObj.UpdatedAt())
			}
		})
	}
}

func createParams(text question.QText, desc question.QDescription, imageID uuid.UUID) Params {
	return Params{
		Text:        text,
		Description: desc,
		ImageID:     imageID,
	}
}

func createText(s string) question.QText {
	text, err := question.NewQText(s)
	if err != nil {
		panic(err)
	}

	return text
}

func createDescription(s string) question.QDescription {
	desc, err := question.NewQDescription(s)
	if err != nil {
		panic(err)
	}

	return desc
}
