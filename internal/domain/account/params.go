package account

import (
	"gitflic.ru/lms/internal/domain/role"
	"github.com/google/uuid"
)

type Params struct {
	Login        string
	PasswordHash string
	Role         role.Role
	PersonID     uuid.UUID
}
