package account

import (
	"gitflic.ru/lms/internal/domain/account/login"
	"gitflic.ru/lms/internal/domain/account/passhash"
	"gitflic.ru/lms/internal/domain/role"
	"github.com/google/uuid"
)

type Params struct {
	Login        login.Login
	PasswordHash passhash.Hash
	Role         role.Role
	PersonID     uuid.UUID
}
