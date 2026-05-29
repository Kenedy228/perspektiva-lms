package enrollment

import (
	"context"

	"github.com/google/uuid"
)

type EnrollmentView struct {
	ID            string `json:"id"`
	AccountID     string `json:"account_id"`
	CourseID      string `json:"course_id"`
	ActivatedAt   string `json:"activated_at"`
	DeactivatedAt string `json:"deactivated_at"`
	Status        string `json:"status"`
	StatusTitle   string `json:"status_title"`
}

type ListFilter struct {
	AccountID uuid.UUID
	CourseID  uuid.UUID
	Limit     int
	Offset    int
}

type QueryService interface {
	GetByID(ctx context.Context, id uuid.UUID) (EnrollmentView, error)
	List(ctx context.Context, filter ListFilter) ([]EnrollmentView, error)
}
