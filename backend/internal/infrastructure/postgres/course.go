package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	courseports "gitflic.ru/lms/backend/internal/application/ports/course"
	coursedomain "gitflic.ru/lms/backend/internal/domain/course"
	"gitflic.ru/lms/backend/internal/domain/course/block"
	blocktitle "gitflic.ru/lms/backend/internal/domain/course/block/title"
	"gitflic.ru/lms/backend/internal/domain/course/element"
	elementtitle "gitflic.ru/lms/backend/internal/domain/course/element/title"
	"gitflic.ru/lms/backend/internal/domain/course/progress"
	coursetitle "gitflic.ru/lms/backend/internal/domain/course/title"
	"github.com/google/uuid"
)

var (
	_ courseports.CourseRepository   = (*CourseRepository)(nil)
	_ courseports.BlockRepository    = (*BlockRepository)(nil)
	_ courseports.ElementRepository  = (*ElementRepository)(nil)
	_ courseports.ProgressRepository = (*ProgressRepository)(nil)
	_ courseports.EnrollmentAccess   = (*CoursePolicy)(nil)
	_ courseports.QueryService       = (*CourseQueryService)(nil)
)

// CourseRepository persists course aggregates.
type CourseRepository struct{ db *sql.DB }

func NewCourseRepository(db *sql.DB) *CourseRepository { return &CourseRepository{db: db} }

func (r *CourseRepository) FindByID(ctx context.Context, id uuid.UUID) (*coursedomain.Course, error) {
	var titleValue string
	if err := r.db.QueryRowContext(ctx, `SELECT title FROM courses WHERE id = $1`, id).Scan(&titleValue); err != nil {
		return nil, err
	}
	rows, err := r.db.QueryContext(ctx, `
		SELECT block_id
		FROM course_blocks_links
		WHERE course_id = $1
		ORDER BY position`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var blockIDs []uuid.UUID
	for rows.Next() {
		var blockID uuid.UUID
		if err := rows.Scan(&blockID); err != nil {
			return nil, err
		}
		blockIDs = append(blockIDs, blockID)
	}
	t, err := coursetitle.New(titleValue)
	if err != nil {
		return nil, err
	}
	return coursedomain.Restore(id, t, blockIDs)
}

func (r *CourseRepository) Save(ctx context.Context, c *coursedomain.Course) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err = tx.ExecContext(ctx, `
		INSERT INTO courses (id, title, updated_at) VALUES ($1, $2, now())
		ON CONFLICT (id) DO UPDATE SET title = EXCLUDED.title, updated_at = now()`,
		c.ID(), c.Title().Value()); err != nil {
		return fmt.Errorf("save course: %w", err)
	}
	if _, err = tx.ExecContext(ctx, `DELETE FROM course_blocks_links WHERE course_id = $1`, c.ID()); err != nil {
		return err
	}
	for idx, blockID := range c.BlockIDs() {
		if _, err = tx.ExecContext(ctx, `
			INSERT INTO course_blocks_links (course_id, block_id, position)
			VALUES ($1, $2, $3)`, c.ID(), blockID, idx); err != nil {
			return fmt.Errorf("link course block: %w", err)
		}
	}
	return tx.Commit()
}

// BlockRepository persists course block aggregates.
type BlockRepository struct{ db *sql.DB }

func NewBlockRepository(db *sql.DB) *BlockRepository { return &BlockRepository{db: db} }

func (r *BlockRepository) FindByID(ctx context.Context, id uuid.UUID) (*block.Block, error) {
	var titleValue string
	if err := r.db.QueryRowContext(ctx, `SELECT title FROM course_blocks WHERE id = $1`, id).Scan(&titleValue); err != nil {
		return nil, err
	}
	rows, err := r.db.QueryContext(ctx, `SELECT element_id FROM course_block_elements WHERE block_id = $1 ORDER BY position`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var elementIDs []uuid.UUID
	for rows.Next() {
		var elementID uuid.UUID
		if err := rows.Scan(&elementID); err != nil {
			return nil, err
		}
		elementIDs = append(elementIDs, elementID)
	}
	t, err := blocktitle.New(titleValue)
	if err != nil {
		return nil, err
	}
	return block.Restore(id, t, elementIDs)
}

func (r *BlockRepository) Save(ctx context.Context, b *block.Block) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err = tx.ExecContext(ctx, `
		INSERT INTO course_blocks (id, title, updated_at) VALUES ($1, $2, now())
		ON CONFLICT (id) DO UPDATE SET title = EXCLUDED.title, updated_at = now()`,
		b.ID(), b.Title().Value()); err != nil {
		return fmt.Errorf("save block: %w", err)
	}
	if _, err = tx.ExecContext(ctx, `DELETE FROM course_block_elements WHERE block_id = $1`, b.ID()); err != nil {
		return err
	}
	for idx, elementID := range b.ElementIDs() {
		if _, err = tx.ExecContext(ctx, `
			INSERT INTO course_block_elements (block_id, element_id, position)
			VALUES ($1, $2, $3)`, b.ID(), elementID, idx); err != nil {
			return fmt.Errorf("link block element: %w", err)
		}
	}
	return tx.Commit()
}

// ElementRepository persists course elements through explicit content DTOs.
type ElementRepository struct{ db *sql.DB }

func NewElementRepository(db *sql.DB) *ElementRepository { return &ElementRepository{db: db} }

func (r *ElementRepository) FindByID(ctx context.Context, id uuid.UUID) (*element.Element, error) {
	var titleValue, contentType, completionModeValue string
	var payload []byte
	err := r.db.QueryRowContext(ctx, `
		SELECT title, type, payload, completion_mode
		FROM course_elements
		WHERE id = $1`, id).Scan(&titleValue, &contentType, &payload, &completionModeValue)
	if err != nil {
		return nil, err
	}
	t, err := elementtitle.New(titleValue)
	if err != nil {
		return nil, err
	}
	content, mode, err := unmarshalElementPayload(contentType, payload)
	if err != nil {
		return nil, err
	}
	if completionModeValue != "" {
		mode = element.CompletionMode(completionModeValue)
	}
	return element.RestoreWithCompletionMode(id, t, content, mode)
}

func (r *ElementRepository) Save(ctx context.Context, e *element.Element) error {
	payload, objectKey, quizID, storageType, err := marshalElementPayload(e.Content(), e.CompletionMode())
	if err != nil {
		return err
	}
	var quizIDValue any
	if quizID != uuid.Nil {
		quizIDValue = quizID
	}
	_, err = r.db.ExecContext(ctx, `
		INSERT INTO course_elements (id, type, title, object_key, quiz_id, payload, completion_mode, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, now())
		ON CONFLICT (id) DO UPDATE SET
			type = EXCLUDED.type,
			title = EXCLUDED.title,
			object_key = EXCLUDED.object_key,
			quiz_id = EXCLUDED.quiz_id,
			payload = EXCLUDED.payload,
			completion_mode = EXCLUDED.completion_mode,
			updated_at = now()`,
		e.ID(), storageType, e.Title().Value(), nullString(objectKey), quizIDValue, payload, e.CompletionMode().String())
	if err != nil {
		return fmt.Errorf("save element: %w", err)
	}
	return nil
}

// ProgressRepository persists student progress markers.
type ProgressRepository struct{ db *sql.DB }

func NewProgressRepository(db *sql.DB) *ProgressRepository { return &ProgressRepository{db: db} }

func (r *ProgressRepository) FindByEnrollmentID(ctx context.Context, enrollmentID uuid.UUID) (*progress.Progress, error) {
	var id uuid.UUID
	if err := r.db.QueryRowContext(ctx, `
		SELECT id FROM course_progress WHERE enrollment_id = $1`, enrollmentID).Scan(&id); err != nil {
		return nil, err
	}
	rows, err := r.db.QueryContext(ctx, `
		SELECT element_id, marker_type, completed_at
		FROM course_progress_markers
		WHERE progress_id = $1`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	markers := map[uuid.UUID]progress.Marker{}
	for rows.Next() {
		var elementID uuid.UUID
		var markerType string
		var completedAt time.Time
		if err := rows.Scan(&elementID, &markerType, &completedAt); err != nil {
			return nil, err
		}
		p, err := progress.New(enrollmentID)
		if err != nil {
			return nil, err
		}
		if err := p.MarkElement(elementID, progress.MarkerType(markerType), completedAt); err != nil {
			return nil, err
		}
		markers[elementID] = p.Markers()[elementID]
	}
	return progress.Restore(id, enrollmentID, markers)
}

func (r *ProgressRepository) Save(ctx context.Context, p *progress.Progress) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err = tx.ExecContext(ctx, `
		INSERT INTO course_progress (id, enrollment_id, completed_elements, updated_at)
		VALUES ($1, $2, $3, now())
		ON CONFLICT (enrollment_id) DO UPDATE SET
			completed_elements = EXCLUDED.completed_elements,
			updated_at = now()`,
		p.ID(), p.EnrollmentID(), p.CompletedCount()); err != nil {
		return fmt.Errorf("save progress: %w", err)
	}
	if _, err = tx.ExecContext(ctx, `DELETE FROM course_progress_markers WHERE progress_id = $1`, p.ID()); err != nil {
		return err
	}
	for elementID, marker := range p.Markers() {
		if _, err = tx.ExecContext(ctx, `
			INSERT INTO course_progress_markers (progress_id, element_id, marker_type, completed_at)
			VALUES ($1, $2, $3, $4)`, p.ID(), elementID, marker.Type(), marker.CompletedAt()); err != nil {
			return fmt.Errorf("save progress marker: %w", err)
		}
	}
	return tx.Commit()
}

// CoursePolicy implements course access checks from enrollment/version state.
type CoursePolicy struct{ db *sql.DB }

func NewCoursePolicy(db *sql.DB) *CoursePolicy { return &CoursePolicy{db: db} }

func (p *CoursePolicy) CanViewCourse(ctx context.Context, accountID, courseID uuid.UUID, at time.Time) (bool, error) {
	var ok bool
	err := p.db.QueryRowContext(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM enrollments
			WHERE account_id = $1 AND course_id = $2
				AND enrolled_at <= $3 AND completed_at >= $3
		)`, accountID, courseID, at).Scan(&ok)
	return ok, err
}

func (p *CoursePolicy) CanEnrollCourse(ctx context.Context, courseID uuid.UUID) (bool, error) {
	var ok bool
	err := p.db.QueryRowContext(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM courses
			WHERE id = $1
		)`, courseID).Scan(&ok)
	return ok, err
}

// CourseQueryService serves course read models and student statistics.
type CourseQueryService struct{ db *sql.DB }

func NewCourseQueryService(db *sql.DB) *CourseQueryService { return &CourseQueryService{db: db} }

func (q *CourseQueryService) ListManageable(ctx context.Context, filter courseports.Filter) ([]courseports.ShortView, error) {
	return q.listCourses(ctx, filter, uuid.Nil)
}

func (q *CourseQueryService) ListVisibleForStudent(ctx context.Context, filter courseports.Filter) ([]courseports.ShortView, error) {
	return q.listCourses(ctx, filter, filter.AccountID)
}

func (q *CourseQueryService) GetDetailsByID(ctx context.Context, id uuid.UUID) (courseports.DetailedView, error) {
	var view courseports.DetailedView
	err := q.db.QueryRowContext(ctx, `
		SELECT c.id::text, c.title
		FROM courses c
		WHERE c.id = $1`, id).Scan(&view.ID, &view.Title)
	if err != nil {
		return courseports.DetailedView{}, err
	}
	return view, nil
}

func (q *CourseQueryService) ListRatings(ctx context.Context, courseID uuid.UUID, limit, offset int) ([]courseports.StudentRatingView, error) {
	filter := courseports.StudentStatisticsFilter{Limit: limit, Offset: offset}
	return q.listStudentStatistics(ctx, filter, courseID)
}

func (q *CourseQueryService) ListStudentStatistics(ctx context.Context, filter courseports.StudentStatisticsFilter) ([]courseports.StudentRatingView, error) {
	return q.listStudentStatistics(ctx, filter, uuid.Nil)
}

func (q *CourseQueryService) listCourses(ctx context.Context, filter courseports.Filter, visibleAccountID uuid.UUID) ([]courseports.ShortView, error) {
	if filter.Limit <= 0 {
		filter.Limit = 50
	}
	rows, err := q.db.QueryContext(ctx, `
		SELECT c.id::text, c.title,
			true AS published,
			count(cbl.block_id)::int AS blocks_count
		FROM courses c
		LEFT JOIN course_blocks_links cbl ON cbl.course_id = c.id
		WHERE ($1 = '' OR lower(c.title) LIKE lower('%%' || $1 || '%%'))
			AND ($3::uuid IS NULL OR EXISTS (
				SELECT 1 FROM enrollments e WHERE e.course_id = c.id AND e.account_id = $3
			))
		GROUP BY c.id, c.title
		ORDER BY c.title, c.id
		LIMIT $4 OFFSET $5`, filter.TitleContains, nullUUID(visibleAccountID), filter.Limit, filter.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var views []courseports.ShortView
	for rows.Next() {
		var view courseports.ShortView
		if err := rows.Scan(&view.ID, &view.Title, &view.Published, &view.BlocksCount); err != nil {
			return nil, err
		}
		views = append(views, view)
	}
	return views, rows.Err()
}

func (q *CourseQueryService) listStudentStatistics(ctx context.Context, filter courseports.StudentStatisticsFilter, courseID uuid.UUID) ([]courseports.StudentRatingView, error) {
	if filter.Limit <= 0 {
		filter.Limit = 50
	}
	rows, err := q.db.QueryContext(ctx, `
		SELECT e.account_id::text, e.id::text, e.course_id::text,
			CASE WHEN cp.total_elements = 0 THEN 0 ELSE (cp.completed_elements * 100 / cp.total_elements)::int END,
			cp.completed_elements, cp.total_elements
		FROM enrollments e
		JOIN course_progress cp ON cp.enrollment_id = e.id
		JOIN accounts a ON a.id = e.account_id
		JOIN persons p ON p.id = a.person_id
		WHERE ($1::uuid IS NULL OR e.account_id = $1)
			AND ($2::uuid IS NULL OR p.organization_id = $2)
			AND ($3::uuid IS NULL OR e.course_id = $3)
		ORDER BY e.account_id, e.course_id
		LIMIT $4 OFFSET $5`, nullUUID(filter.AccountID), nullUUID(filter.OrganizationID), nullUUID(courseID), filter.Limit, filter.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var views []courseports.StudentRatingView
	for rows.Next() {
		var view courseports.StudentRatingView
		if err := rows.Scan(&view.AccountID, &view.EnrollmentID, &view.CourseID, &view.CompletionPercent, &view.CompletedItems, &view.TotalItems); err != nil {
			return nil, err
		}
		views = append(views, view)
	}
	return views, rows.Err()
}
