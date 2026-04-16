package action

import "testing"

func TestIsValid(t *testing.T) {
	t.Run("valid action", func(t *testing.T) {
		action := Action(2)

		if !action.IsValid() {
			t.Errorf("expected valid action")
		}
	})

	t.Run("invalid action", func(t *testing.T) {
		actions := []Action{unknown, count, Action(100), Action(-1)}

		for i := range actions {
			if actions[i].IsValid() {
				t.Errorf("expected invalid action")
			}
		}
	})
}
