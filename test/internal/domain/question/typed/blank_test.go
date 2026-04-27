package typed_test

import (
	"fmt"
	"testing"

	"gitflic.ru/lms/internal/domain/question/typed"
	"github.com/stretchr/testify/assert"
)

const maxVariants = 100

func TestNewBlank(t *testing.T) {
	t.Run("invalid placeholder", func(t *testing.T) {
		//Arrange-Assert
		newBlankBuilder().withInvalidPlaceholder("invalid").
			withVariant("answer").
			build(t, typed.ErrInvalidPlaceholder)
	})

	t.Run("len answers less than minAnswers", func(t *testing.T) {
		//Arrange-Assert
		newBlankBuilder().withPlaceholder("valid").
			build(t, typed.ErrInvalidVariants)
	})

	t.Run("len answers greater than maxAnswers", func(t *testing.T) {
		//Arrange-Assert
		b := newBlankBuilder().withPlaceholder("valid")

		for i := range maxVariants {
			b = b.withVariant(fmt.Sprintf("%d", i))
		}

		b.build(t, typed.ErrInvalidVariants)
	})

	t.Run("invalid answer format", func(t *testing.T) {
		//Arrange-Assert
		newBlankBuilder().withInvalidVariant().
			withPlaceholder("valid").
			build(t, typed.ErrInvalidVariants)
	})

	t.Run("valid", func(t *testing.T) {
		//Arrange
		blank := newBlankBuilder().withPlaceholder("valid").
			withVariant("answer").
			withVariant("answer2").
			build(t, nil)

		//Assert
		assert.Equal(t, blank.Placeholder(), "{{valid}}")
		assert.Equal(t, len(blank.Variants()), 2)
		assert.True(t, blank.Variants()[0].Equal(makeContent("answer")))
		assert.True(t, blank.Variants()[1].Equal(makeContent("answer2")))
	})
}
