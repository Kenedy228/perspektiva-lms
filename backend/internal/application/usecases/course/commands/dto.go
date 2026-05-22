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

type CreateVersionInput struct {
	ActorRole role.Role
	CourseID  string
	Title     string
}

type VersionIDInput struct {
	ActorRole role.Role
	VersionID string
}

type AddBlockInput struct {
	ActorRole role.Role
	VersionID string
	Title     string
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
