//go:build integration

package postgres_test

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	orgcommands "gitflic.ru/lms/backend/internal/application/usecases/organization/commands"
	orgqueries "gitflic.ru/lms/backend/internal/application/usecases/organization/queries"
	attemptdomain "gitflic.ru/lms/backend/internal/domain/attempt"
	"gitflic.ru/lms/backend/internal/domain/question"
	selectablequestion "gitflic.ru/lms/backend/internal/domain/question/selectable"
	selectableanswer "gitflic.ru/lms/backend/internal/domain/question/selectable/answer"
	selectableoption "gitflic.ru/lms/backend/internal/domain/question/selectable/option"
	"gitflic.ru/lms/backend/internal/domain/quiz/limit"
	"gitflic.ru/lms/backend/internal/domain/role"
	"gitflic.ru/lms/backend/internal/domain/shared/text"
	"gitflic.ru/lms/backend/internal/domain/shared/title"
	"gitflic.ru/lms/backend/internal/infrastructure/postgres"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
)

func TestOrganizationUseCaseLifecycleWithPostgres(t *testing.T) {
	ctx, db := startPostgres(t)

	repo := postgres.NewOrganizationRepository(db)
	admin := role.NewAdmin()

	create := orgcommands.NewCreateUseCase(repo, repo)
	created, err := create.Execute(ctx, orgcommands.CreateInput{
		ActorRole: admin,
		INN:       "1030000000",
		INNType:   "organization",
		Name:      "ООО Ромашка",
	})
	require.NoError(t, err)
	require.NotEmpty(t, created.ID)

	rename := orgcommands.NewRenameUseCase(repo, repo)
	_, err = rename.Execute(ctx, orgcommands.RenameInput{
		ActorRole:      admin,
		OrganizationID: created.ID,
		Name:           "ООО Академия",
	})
	require.NoError(t, err)

	changeINN := orgcommands.NewChangeINNUseCase(repo, repo)
	_, err = changeINN.Execute(ctx, orgcommands.ChangeINNInput{
		ActorRole:      admin,
		OrganizationID: created.ID,
		INN:            "3664069397",
		INNType:        "organization",
	})
	require.NoError(t, err)

	get := orgqueries.NewGetDetailsByIDQuery(repo)
	details, err := get.Execute(ctx, orgqueries.GetDetailsByIDInput{ActorRole: admin, ID: created.ID})
	require.NoError(t, err)
	require.Equal(t, "ООО Академия", details.View.OrganizationName)
	require.Equal(t, "3664069397", details.View.INN)

	remove := orgcommands.NewDeleteByIDUseCase(repo, repo)
	require.NoError(t, remove.Execute(ctx, orgcommands.DeleteByIDInput{
		ActorRole:      admin,
		OrganizationID: created.ID,
	}))

	_, err = get.Execute(ctx, orgqueries.GetDetailsByIDInput{ActorRole: admin, ID: created.ID})
	require.Error(t, err)
}

func TestQuestionAndAttemptRepositoryLifecycleWithPostgres(t *testing.T) {
	ctx, db := startPostgres(t)

	questionRepo := postgres.NewQuestionRepository(db)
	attemptRepo := postgres.NewAttemptRepository(db)

	q := newSelectableQuestion(t)
	require.NoError(t, questionRepo.Save(ctx, q))

	foundQuestion, err := questionRepo.FindByID(ctx, q.ID())
	require.NoError(t, err)
	require.Equal(t, q.ID(), foundQuestion.ID())
	require.Equal(t, q.Title().Value(), foundQuestion.Title().Value())

	insertAttemptDependencies(t, ctx, db, q.ID())

	timeLimit, err := limit.NewTime(600)
	require.NoError(t, err)
	startedAt := time.Date(2026, 5, 9, 10, 0, 0, 0, time.UTC)
	attempt, err := attemptdomain.New(attemptdomain.Params{
		EnrollmentID: uuidFromString("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaa0004"),
		QuizID:       uuidFromString("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaa0005"),
		TimeLimit:    timeLimit,
		Questions:    []question.Question{q},
	}, startedAt)
	require.NoError(t, err)

	optionID, err := selectableanswer.NewOptionID(q.Options()[0].ID())
	require.NoError(t, err)
	ans, err := selectableanswer.New([]selectableanswer.OptionID{optionID})
	require.NoError(t, err)
	require.NoError(t, attempt.AddAnswer(q.ID(), ans, startedAt.Add(time.Minute)))
	require.NoError(t, attemptRepo.Save(ctx, attempt))

	restored, err := attemptRepo.FindByID(ctx, attempt.ID())
	require.NoError(t, err)
	require.Equal(t, attempt.ID(), restored.ID())
	require.Equal(t, 1, restored.CountItems())
	require.Equal(t, 1, restored.CountAnswers())
}

func startPostgres(t *testing.T) (context.Context, *sql.DB) {
	t.Helper()
	ctx := context.Background()

	container, err := tcpostgres.Run(ctx,
		"postgres:18",
		tcpostgres.WithDatabase("lms"),
		tcpostgres.WithUsername("lms"),
		tcpostgres.WithPassword("lms"),
		tcpostgres.BasicWaitStrategies(),
	)
	testcontainers.CleanupContainer(t, container)
	require.NoError(t, err)

	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	db, err := sql.Open("pgx", dsn)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, db.Close()) })
	require.NoError(t, db.PingContext(ctx))
	require.NoError(t, applyMigrations(ctx, db))
	return ctx, db
}

func applyMigrations(ctx context.Context, db *sql.DB) error {
	dir := filepath.Join("..", "..", "..", "..", "migrations")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
			names = append(names, entry.Name())
		}
	}
	sort.Strings(names)

	for _, name := range names {
		raw, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			return err
		}
		upSQL := strings.Split(string(raw), "-- +goose Down")[0]
		upSQL = strings.ReplaceAll(upSQL, "-- +goose Up", "")
		for _, stmt := range strings.Split(upSQL, ";") {
			stmt = strings.TrimSpace(stmt)
			if stmt == "" {
				continue
			}
			if _, err := db.ExecContext(ctx, stmt); err != nil {
				return err
			}
		}
	}
	return nil
}

func newSelectableQuestion(t *testing.T) *selectablequestion.Question {
	t.Helper()
	qTitle, err := title.New("Question title")
	require.NoError(t, err)
	firstText, err := text.New("first")
	require.NoError(t, err)
	secondText, err := text.New("second")
	require.NoError(t, err)
	first, err := selectableoption.New(firstText, true)
	require.NoError(t, err)
	second, err := selectableoption.New(secondText, false)
	require.NoError(t, err)
	q, err := selectablequestion.New(qTitle, []selectableoption.Option{first, second})
	require.NoError(t, err)
	return q
}

func insertAttemptDependencies(t *testing.T, ctx context.Context, db *sql.DB, questionID uuid.UUID) {
	t.Helper()
	execSQL(t, ctx, db, `INSERT INTO persons (id, first_name, last_name) VALUES ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaa0001', 'Ivan', 'Ivanov')`)
	execSQL(t, ctx, db, `INSERT INTO accounts (id, person_id, login, password_hash, role, status) VALUES ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaa0002', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaa0001', 'student1', 'hash-value', 'student', 'active')`)
	execSQL(t, ctx, db, `INSERT INTO question_banks (id, title) VALUES ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaa0003', 'Bank')`)
	execSQL(t, ctx, db, `INSERT INTO question_bank_questions (bank_id, question_id, position) VALUES ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaa0003', $1, 0)`, questionID)
	execSQL(t, ctx, db, `INSERT INTO quizzes (id, title, time_limit_seconds, attempts_limit, shuffle_questions) VALUES ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaa0005', 'Quiz', 600, 1, false)`)
	execSQL(t, ctx, db, `INSERT INTO quiz_sources (quiz_id, bank_id, criteria_type, question_count, question_ids, position) VALUES ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaa0005', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaa0003', 'manual', 1, ARRAY[$1]::uuid[], 0)`, questionID)
	execSQL(t, ctx, db, `INSERT INTO courses (id, title) VALUES ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaa0006', 'Course')`)
	execSQL(t, ctx, db, `INSERT INTO course_versions (id, title, status) VALUES ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaa0007', 'Version', 'published')`)
	execSQL(t, ctx, db, `INSERT INTO course_version_links (course_id, version_id, position) VALUES ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaa0006', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaa0007', 0)`)
	execSQL(t, ctx, db, `INSERT INTO enrollments (id, account_id, course_id, version_id, status, enrolled_at, completed_at) VALUES ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaa0004', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaa0002', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaa0006', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaa0007', 'active', '2026-01-01', '2026-12-31')`)
}

func uuidFromString(value string) uuid.UUID {
	return uuid.MustParse(value)
}

func execSQL(t *testing.T, ctx context.Context, db *sql.DB, stmt string, args ...any) {
	t.Helper()
	_, err := db.ExecContext(ctx, stmt, args...)
	require.NoError(t, err)
}
