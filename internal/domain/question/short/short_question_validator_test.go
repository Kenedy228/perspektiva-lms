package short

import (
	"errors"
	"fmt"
	"testing"
)

func TestValidateAnswers(t *testing.T) {
	type when struct {
		answers         []string
		allowDuplicates bool
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
			name: "empty answers",
			when: when{
				answers:         []string{},
				allowDuplicates: false,
			},
			want: want{
				err: ErrNoAnswers,
			},
		},
		{
			name: "contains empty answer",
			when: when{
				answers:         []string{"answer", "", "answer2"},
				allowDuplicates: false,
			},
			want: want{
				err: ErrEmptyAnswer,
			},
		},
		{
			name: "contains whitespaces answer",
			when: when{
				answers:         []string{"answer", "    ", "answer2"},
				allowDuplicates: false,
			},
			want: want{
				err: ErrEmptyAnswer,
			},
		},
		{
			name: "contains identical answers with allowDuplicates false",
			when: when{
				answers:         []string{"answer", "answer"},
				allowDuplicates: false,
			},
			want: want{
				err: ErrDuplicateAnswer,
			},
		},
		{
			name: "contains identical answers with allowDuplicates true",
			when: when{
				answers:         []string{"answer", "answer"},
				allowDuplicates: true,
			},
			want: want{
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateAnswers(tt.when.answers, tt.when.allowDuplicates)

			if !errors.Is(err, tt.want.err) {
				t.Errorf("expected err %v, got %v", tt.want.err, err)
			}
		})
	}
}

func TestValidateAnswersWithLimits(t *testing.T) {
	type when struct {
		size int
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
			name: "size 20",
			when: when{
				size: 20,
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "size 21",
			when: when{
				size: 21,
			},
			want: want{
				err: ErrTooManyAnswers,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			answers := make([]string, tt.when.size)

			for i := range tt.when.size {
				answers[i] = fmt.Sprintf("answer%d", i)
			}

			err := validateAnswers(answers, false)

			if !errors.Is(err, tt.want.err) {
				t.Errorf("expected err %v, got %v", tt.want.err, err)
			}
		})
	}
}
