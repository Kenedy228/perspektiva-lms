package base

import "testing"

func TestValidateText(t *testing.T) {
	type when struct {
		text string
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
				text: "",
			},
			want: want{
				err: ErrEmptyText,
			},
		},
		{
			name: "whitespaces text",
			when: when{
				text: "     ",
			},
			want: want{
				err: ErrEmptyText,
			},
		},
		{
			name: "valid text",
			when: when{
				text: "text",
			},
			want: want{
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateText(tt.when.text)

			if err != tt.want.err {
				t.Errorf("expected err %v, got %v", tt.want.err, err)
			}
		})
	}
}

func TestValidateDescription(t *testing.T) {
	type when struct {
		description string
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
			name: "empty description",
			when: when{
				description: "",
			},
			want: want{
				err: ErrEmptyDescription,
			},
		},
		{
			name: "whitespaces description",
			when: when{
				description: "     ",
			},
			want: want{
				err: ErrEmptyDescription,
			},
		},
		{
			name: "valid description",
			when: when{
				description: "description",
			},
			want: want{
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateDescription(tt.when.description)

			if err != tt.want.err {
				t.Errorf("expected err %v, got %v", tt.want.err, err)
			}
		})
	}
}
