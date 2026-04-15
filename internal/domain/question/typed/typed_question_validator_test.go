package typed

import (
	"errors"
	"testing"
)

func TextValidatePlaceholders(t *testing.T) {
	type when struct {
		text              string
		blanks            map[string][]string
		placeholdersCount int
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
			name: "placeholder count different from blanks len",
			when: when{
				text:              "text",
				blanks:            map[string][]string{"blank1": {"hi"}},
				placeholdersCount: 5,
			},
			want: want{
				err: ErrPlaceholderCountMismatch,
			},
		},
		{
			name: "placeholder missing",
			when: when{
				text:              "text [blank1]",
				blanks:            map[string][]string{"[blank2]": {"answer"}},
				placeholdersCount: 1,
			},
			want: want{
				err: ErrPlaceholderMissing,
			},
		},
		{
			name: "blank duplicate",
			when: when{
				text:              "text [blank1] [blank1]",
				blanks:            map[string][]string{"[blank1]": {"answer"}},
				placeholdersCount: 1,
			},
			want: want{
				err: ErrMarkDuplicate,
			},
		},
		{
			name: "valid blanks",
			when: when{
				text:              "text [blank1] [blank2]",
				blanks:            map[string][]string{"[blank1]": {"answer"}, "[blank2]": {"answer", "answer"}},
				placeholdersCount: 2,
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "too many placeholders",
			when: when{
				text:              "text [blank1] [blank2]",
				blanks:            map[string][]string{"[blank1]": {"answer"}, "[blank2]": {"answer", "answer"}},
				placeholdersCount: 21,
			},
			want: want{
				err: ErrTooManyPlaceholders,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validatePlaceholders(tt.when.text, tt.when.placeholdersCount, tt.when.blanks)

			if !errors.Is(err, tt.want.err) {
				t.Errorf("expected err %v, got %v", tt.want.err, err)
			}
		})
	}

}

func TestValidateBlank(t *testing.T) {
	type when struct {
		mark    string
		answers []string
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
			name: "empty mark",
			when: when{
				mark:    "",
				answers: []string{"answer"},
			},
			want: want{
				err: ErrEmptyMark,
			},
		},
		{
			name: "whitespaces mark",
			when: when{
				mark:    "   ",
				answers: []string{"answer"},
			},
			want: want{
				err: ErrEmptyMark,
			},
		},
		{
			name: "empty answers",
			when: when{
				mark:    "mark",
				answers: []string{},
			},
			want: want{
				err: ErrNoBlankAnswers,
			},
		},
		{
			name: "empty answer",
			when: when{
				mark:    "mark",
				answers: []string{""},
			},
			want: want{
				err: ErrEmptyBlankAnswer,
			},
		},
		{
			name: "whitespaces answer",
			when: when{
				mark:    "mark",
				answers: []string{"   "},
			},
			want: want{
				err: ErrEmptyBlankAnswer,
			},
		},
		{
			name: "valid blank",
			when: when{
				mark:    "mark",
				answers: []string{"answer", "answer"},
			},
			want: want{
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateBlank(tt.when.mark, tt.when.answers)

			if err != tt.want.err {
				t.Errorf("expected err %v, got %v", tt.want.err, err)
			}
		})
	}
}
