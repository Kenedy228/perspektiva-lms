package account

import "gitflic.ru/lms/internal/domain/account/passhash"

type PasswordComparer interface {
	Compare(hash passhash.Hash, plain string) bool
}
