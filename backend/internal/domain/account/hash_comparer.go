package account

import (
	"gitflic.ru/lms/backend/internal/domain/account/passhash"
)

// PasswordComparer экспортируемый интерфейс из домена в слой инфраструктуры
// для реализации сравнения хеша и открытого пароля. Алгоритм
// остается в слое инфраструктуры.
type PasswordComparer interface {
	Compare(hash passhash.Hash, plain string) bool
}
