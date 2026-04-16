package resource

import "testing"

func TestIsValid(t *testing.T) {
	t.Run("valid resource", func(t *testing.T) {
		resource := Resource(2)

		if !resource.IsValid() {
			t.Errorf("expected valid resource")
		}
	})

	t.Run("invalid resouce", func(t *testing.T) {
		resources := []Resource{unknown, count, Resource(100), Resource(-1)}

		for i := range resources {
			if resources[i].IsValid() {
				t.Errorf("expected invalid resource")
			}
		}
	})
}
