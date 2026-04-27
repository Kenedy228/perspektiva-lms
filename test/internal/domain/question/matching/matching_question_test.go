package matching_test

import (
	"fmt"
	"testing"
	"testing/synctest"
	"time"

	"gitflic.ru/lms/internal/domain/question"
	"gitflic.ru/lms/internal/domain/question/matching"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

const maxPairs = 200

func TestNew(t *testing.T) {
	t.Run("ErrEmptyPairs", func(t *testing.T) {
		//Arrange-Assert
		newQuestionBuilder().build(t, matching.ErrInvalidPairs)
	})

	t.Run("ErrNotEnoughPairs", func(t *testing.T) {
		//Arrange-Assert
		newQuestionBuilder().withPairParamAsText("пара", "значение").
			build(t, matching.ErrInvalidPairs)
	})

	t.Run("ErrTooManyPairs", func(t *testing.T) {
		//Arrange
		b := newQuestionBuilder()

		for i := range maxPairs + 1 {
			v := fmt.Sprintf("%d", i)
			b = b.withPairParamAsText(v, v)
		}

		//Assert
		b.build(t, matching.ErrInvalidPairs)
	})

	t.Run("ErrDuplicatePrompt", func(t *testing.T) {
		//Arrange-Assert
		newQuestionBuilder().withPairParamAsText("дубликат", "значение 1").
			withPairParamAsText("дубликат", "значение 2").
			build(t, matching.ErrInvalidPairs)
	})

	t.Run("ErrDuplicateOption", func(t *testing.T) {
		//Arrange-Assert
		newQuestionBuilder().withPairParamAsText("пара 1", "значение").
			withPairParamAsText("пара 2", "значение").
			build(t, matching.ErrInvalidPairs)
	})

	t.Run("Valid", func(t *testing.T) {
		t.Run("plain creation", func(t *testing.T) {
			//Arrange
			imageID := uuid.New()
			q := castQuestion(t, newQuestionBuilder().withText("text").
				withImageID(imageID).
				withPairParam("пара изображение", "path/to-image.png", question.ContentTypeImage).
				withPairParam("пара аудио", "path/to-audio.wav", question.ContentTypeAudio).
				build(t, nil))

			//Assert
			assert.Equal(t, q.Text(), "text")
			assert.Equal(t, q.Instruction(), question.TypeMatching.DefaultInstruction())
			assert.Equal(t, q.ImageID(), imageID)
			assert.True(t, q.HasImage())
			assert.Equal(t, len(q.Pairs()), 2)
		})

		t.Run("without image", func(t *testing.T) {
			//Arrange
			q := castQuestion(t, newQuestionBuilder().withText("text").
				withPairParam("пара изображение", "path/to-image.png", question.ContentTypeImage).
				withPairParamAsText("пара текст", "text").
				build(t, nil))

			//Assert
			assert.Equal(t, q.Text(), "text")
			assert.Equal(t, q.Instruction(), question.TypeMatching.DefaultInstruction())
			assert.Equal(t, q.ImageID(), uuid.Nil)
			assert.False(t, q.HasImage())
			assert.Equal(t, len(q.Pairs()), 2)
		})
	})
}

func TestMatchingQuestion_UpdatePairs(t *testing.T) {
	t.Run("success update should change updated at", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			//Arrange
			q := castQuestion(t, newQuestionBuilder().withText("Сопоставьте страны и столицы").
				withPairParamAsText("Россия", "1").
				withPairParamAsText("Япония", "2").
				build(t, nil),
			)
			pairs := mockPairs()

			//Act
			time.Sleep(time.Second)
			err := q.UpdatePairs(pairs)
			assert.Nil(t, err)

			//Assert
			assert.Equal(t, q.Pairs(), pairs)
			assert.True(t, q.UpdatedAt().After(q.CreatedAt()))
		})
	})

	t.Run("error should not change updated at pairs", func(t *testing.T) {
		synctest.Test(t, func(t *testing.T) {
			//Arrange
			q := castQuestion(t, newQuestionBuilder().withText("Сопоставьте страны и столицы").
				withPairParamAsText("Россия", "1").
				withPairParamAsText("Япония", "2").
				build(t, nil),
			)

			//Act
			time.Sleep(time.Second)
			err := q.UpdatePairs(nil)
			assert.ErrorIs(t, err, matching.ErrInvalidPairs)

			//Assert
			assert.Equal(t, q.Pairs(),
				[]matching.Pair{
					newPairBuilder().withPrompt("Россия").withContentText("1").buildNoTest(),
					newPairBuilder().withPrompt("Япония").withContentText("2").buildNoTest(),
				},
			)
			assert.Equal(t, q.UpdatedAt(), q.CreatedAt())
		})
	})
}

func TestClone(t *testing.T) {
	//Arrange
	q := castQuestion(t, newQuestionBuilder().withPairParamAsText("prompt", "val").
		withPairParamAsText("prompt2", "val2").build(t, nil),
	)
	clone := castQuestion(t, q.Clone())

	//Assert
	assert.Equal(t, clone.ID(), q.ID())
	assert.Equal(t, clone.Text(), q.Text())
	assert.Equal(t, clone.Instruction(), q.Instruction())
	assert.Equal(t, clone.ImageID(), q.ImageID())
	assert.Equal(t, clone.Pairs(), q.Pairs())
	assert.Equal(t, clone.CreatedAt(), q.CreatedAt())
	assert.Equal(t, clone.UpdatedAt(), q.UpdatedAt())
	assert.NotSame(t, &clone.Pairs()[0], &q.Pairs()[0])
}

func TestCheckAnswer(t *testing.T) {
	q := castQuestion(t, newQuestionBuilder().withPairParamAsText("prompt", "val").
		withPairParamAsText("prompt2", "val2").
		withPairParamAsText("prompt3", "val3").
		withPairParamAsText("prompt4", "val4").
		build(t, nil),
	)

	t.Run("answer pairs length equal to q pairs length, but prompts different", func(t *testing.T) {
		//Arrange
		answer := newAnswerBuilder().withTextPair("prompt", "val").
			withTextPair("prompt2", "val2").
			withTextPair("prompt3", "val3").
			withTextPair("prompt5", "val4").
			build(t)

		//Act
		check := q.CheckAnswer(answer)

		//Assert
		assert.False(t, check)
	})

	t.Run("answer pairs length equal to q pairs length, but values different", func(t *testing.T) {
		//Arrange
		answer := newAnswerBuilder().withTextPair("prompt", "val").
			withTextPair("prompt2", "val2").
			withTextPair("prompt3", "val3").
			withTextPair("prompt4", "val5").
			build(t)

		//Act
		check := q.CheckAnswer(answer)

		//Assert
		assert.False(t, check)
	})

	t.Run("answer pairs length equal to q pairs length, but contains all duplicates", func(t *testing.T) {
		//Arrange
		answer := newAnswerBuilder().withTextPair("prompt", "val").
			withTextPair("prompt", "val").
			withTextPair("prompt", "val").
			withTextPair("prompt", "val").
			build(t)

		//Act
		check := q.CheckAnswer(answer)

		//Assert
		assert.False(t, check)
	})

	t.Run("answer pairs length less than q pairs length, but contains present answers", func(t *testing.T) {
		//Arrange
		answer := newAnswerBuilder().withTextPair("prompt", "val").
			withTextPair("prompt2", "val2").
			withTextPair("prompt3", "val3").
			build(t)

		//Act
		check := q.CheckAnswer(answer)

		//Assert
		assert.False(t, check)
	})

	t.Run("answer pairs length less than q pairs length, but contains one invalid answer", func(t *testing.T) {
		//Arrange
		answer := newAnswerBuilder().withTextPair("prompt", "val").
			withTextPair("prompt2", "val2").
			withTextPair("5rompt", "val3").
			build(t)

		//Act
		check := q.CheckAnswer(answer)

		//Assert
		assert.False(t, check)
	})

	t.Run("answer pairs length greater than q pairs length, but contains all duplicates valid answers", func(t *testing.T) {
		//Arrange
		answer := newAnswerBuilder().withTextPair("prompt", "val").
			withTextPair("prompt", "val").
			withTextPair("prompt", "val").
			withTextPair("prompt", "val").
			withTextPair("prompt", "val").
			withTextPair("prompt", "val").
			build(t)

		//Act
		check := q.CheckAnswer(answer)

		//Assert
		assert.False(t, check)
	})

	t.Run("answer pairs length greater than q pairs length, but contains invalid answers", func(t *testing.T) {
		//Arrange
		answer := newAnswerBuilder().withTextPair("prompt", "val").
			withTextPair("prompt1", "val").
			withTextPair("prompt2", "val").
			withTextPair("prompt3", "val").
			withTextPair("prompt4", "val").
			withTextPair("prompt5", "val").
			build(t)

		//Act
		check := q.CheckAnswer(answer)

		//Assert
		assert.False(t, check)
	})

	t.Run("answers are valid", func(t *testing.T) {
		//Arrange
		answer := newAnswerBuilder().withTextPair("prompt", "val").
			withTextPair("prompt2", "val2").
			withTextPair("prompt3", "val3").
			withTextPair("prompt4", "val4").
			build(t)

		//Act
		check := q.CheckAnswer(answer)

		//Assert
		assert.True(t, check)
	})
}
