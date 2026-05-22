package account

import "fmt"

func ensureNotDeleted(a *Account) error {
	if a.Status() == StatusDeleted {
		return fmt.Errorf("%w: над удаленным аккаунтом нельзя совершать какие-либо действия", ErrInvalid)
	}

	return nil
}
