package commands

import (
	"time"

	"gitflic.ru/lms/backend/internal/domain/course/progress"
	"gitflic.ru/lms/backend/internal/domain/role"
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

type MoveBlockElementInput struct {
	ActorRole role.Role
	BlockID   string
	From      int
	To        int
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

type Output struct {
	ID string
}
