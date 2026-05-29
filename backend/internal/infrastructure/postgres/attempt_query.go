package postgres

import (
	"context"
	"database/sql"

	attemptports "gitflic.ru/lms/backend/internal/application/ports/attempt"
	"github.com/google/uuid"
)

var _ attemptports.QueryService = (*AttemptQueryService)(nil)

type AttemptQueryService struct {
	db *sql.DB
}

func NewAttemptQueryService(db *sql.DB) *AttemptQueryService {
	return &AttemptQueryService{db: db}
}

func (q *AttemptQueryService) ListByEnrollmentID(ctx context.Context, enrollmentID uuid.UUID, limit, offset int) ([]attemptports.AttemptSummaryView, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	rows, err := q.db.QueryContext(ctx, `
		SELECT a.id::text, a.enrollment_id::text, a.quiz_id::text, a.status,
		       a.started_at, a.deadline_at, a.submitted_at,
		       COUNT(DISTINCT ai.question_id)::int AS questions_count,
		       COUNT(DISTINCT aa.question_id)::int AS answers_count
		FROM attempts a
		LEFT JOIN attempt_items ai ON ai.attempt_id = a.id
		LEFT JOIN attempt_answers aa ON aa.attempt_id = a.id
		WHERE a.enrollment_id = $1
		GROUP BY a.id, a.enrollment_id, a.quiz_id, a.status, a.started_at, a.deadline_at, a.submitted_at
		ORDER BY a.started_at DESC
		LIMIT $2 OFFSET $3`, enrollmentID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []attemptports.AttemptSummaryView
	for rows.Next() {
		var v attemptports.AttemptSummaryView
		var deadlineAt, finishedAt sql.NullTime
		if err := rows.Scan(
			&v.ID, &v.EnrollmentID, &v.QuizID, &v.Status,
			&v.StartedAt, &deadlineAt, &finishedAt,
			&v.QuestionsCount, &v.AnswersCount,
		); err != nil {
			return nil, err
		}
		v.DeadlineAt = deadlineAt.Time
		v.FinishedAt = finishedAt.Time
		result = append(result, v)
	}
	return result, rows.Err()
}
