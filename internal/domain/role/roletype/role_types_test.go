package roletype

import "testing"

func TestIsValid(t *testing.T) {
	t.Run("valid type", func(t *testing.T) {
		rt := RoleType(2)

		if !rt.IsValid() {
			t.Errorf("expected valid action")
		}
	})

	t.Run("invalid type", func(t *testing.T) {
		rts := []RoleType{unknown, count, RoleType(100), RoleType(-1)}

		for i := range rts {
			if rts[i].IsValid() {
				t.Errorf("expected invalid action")
			}
		}
	})

}
