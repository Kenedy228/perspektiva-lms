package criteria_test

import (
	"fmt"
	"testing"

	"gitflic.ru/lms/internal/domain/quiz/criteria"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewManual(t *testing.T) {
	t.Run("should return err if question ids empty", func(t *testing.T) {
		//Arrange-Assert
		newManualBuilder().build(t, criteria.ErrInvalidQuestions)
	})

	t.Run("should return err if there is nil id", func(t *testing.T) {
		//Arrange-Assert
		newManualBuilder().withQuestionID(uuid.New()).
			withQuestionID(uuid.Nil).
			build(t, criteria.ErrInvalidQuestions)
	})

	t.Run("should return err if there is duplicates of questionIDs", func(t *testing.T) {
		//Arrange-Assert
		duplicate := uuid.New()
		newManualBuilder().withQuestionID(duplicate).
			withQuestionID(duplicate).
			build(t, criteria.ErrInvalidQuestions)
	})

	t.Run("should return err if size of questions > maxSize", func(t *testing.T) {
		//Arrange-Assert
		newManualBuilder().withMaxSizeQuestions().
			build(t, criteria.ErrInvalidQuestions)
	})

	t.Run("valid", func(t *testing.T) {
		//Arrange
		c := castManual(t, newManualBuilder().withQuestionID(uuid.New()).
			withQuestionID(uuid.New()).
			build(t, nil))

		//Assert
		assert.Equal(t, c.Type(), criteria.TypeManual)
		assert.Equal(t, c.QuestionCount(), 2)
		assert.Equal(t, len(c.QuestionIDs()), 2)
	})
}

func BenchmarkValidateQuestionIDs(b *testing.B) {
	sizes := []int{1e2, 1e3, 1e4, 1e5}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("size %d", size), func(b *testing.B) {
			builder := newManualBuilder()

			for range size {
				builder = builder.withQuestionID(uuid.New())
			}

			for b.Loop() {
				builder.buildNoTest()
			}
		})
	}
}
