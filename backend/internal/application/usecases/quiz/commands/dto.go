package commands

import "gitflic.ru/lms/backend/internal/domain/role"

type SourceInput struct {
	BankID        string
	CriteriaType  string
	QuestionCount int
	QuestionIDs   []string
}

type CreateInput struct {
	ActorRole        role.Role
	Title            string
	MaxAttempts      int
	TimeLimitSeconds int
	ShuffleQuestions bool
	Sources          []SourceInput
}

type QuizIDInput struct {
	ActorRole role.Role
	QuizID    string
}

type RenameInput struct {
	ActorRole role.Role
	QuizID    string
	Title     string
}

type ChangeLimitsInput struct {
	ActorRole        role.Role
	QuizID           string
	MaxAttempts      int
	TimeLimitSeconds int
}

type ChangeShufflePolicyInput struct {
	ActorRole        role.Role
	QuizID           string
	ShuffleQuestions bool
}

type ChangeSourceInput struct {
	ActorRole role.Role
	QuizID    string
	Source    SourceInput
}

type ReplaceSourcesInput struct {
	ActorRole role.Role
	QuizID    string
	Sources   []SourceInput
}

type Output struct {
	ID string
}
