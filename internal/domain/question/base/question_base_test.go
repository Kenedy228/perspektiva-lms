package base

import (
	"testing"

	"github.com/google/uuid"
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

			_, err := New(params)

			if err != tt.want.err {
				t.Errorf("expected err %v, got %v", tt.want.err, err)
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

			base, err := New(params)

			if err != nil {
				t.Errorf("expected no errors, got %v", err)
			}

			err = base.UpdateText(tt.when.text)

			if err != tt.want.err {
				t.Errorf("expected err %v, got %v", tt.want.err, err)
			}

			if base.Text() != tt.want.text {
				t.Errorf("expected text %v, got %v", tt.want.text, base.Text())
			}
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

			base, err := New(params)

			if err != nil {
				t.Errorf("expected no errors, got %v", err)
			}

			base.UpdateImage(tt.when.image)

			if base.HasImage() != tt.want.hasImage {
				t.Errorf("expected hasImage %v, got %v", tt.want.hasImage, base.HasImage())
			}
		})
	}
}
