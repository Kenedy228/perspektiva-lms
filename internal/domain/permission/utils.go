package permission

func reviewActions(actions []Action) error {
	for i := range actions {
		if !actions[i].IsValid() {
			return ErrInvalidAction
		}

		for j := range i {
			if actions[i] == actions[j] {
				return ErrDuplicatedActions
			}
		}
	}
	return nil
}
