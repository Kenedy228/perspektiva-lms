package bank

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name  string
		title string
		err   error
	}{
		{
			name:  "empty title",
			title: "",
			err:   ErrEmptyTitle,
		},
		{
			name:  "whitespaces title",
			title: " ",
			err:   ErrEmptyTitle,
		},
		{
			name:  "valid title",
			title: "valid",
			err:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := New(tt.title)

			assert.ErrorIs(t, err, tt.err)

			if b != nil {
				assert.NotEqual(t, uuid.Nil, b.ID(), "expected id non nil")
				assert.Equal(t, tt.title, b.Title())
				assert.Empty(t, b.Questions(), "expected len of questions 0")
				assert.Equal(t, b.CreatedAt(), b.UpdatedAt(), "createdAt should equal to updatedAt")
			}
		})
	}
}

func TestRename(t *testing.T) {
	tests := []struct {
		name string
		old  string
		new  string
		err  error
	}{
		{
			name: "empty new",
			old:  "title",
			new:  "",
			err:  ErrEmptyTitle,
		},
		{
			name: "whitespaces new",
			old:  "title",
			new:  "  ",
			err:  ErrEmptyTitle,
		},
		{
			name: "valid new",
			old:  "title",
			new:  "new title",
			err:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := New(tt.old)
			require.NoError(t, err, "expected err nil on setup")

			oldUpdatedAt := b.UpdatedAt()
			err = b.Rename(tt.new)

			assert.ErrorIs(t, err, tt.err)

			if err == nil {
				assert.NotEqual(t, tt.old, b.Title(), "expected rename title, but got title unchanged")
				assert.Equal(t, tt.new, b.Title())
				assert.False(t, oldUpdatedAt.After(b.UpdatedAt()), "expected updatedAt change")
			}
		})
	}
}

func TestAddQuestionsWithNoErrors(t *testing.T) {
	tests := []struct {
		name      string
		title     string
		questions []uuid.UUID
		total     int
	}{
		{
			name:      "no errors questions",
			title:     "title",
			questions: []uuid.UUID{uuid.New(), uuid.New(), uuid.New()},
			total:     3,
		},
		{
			name:      "no questions",
			title:     "title",
			questions: []uuid.UUID{},
			total:     0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := New(tt.title)
			require.NoError(t, err)

			total, err := b.AddQuestions(tt.questions...)

			assert.NoError(t, err)
			assert.Equal(t, tt.total, total)
			assert.Len(t, b.Questions(), total)

			if len(b.Questions()) != 0 {
				firstQCopy := b.Questions()
				secondQCopy := b.Questions()

				// Проверяем, что слайсы указывают на разные области памяти (инкапсуляция)
				assert.NotSame(t, &firstQCopy[0], &secondQCopy[0], "expected slices with different addresses")
			}
		})
	}
}

func TestAddQuestionsWithError(t *testing.T) {
	bank, err := New("title")
	require.NoError(t, err)

	q := uuid.New()

	total, err := bank.AddQuestions(q, q)

	assert.ErrorIs(t, err, ErrQuestionDuplicate)
	assert.Equal(t, 0, total)
}

func TestRemoveQuestions(t *testing.T) {
	bank, err := New("title")
	require.NoError(t, err)

	toDel := uuid.New()

	_, err = bank.AddQuestions(uuid.New(), uuid.New(), toDel)
	require.NoError(t, err)

	bank.RemoveQuestions(toDel)

	// assert.NotContains проверяет, что элемента нет в слайсе (заменяет проверку через slices.Index)
	assert.NotContains(t, bank.Questions(), toDel)
}

func TestClearQuestions(t *testing.T) {
	bank, err := New("title")
	require.NoError(t, err)

	_, err = bank.AddQuestions(uuid.New(), uuid.New(), uuid.New())
	require.NoError(t, err)

	bank.ClearQuestions()

	// assert.Empty заменяет проверку длины на 0
	assert.Empty(t, bank.Questions(), "expected to clear questions")
}

func BenchmarkAddQuestionsWithNoError(b *testing.B) {
	sizes := []int{10, 1000, 100000}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("size is %d", size), func(b *testing.B) {
			questions := make([]uuid.UUID, size)

			for i := range questions {
				questions[i] = uuid.New()
			}

			bank, _ := New("title")

			for b.Loop() {
				bank.AddQuestions(questions...)
			}
		})
	}
}

func BenchmarkAddQuestionsWithError(b *testing.B) {
	sizes := []int{9, 999, 99999}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("size is %d", size), func(b *testing.B) {
			questions := make([]uuid.UUID, size+1)

			for i := range questions {
				questions[i] = uuid.New()
			}

			questions[len(questions)-1] = questions[len(questions)-2]

			bank, _ := New("title")

			for b.Loop() {
				bank.AddQuestions(questions...)
			}
		})
	}
}

