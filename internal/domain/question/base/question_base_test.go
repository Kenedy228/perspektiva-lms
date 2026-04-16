package base

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	type when struct {
		text        string
		description string
		image       uuid.UUID
	}

	type want struct {
		err error
	}

	tests := []struct {
		name string
		when
		want
	}{
		{
			name: "empty text",
			when: when{
				text:        "",
				description: "desc",
				image:       uuid.Nil,
			},
			want: want{
				err: ErrEmptyText,
			},
		},
		{
			name: "whitespaces text",
			when: when{
				text:        "    ",
				description: "desc",
				image:       uuid.Nil,
			},
			want: want{
				err: ErrEmptyText,
			},
		},
		{
			name: "empty description",
			when: when{
				text:        "text",
				description: "",
				image:       uuid.Nil,
			},
			want: want{
				err: ErrEmptyDescription,
			},
		},
		{
			name: "whitespaces description",
			when: when{
				text:        "text",
				description: "    ",
				image:       uuid.Nil,
			},
			want: want{
				err: ErrEmptyDescription,
			},
		},
		{
			name: "nil image",
			when: when{
				text:        "text",
				description: "description",
				image:       uuid.Nil,
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "non-nil image",
			when: when{
				text:        "text",
				description: "description",
				image:       uuid.New(),
			},
			want: want{
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := &Params{
				Text:        tt.when.text,
				Description: tt.when.description,
				Image:       tt.when.image,
			}

			baseObj, err := New(params)

			assert.ErrorIs(t, err, tt.want.err)

			// Если ошибки не ожидалось, дополнительно проверяем, что объект создался корректно
			if tt.want.err == nil {
				require.NotNil(t, baseObj)
				assert.NotEqual(t, uuid.Nil, baseObj.ID())
				assert.Equal(t, tt.when.text, baseObj.Text())
				assert.Equal(t, tt.when.description, baseObj.Description())

				if tt.when.image != uuid.Nil {
					assert.True(t, baseObj.HasImage())
					assert.Equal(t, tt.when.image, baseObj.Image())
				} else {
					assert.False(t, baseObj.HasImage())
				}
			}
		})
	}
}

func TestUpdateText(t *testing.T) {
	type given struct {
		text        string
		description string
		image       uuid.UUID
	}

	type when struct {
		text string
	}

	type want struct {
		text string
		err  error
	}

	tests := []struct {
		name string
		given
		when
		want
	}{
		{
			name: "empty text",
			given: given{
				text:        "old text",
				description: "desc",
				image:       uuid.Nil,
			},
			when: when{
				text: "",
			},
			want: want{
				text: "old text",
				err:  ErrEmptyText,
			},
		},
		{
			name: "whitespaces text",
			given: given{
				text:        "old text",
				description: "desc",
				image:       uuid.Nil,
			},
			when: when{
				text: "    ",
			},
			want: want{
				text: "old text",
				err:  ErrEmptyText,
			},
		},
		{
			name: "valid text",
			given: given{
				text:        "old text",
				description: "desc",
				image:       uuid.Nil,
			},
			when: when{
				text: "new text",
			},
			want: want{
				text: "new text",
				err:  nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := &Params{
				Text:        tt.given.text,
				Description: tt.given.description,
				Image:       tt.given.image,
			}

			baseObj, err := New(params)
			require.NoError(t, err, "expected no errors on setup")

			err = baseObj.UpdateText(tt.when.text)

			assert.ErrorIs(t, err, tt.want.err)
			assert.Equal(t, tt.want.text, baseObj.Text())
		})
	}
}

func TestUpdateImage(t *testing.T) {
	type given struct {
		text        string
		description string
		image       uuid.UUID
	}

	type when struct {
		image uuid.UUID
	}

	type want struct {
		hasImage bool
	}

	tests := []struct {
		name string
		given
		when
		want
	}{
		{
			name: "nil image",
			given: given{
				text:        "old text",
				description: "desc",
				image:       uuid.Nil,
			},
			when: when{
				image: uuid.Nil,
			},
			want: want{
				hasImage: false,
			},
		},
		{
			name: "non-nil image",
			given: given{
				text:        "old text",
				description: "desc",
				image:       uuid.Nil,
			},
			when: when{
				image: uuid.New(),
			},
			want: want{
				hasImage: true,
			},
		},
		{
			name: "non-nil image with existing image",
			given: given{
				text:        "old text",
				description: "desc",
				image:       uuid.New(),
			},
			when: when{
				image: uuid.New(),
			},
			want: want{
				hasImage: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := &Params{
				Text:        tt.given.text,
				Description: tt.given.description,
				Image:       tt.given.image,
			}

			baseObj, err := New(params)
			require.NoError(t, err, "expected no errors on setup")

			baseObj.UpdateImage(tt.when.image)

			assert.Equal(t, tt.want.hasImage, baseObj.HasImage())

			// Дополнительно проверяем, что сама картинка установилась (или сбросилась)
			if tt.want.hasImage {
				assert.Equal(t, tt.when.image, baseObj.Image())
			}
		})
	}
}
