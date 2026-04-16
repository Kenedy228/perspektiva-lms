package action

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValid(t *testing.T) {
	t.Run("valid action", func(t *testing.T) {
		action := Action(2)

		assert.True(t, action.IsValid(), "expected valid action")
	})

	t.Run("invalid action", func(t *testing.T) {
		actions := []Action{unknown, count, Action(100), Action(-1)}

		for _, a := range actions {
			assert.False(t, a.IsValid(), "expected invalid action for value: %v", a)
		}
	})
}

