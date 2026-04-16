package resource

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValid(t *testing.T) {
	t.Run("valid resource", func(t *testing.T) {
		resource := Resource(2)

		assert.True(t, resource.IsValid(), "expected valid resource")
	})

	t.Run("invalid resource", func(t *testing.T) {
		resources := []Resource{unknown, count, Resource(100), Resource(-1)}

		for _, r := range resources {
			assert.False(t, r.IsValid(), "expected invalid resource for value: %v", r)
		}
	})
}

