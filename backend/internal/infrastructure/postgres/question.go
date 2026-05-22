package postgres

import (
	"context"
	"database/sql"
	"fmt"

	attemptports "gitflic.ru/lms/backend/internal/application/ports/attempt"
	questionports "gitflic.ru/lms/backend/internal/application/ports/question"
	questdomain "gitflic.ru/lms/backend/internal/domain/question"
	"github.com/google/uuid"
)

var (
	_ questionports.Repository         = (*QuestionRepository)(nil)
	_ attemptports.QuestionSetProvider = (*QuestionRepository)(nil)
)

// QuestionRepository stores questions using explicit JSON DTOs for each
// polymorphic question type.
type QuestionRepository struct {
	db *sql.DB
}

// NewQuestionRepository creates a PostgreSQL question adapter.
func NewQuestionRepository(db *sql.DB) *QuestionRepository {
	return &QuestionRepository{db: db}
}

func (r *QuestionRepository) FindByID(ctx context.Context, id uuid.UUID) (questdomain.Question, error) {
	var qType string
	var payload []byte
	err := r.db.QueryRowContext(ctx, `
		SELECT type, payload
		FROM questions
		WHERE id = $1 AND deleted_at IS NULL`, id).Scan(&qType, &payload)
	if err != nil {
		return nil, err
	}
	return unmarshalQuestion(id, qType, payload)
}

func (r *QuestionRepository) Save(ctx context.Context, q questdomain.Question) error {
	payload, err := marshalQuestion(q)
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx, `
		INSERT INTO questions (id, type, payload, updated_at)
		VALUES ($1, $2, $3, now())
		ON CONFLICT (id) DO UPDATE SET
			type = EXCLUDED.type,
			payload = EXCLUDED.payload,
			deleted_at = NULL,
			updated_at = now()`,
		q.ID(), q.Type().String(), payload)
	if err != nil {
		return fmt.Errorf("save question: %w", err)
	}
	return nil
}

func (r *QuestionRepository) DeleteByID(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE questions SET deleted_at = now(), updated_at = now()
		WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete question: %w", err)
	}
	return nil
}

func (r *QuestionRepository) FindQuestionsByIDs(ctx context.Context, bankID uuid.UUID, questionIDs []uuid.UUID) ([]questdomain.Question, error) {
	if len(questionIDs) == 0 {
		return []questdomain.Question{}, nil
	}
	rows, err := r.db.QueryContext(ctx, `
		SELECT q.id, q.type, q.payload
		FROM questions q
		JOIN question_bank_questions bq ON bq.question_id = q.id
		WHERE bq.bank_id = $1 AND q.id = ANY($2::uuid[]) AND q.deleted_at IS NULL
		ORDER BY array_position($2::uuid[], q.id)`,
		bankID, uuidArray(questionIDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanQuestions(rows)
}

func (r *QuestionRepository) SelectRandomQuestions(ctx context.Context, bankID uuid.UUID, count int) ([]questdomain.Question, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT q.id, q.type, q.payload
		FROM questions q
		JOIN question_bank_questions bq ON bq.question_id = q.id
		WHERE bq.bank_id = $1 AND q.deleted_at IS NULL
		ORDER BY random()
		LIMIT $2`, bankID, count)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanQuestions(rows)
}

func scanQuestions(rows *sql.Rows) ([]questdomain.Question, error) {
	var questions []questdomain.Question
	for rows.Next() {
		var id uuid.UUID
		var qType string
		var payload []byte
		if err := rows.Scan(&id, &qType, &payload); err != nil {
			return nil, err
		}
		q, err := unmarshalQuestion(id, qType, payload)
		if err != nil {
			return nil, err
		}
		questions = append(questions, q)
	}
	return questions, rows.Err()
}
