package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	quizports "gitflic.ru/lms/backend/internal/application/ports/quiz"
	quizdomain "gitflic.ru/lms/backend/internal/domain/quiz"
	"gitflic.ru/lms/backend/internal/domain/quiz/limit"
	"gitflic.ru/lms/backend/internal/domain/quiz/source"
	"gitflic.ru/lms/backend/internal/domain/quiz/source/criteria"
	quiztitle "gitflic.ru/lms/backend/internal/domain/quiz/title"
	"github.com/google/uuid"
)

var _ quizports.Repository = (*QuizRepository)(nil)

// QuizRepository persists quiz aggregates.
type QuizRepository struct {
	db *sql.DB
}

// NewQuizRepository creates a PostgreSQL quiz adapter.
func NewQuizRepository(db *sql.DB) *QuizRepository {
	return &QuizRepository{db: db}
}

func (r *QuizRepository) FindByID(ctx context.Context, id uuid.UUID) (*quizdomain.Quiz, error) {
	var titleValue string
	var seconds, attempts int
	var shuffle bool
	err := r.db.QueryRowContext(ctx, `
		SELECT title, time_limit_seconds, attempts_limit, shuffle_questions
		FROM quizzes
		WHERE id = $1 AND deleted_at IS NULL`, id).Scan(&titleValue, &seconds, &attempts, &shuffle)
	if err != nil {
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT bank_id, criteria_type, question_count, array_to_string(question_ids, ',')
		FROM quiz_sources
		WHERE quiz_id = $1
		ORDER BY position`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sources []source.Source
	for rows.Next() {
		var bankID uuid.UUID
		var criteriaType string
		var questionCount int
		var questionIDsCSV string
		if err := rows.Scan(&bankID, &criteriaType, &questionCount, &questionIDsCSV); err != nil {
			return nil, err
		}
		c, err := restoreCriteria(criteriaType, questionCount, questionIDsCSV)
		if err != nil {
			return nil, err
		}
		s, err := source.NewSource(bankID, c)
		if err != nil {
			return nil, err
		}
		sources = append(sources, s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	t, err := quiztitle.New(titleValue)
	if err != nil {
		return nil, err
	}
	timeLimit, err := limit.NewTime(seconds)
	if err != nil {
		return nil, err
	}
	maxAttempts, err := limit.NewAttempts(attempts)
	if err != nil {
		return nil, err
	}
	return quizdomain.Restore(id, quizdomain.Params{
		Title:            t,
		MaxAttempts:      maxAttempts,
		TimeLimit:        timeLimit,
		ShuffleQuestions: shuffle,
		Sources:          sources,
	})
}

func (r *QuizRepository) Save(ctx context.Context, q *quizdomain.Quiz) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
		INSERT INTO quizzes (id, title, time_limit_seconds, attempts_limit, shuffle_questions, updated_at)
		VALUES ($1, $2, $3, $4, $5, now())
		ON CONFLICT (id) DO UPDATE SET
			title = EXCLUDED.title,
			time_limit_seconds = EXCLUDED.time_limit_seconds,
			attempts_limit = EXCLUDED.attempts_limit,
			shuffle_questions = EXCLUDED.shuffle_questions,
			deleted_at = NULL,
			updated_at = now()`,
		q.ID(), q.Title().Value(), q.Time().Seconds(), q.Attempts().Count(), q.ShuffleQuestions())
	if err != nil {
		return fmt.Errorf("save quiz: %w", err)
	}
	if _, err = tx.ExecContext(ctx, `DELETE FROM quiz_sources WHERE quiz_id = $1`, q.ID()); err != nil {
		return fmt.Errorf("replace quiz sources: %w", err)
	}
	for idx, src := range q.Sources() {
		questionIDs := ""
		if manual, ok := src.Criteria().(criteria.Manual); ok {
			questionIDs = uuidArray(manual.QuestionIDs())
		} else {
			questionIDs = "{}"
		}
		_, err = tx.ExecContext(ctx, `
			INSERT INTO quiz_sources (
				quiz_id, bank_id, criteria_type, question_count, question_ids, position
			)
			VALUES ($1, $2, $3, $4, $5::uuid[], $6)`,
			q.ID(), src.BankID(), src.Criteria().Type().String(), src.Criteria().QuestionCount(), questionIDs, idx)
		if err != nil {
			return fmt.Errorf("insert quiz source: %w", err)
		}
	}
	return tx.Commit()
}

func (r *QuizRepository) DeleteByID(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE quizzes SET deleted_at = now(), updated_at = now()
		WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete quiz: %w", err)
	}
	return nil
}

func restoreCriteria(criteriaType string, questionCount int, questionIDsCSV string) (criteria.Criteria, error) {
	switch criteria.Type(criteriaType) {
	case criteria.TypeRandom:
		return criteria.NewRandom(questionCount)
	case criteria.TypeManual:
		parts := strings.Split(questionIDsCSV, ",")
		ids := make([]uuid.UUID, 0, len(parts))
		for _, part := range parts {
			if part == "" {
				continue
			}
			id, err := uuid.Parse(part)
			if err != nil {
				return nil, err
			}
			ids = append(ids, id)
		}
		return criteria.NewManual(ids)
	default:
		return nil, fmt.Errorf("%w: unknown criteria type %q", ErrUnsupported, criteriaType)
	}
}
