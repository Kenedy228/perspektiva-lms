package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	attemptports "gitflic.ru/lms/backend/internal/application/ports/attempt"
	attemptdomain "gitflic.ru/lms/backend/internal/domain/attempt"
	attemptanswer "gitflic.ru/lms/backend/internal/domain/attempt/answer"
	questdomain "gitflic.ru/lms/backend/internal/domain/question"
	"github.com/google/uuid"
)

var (
	_ attemptports.Repository       = (*AttemptRepository)(nil)
	_ attemptports.EnrollmentPolicy = (*AttemptPolicy)(nil)
)

// AttemptRepository persists attempts with question snapshots referenced from
// the question table and explicit JSON DTOs for answers.
type AttemptRepository struct {
	db *sql.DB
}

func NewAttemptRepository(db *sql.DB) *AttemptRepository { return &AttemptRepository{db: db} }

func (r *AttemptRepository) FindByID(ctx context.Context, id uuid.UUID) (*attemptdomain.Attempt, error) {
	var enrollmentID, quizID uuid.UUID
	var status string
	var startedAt time.Time
	var deadlineAt, finishedAt sql.NullTime
	err := r.db.QueryRowContext(ctx, `
		SELECT enrollment_id, quiz_id, status, started_at, deadline_at, submitted_at
		FROM attempts
		WHERE id = $1`, id).Scan(&enrollmentID, &quizID, &status, &startedAt, &deadlineAt, &finishedAt)
	if err != nil {
		return nil, err
	}

	questions, err := r.findAttemptQuestions(ctx, id)
	if err != nil {
		return nil, err
	}
	answers, err := r.findAttemptAnswers(ctx, id, questions)
	if err != nil {
		return nil, err
	}

	return attemptdomain.Restore(id, attemptdomain.RestoreParams{
		EnrollmentID: enrollmentID,
		QuizID:       quizID,
		Status:       attemptdomain.Status(status),
		StartedAt:    startedAt,
		DeadlineAt:   deadlineAt.Time,
		FinishedAt:   finishedAt.Time,
		Questions:    questions,
		Answers:      answers,
	})
}

func (r *AttemptRepository) Save(ctx context.Context, a *attemptdomain.Attempt) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
		INSERT INTO attempts (
			id, enrollment_id, quiz_id, account_id, status, started_at, deadline_at, submitted_at, updated_at
		)
		SELECT $1, $2, $3, e.account_id, $4, $5, $6, $7, now()
		FROM enrollments e
		WHERE e.id = $2
		ON CONFLICT (id) DO UPDATE SET
			status = EXCLUDED.status,
			deadline_at = EXCLUDED.deadline_at,
			submitted_at = EXCLUDED.submitted_at,
			updated_at = now()`,
		a.ID(), a.EnrollmentID(), a.QuizID(), a.Status().String(), a.StartedAt(), nullTime(a.DeadlineAt()), nullTime(a.FinishedAt()))
	if err != nil {
		return fmt.Errorf("save attempt: %w", err)
	}

	if _, err = tx.ExecContext(ctx, `DELETE FROM attempt_items WHERE attempt_id = $1`, a.ID()); err != nil {
		return err
	}
	for idx, item := range a.Items() {
		if _, err = tx.ExecContext(ctx, `
			INSERT INTO attempt_items (attempt_id, question_id, position)
			VALUES ($1, $2, $3)`, a.ID(), item.ID(), idx); err != nil {
			return fmt.Errorf("save attempt item: %w", err)
		}
	}

	if _, err = tx.ExecContext(ctx, `DELETE FROM attempt_answers WHERE attempt_id = $1`, a.ID()); err != nil {
		return err
	}
	for questionID, entry := range a.Answers() {
		payload, err := marshalAnswer(entry.Answer())
		if err != nil {
			return err
		}
		if _, err = tx.ExecContext(ctx, `
			INSERT INTO attempt_answers (attempt_id, question_id, answer_payload, answered_at)
			VALUES ($1, $2, $3, $4)`, a.ID(), questionID, payload, entry.AnsweredAt()); err != nil {
			return fmt.Errorf("save attempt answer: %w", err)
		}
	}

	return tx.Commit()
}

func (r *AttemptRepository) CountByEnrollmentAndQuiz(ctx context.Context, enrollmentID, quizID uuid.UUID) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, `
		SELECT count(*)::int FROM attempts
		WHERE enrollment_id = $1 AND quiz_id = $2`, enrollmentID, quizID).Scan(&count)
	return count, err
}

// AttemptPolicy checks whether an enrolled student can start a quiz.
type AttemptPolicy struct {
	db *sql.DB
}

func NewAttemptPolicy(db *sql.DB) *AttemptPolicy { return &AttemptPolicy{db: db} }

func (p *AttemptPolicy) CanStartQuiz(ctx context.Context, accountID, enrollmentID, quizID uuid.UUID, at time.Time) (bool, error) {
	var ok bool
	err := p.db.QueryRowContext(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM enrollments e
			JOIN course_version_links cvl ON cvl.course_id = e.course_id AND cvl.version_id = e.version_id
			JOIN course_version_blocks cvb ON cvb.version_id = cvl.version_id
			JOIN course_block_elements cbe ON cbe.block_id = cvb.block_id
			JOIN course_elements ce ON ce.id = cbe.element_id
			WHERE e.id = $1
				AND e.account_id = $2
				AND e.enrolled_at <= $4
				AND e.completed_at >= $4
				AND ce.quiz_id = $3
		)`, enrollmentID, accountID, quizID, at).Scan(&ok)
	return ok, err
}

func nullTime(value time.Time) any {
	if value.IsZero() {
		return nil
	}
	return value
}

func (r *AttemptRepository) findAttemptQuestions(ctx context.Context, attemptID uuid.UUID) ([]questdomain.Question, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT q.id, q.type, q.payload
		FROM attempt_items ai
		JOIN questions q ON q.id = ai.question_id
		WHERE ai.attempt_id = $1
		ORDER BY ai.position`, attemptID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanQuestions(rows)
}

func (r *AttemptRepository) findAttemptAnswers(ctx context.Context, attemptID uuid.UUID, questions []questdomain.Question) (map[uuid.UUID]attemptanswer.Entry, error) {
	questionTypes := make(map[uuid.UUID]questdomain.Type, len(questions))
	for _, q := range questions {
		questionTypes[q.ID()] = q.Type()
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT question_id, answer_payload, answered_at
		FROM attempt_answers
		WHERE attempt_id = $1`, attemptID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	answers := make(map[uuid.UUID]attemptanswer.Entry)
	for rows.Next() {
		var questionID uuid.UUID
		var payload []byte
		var answeredAt time.Time
		if err := rows.Scan(&questionID, &payload, &answeredAt); err != nil {
			return nil, err
		}
		qType, ok := questionTypes[questionID]
		if !ok {
			return nil, fmt.Errorf("%w: answer question %s has no attempt item", ErrUnsupported, questionID)
		}
		ans, err := unmarshalAnswer(qType, payload)
		if err != nil {
			return nil, err
		}
		entry, err := attemptanswer.New(questionID, ans, answeredAt)
		if err != nil {
			return nil, err
		}
		answers[questionID] = entry
	}
	return answers, rows.Err()
}
