package selectable

import (
	"errors"
	"fmt"
	"testing"
)

func TestValidateOptionText(t *testing.T) {
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
				err: ErrEmptyOptionText,
			},
		},
		{
			name: "whitespaces text",
			when: when{
				text: "   ",
			},
			want: want{
				err: ErrEmptyOptionText,
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
			if err := validateOptionText(tt.when.text); !errors.Is(err, tt.want.err) {
				t.Fatalf("expected err %v, got %v", tt.want.err, err)
			}
		})
	}
}

func TestValidateOptions(t *testing.T) {
	type when struct {
		options map[string]bool
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
			name: "no options",
			when: when{
				options: map[string]bool{},
			},
			want: want{
				err: ErrEmptyOptions,
			},
		},
		{
			name: "len less than 2",
			when: when{
				options: map[string]bool{"blank": false},
			},
			want: want{
				err: ErrNotEnoughOptions,
			},
		},
		{
			name: "no correct option",
			when: when{
				options: map[string]bool{"blank1": false, "blank2": false},
			},
			want: want{
				err: ErrNoCorrectOption,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateOptions(tt.when.options); !errors.Is(err, tt.want.err) {
				t.Fatalf("expected err %v, got %v", tt.want.err, err)
			}
		})
	}
}

func TestValidateOptionsWithLimitExceeded(t *testing.T) {
	tests := []struct {
		name string
		size int
		err  error
	}{
		{
			name: "19 size",
			size: 19,
			err:  nil,
		},
		{
			name: "20 size",
			size: 20,
			err:  nil,
		},
		{
			name: "21 size",
			size: 21,
			err:  ErrTooManyOptions,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			options := make(map[string]bool, tt.size)
			options["correct"] = true

			for i := range tt.size - 1 {
				options[fmt.Sprintf("%d", i)] = false
			}

			err := validateOptions(options)

			if !errors.Is(err, tt.err) {
				t.Errorf("expected err %v, got %v", tt.err, err)
			}
		})
	}
}
