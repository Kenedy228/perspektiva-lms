package attempt

import (
	"gitflic.ru/lms/internal/domain/question"
	"github.com/google/uuid"
)

type Params struct {
	EnrollmentID uuid.UUID
	CourseID     uuid.UUID
	QuizID       uuid.UUID
	CourseTitle  string
	QuizTitle    string
	QMeta        []QMetaParams
}

type QMetaParams struct {
	QType    question.Type
	Snapshot []byte
}
