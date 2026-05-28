package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	bankports "gitflic.ru/lms/backend/internal/application/ports/bank"
	quizports "gitflic.ru/lms/backend/internal/application/ports/quiz"
	bankdomain "gitflic.ru/lms/backend/internal/domain/bank"
	"gitflic.ru/lms/backend/internal/domain/bank/title"
	"github.com/google/uuid"
)

var (
	_ bankports.Repository            = (*BankRepository)(nil)
	_ bankports.QueryService          = (*BankRepository)(nil)
	_ bankports.AuditRecorder         = (*BankRepository)(nil)
	_ quizports.QuestionBankInspector = (*BankRepository)(nil)
)

// BankRepository persists question banks and exposes question-bank read models.
type BankRepository struct {
	db *sql.DB
}

// NewBankRepository creates a PostgreSQL question bank adapter.
func NewBankRepository(db *sql.DB) *BankRepository {
	return &BankRepository{db: db}
}

func (r *BankRepository) FindByID(ctx context.Context, id uuid.UUID) (*bankdomain.Bank, error) {
	var titleValue string
	err := r.db.QueryRowContext(ctx, `
		SELECT title
		FROM question_banks
		WHERE id = $1 AND deleted_at IS NULL`, id).Scan(&titleValue)
	if err != nil {
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT question_id
		FROM question_bank_questions
		WHERE bank_id = $1
		ORDER BY position`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questionIDs []uuid.UUID
	for rows.Next() {
		var questionID uuid.UUID
		if err := rows.Scan(&questionID); err != nil {
			return nil, err
		}
		questionIDs = append(questionIDs, questionID)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	t, err := title.New(titleValue)
	if err != nil {
		return nil, fmt.Errorf("restore bank title: %w", err)
	}
	return bankdomain.Restore(id, t, questionIDs)
}

func (r *BankRepository) Save(ctx context.Context, b *bankdomain.Bank) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
		INSERT INTO question_banks (id, title, updated_at)
		VALUES ($1, $2, now())
		ON CONFLICT (id) DO UPDATE SET
			title = EXCLUDED.title,
			deleted_at = NULL,
			updated_at = now()`,
		b.ID(), b.Title().Value())
	if err != nil {
		return fmt.Errorf("save bank: %w", err)
	}

	if _, err = tx.ExecContext(ctx, `DELETE FROM question_bank_questions WHERE bank_id = $1`, b.ID()); err != nil {
		return fmt.Errorf("replace bank questions: %w", err)
	}
	for idx, questionID := range b.Questions() {
		_, err = tx.ExecContext(ctx, `
			INSERT INTO question_bank_questions (bank_id, question_id, position)
			VALUES ($1, $2, $3)`, b.ID(), questionID, idx)
		if err != nil {
			return fmt.Errorf("insert bank question: %w", err)
		}
	}
	return tx.Commit()
}

func (r *BankRepository) DeleteByID(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE question_banks
		SET deleted_at = now(), updated_at = now()
		WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete bank: %w", err)
	}
	return nil
}

func (r *BankRepository) List(ctx context.Context, filter bankports.Filter) ([]bankports.ShortView, error) {
	if filter.Limit <= 0 {
		filter.Limit = 50
	}
	rows, err := r.db.QueryContext(ctx, `
		SELECT b.id::text, b.title, count(q.id)::int
		FROM question_banks b
		LEFT JOIN question_bank_questions bq ON bq.bank_id = b.id
		LEFT JOIN questions q ON q.id = bq.question_id AND q.deleted_at IS NULL
		WHERE b.deleted_at IS NULL
			AND ($1 = '' OR lower(b.title) LIKE lower('%' || $1 || '%'))
			AND ($2::uuid IS NULL OR EXISTS (
				SELECT 1 FROM question_bank_questions x
				JOIN questions xq ON xq.id = x.question_id AND xq.deleted_at IS NULL
				WHERE x.bank_id = b.id AND x.question_id = $2
			))
		GROUP BY b.id, b.title
		HAVING ($3 = 0 OR count(q.id) >= $3)
			AND ($4 = 0 OR count(q.id) <= $4)
		ORDER BY b.title, b.id
		LIMIT $5 OFFSET $6`,
		filter.TitleContains, nullUUID(filter.QuestionID), filter.MinQuestions, filter.MaxQuestions, filter.Limit, filter.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var views []bankports.ShortView
	for rows.Next() {
		var view bankports.ShortView
		if err := rows.Scan(&view.ID, &view.Title, &view.QuestionsCount); err != nil {
			return nil, err
		}
		views = append(views, view)
	}
	return views, rows.Err()
}

func (r *BankRepository) GetDetailsByID(ctx context.Context, id uuid.UUID) (bankports.DetailedView, error) {
	var view bankports.DetailedView
	rows, err := r.db.QueryContext(ctx, `
		SELECT b.id::text, b.title, q.id::text, q.type, q.payload
		FROM question_banks b
		LEFT JOIN question_bank_questions bq ON bq.bank_id = b.id
		LEFT JOIN questions q ON q.id = bq.question_id AND q.deleted_at IS NULL
		WHERE b.id = $1 AND b.deleted_at IS NULL
		ORDER BY bq.position`, id)
	if err != nil {
		return bankports.DetailedView{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var questionID sql.NullString
		var questionType sql.NullString
		var payload []byte
		if err := rows.Scan(&view.ID, &view.Title, &questionID, &questionType, &payload); err != nil {
			return bankports.DetailedView{}, err
		}
		if questionID.Valid {
			view.QuestionIDs = append(view.QuestionIDs, questionID.String)
			questionView := bankports.QuestionView{
				ID:   questionID.String,
				Type: questionType.String,
			}

			var p questionPayload
			if len(payload) > 0 {
				if err := json.Unmarshal(payload, &p); err != nil {
					return bankports.DetailedView{}, fmt.Errorf("unmarshal bank question payload: %w", err)
				}
			}
			questionView.Title = p.Title

			for _, option := range p.Selectable {
				questionView.SelectableOptions = append(questionView.SelectableOptions, bankports.SelectableOptionView{
					ID:        option.ID,
					Value:     option.Value,
					IsCorrect: option.IsCorrect,
				})
			}
			for _, option := range p.Sequence {
				questionView.SequenceOptions = append(questionView.SequenceOptions, bankports.SequenceOptionView{
					Value: option,
				})
			}
			for _, pair := range p.Matching {
				questionView.MatchingPairs = append(questionView.MatchingPairs, bankports.MatchingPairView{
					PromptID:   pair.PromptID,
					PromptText: pair.PromptText,
					MatchID:    pair.MatchID,
					MatchText:  pair.MatchText,
				})
			}
			for _, variant := range p.Short {
				questionView.ShortVariants = append(questionView.ShortVariants, bankports.ShortVariantView{
					Value: variant,
				})
			}
			view.Questions = append(view.Questions, questionView)
		}
	}
	if err := rows.Err(); err != nil {
		return bankports.DetailedView{}, err
	}
	if view.ID == "" {
		return bankports.DetailedView{}, sql.ErrNoRows
	}
	return view, nil
}

func (r *BankRepository) CountQuestionsInBank(ctx context.Context, bankID uuid.UUID) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, `
		SELECT count(*)::int
		FROM question_bank_questions
		WHERE bank_id = $1`, bankID).Scan(&count)
	return count, err
}

func (r *BankRepository) QuestionsBelongToBank(ctx context.Context, bankID uuid.UUID, questionIDs []uuid.UUID) (bool, error) {
	if len(questionIDs) == 0 {
		return true, nil
	}
	var count int
	err := r.db.QueryRowContext(ctx, `
		SELECT count(DISTINCT question_id)::int
		FROM question_bank_questions
		WHERE bank_id = $1 AND question_id = ANY($2::uuid[])`,
		bankID, uuidArray(questionIDs)).Scan(&count)
	if err != nil {
		return false, err
	}
	return count == len(questionIDs), nil
}

func (r *BankRepository) RecordBankAudit(ctx context.Context, event bankports.AuditEvent) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal bank audit event: %w", err)
	}
	var entityID any
	if event.BankID != "" {
		entityID = event.BankID
	}
	_, err = r.db.ExecContext(ctx, `
		INSERT INTO audit_events (action, entity_id, actor_role, payload)
		VALUES ($1, $2, $3, $4)`,
		string(event.Action), entityID, event.ActorRole, payload)
	if err != nil {
		return fmt.Errorf("record bank audit: %w", err)
	}
	return nil
}
