package enrollment

import (
	"time"

	"github.com/google/uuid"
)

type Params struct {
	CourseID        uuid.UUID
	CourseVersionID uuid.UUID
	AccountID       uuid.UUID
	ActivatedAt     time.Time
	DeactivatedAt   time.Time
}
