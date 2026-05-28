package commands

import (
	"time"

	"gitflic.ru/lms/backend/internal/domain/course/progress"
	"gitflic.ru/lms/backend/internal/domain/role"
	"github.com/google/uuid"
)

type CreateCourseInput struct {
	ActorRole role.Role
	Title     string
}

type CourseIDInput struct {
	ActorRole role.Role
	CourseID  string
}

type RenameCourseInput struct {
	ActorRole role.Role
	CourseID  string
	Title     string
}

type AddBlockToCourseInput struct {
	ActorRole role.Role
	CourseID  string
	Title     string
}

type RemoveBlockFromCourseInput struct {
	ActorRole role.Role
	CourseID  string
	BlockID   string
}

type MoveCourseBlockInput struct {
	ActorRole role.Role
	CourseID  string
	From      int
	To        int
}

type AddElementToBlockInput struct {
	ActorRole role.Role
	BlockID   string
	Title     string
	Content   ElementContentInput
}

type RemoveElementFromBlockInput struct {
	ActorRole role.Role
	BlockID   string
	ElementID string
}

type MoveBlockElementInput struct {
	ActorRole role.Role
	BlockID   string
	From      int
	To        int
}

type ChangeElementCompletionModeInput struct {
	ActorRole      role.Role
	ElementID      string
	CompletionMode string
}

type ElementContentInput struct {
	Type           string
	FileName       string
	SizeBytes      int64
	QuizID         string
	CompletionMode string
}

type MarkProgressInput struct {
	ActorRole    role.Role
	EnrollmentID string
	ElementID    string
	MarkerType   progress.MarkerType
	At           time.Time
}

type UnmarkElementCompletedInput struct {
	ActorRole    role.Role
	EnrollmentID string
	ElementID    string
}

type GetProgressInput struct {
	ActorRole         role.Role
	EnrollmentID      string
	TotalTrackedItems int
}

type ProgressOutput struct {
	CompletedCount      int
	Percent             int
	TotalTrackedItems   int
	CompletedElementIDs []uuid.UUID
}

type Output struct {
	ID string
}
