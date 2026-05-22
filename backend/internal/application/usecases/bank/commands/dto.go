package commands

import "gitflic.ru/lms/backend/internal/domain/role"

type CreateInput struct {
	ActorRole role.Role
	Title     string
}

type RenameInput struct {
	ActorRole role.Role
	BankID    string
	Title     string
}

type QuestionIDsInput struct {
	ActorRole   role.Role
	BankID      string
	QuestionIDs []string
}

type BankIDInput struct {
	ActorRole role.Role
	BankID    string
}

type Output struct {
	ID string
}
