package element_test

import (
	"testing"

	"gitflic.ru/lms/internal/domain/question/content"
	"gitflic.ru/lms/internal/domain/question/matching/pair/element"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	//Arrange
	e, err := element.New(makeContent(content.TypeText, "text"))

	//Assert
	assert.NoError(t, err)
	assert.Equal(t, e.Value(), "text")
}
