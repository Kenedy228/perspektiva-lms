package main

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"time"

	accountcommands "gitflic.ru/lms/backend/internal/application/usecases/account/commands"
	accountqueries "gitflic.ru/lms/backend/internal/application/usecases/account/queries"
	attemptcommands "gitflic.ru/lms/backend/internal/application/usecases/attempt/commands"
	bankcommands "gitflic.ru/lms/backend/internal/application/usecases/bank/commands"
	bankqueries "gitflic.ru/lms/backend/internal/application/usecases/bank/queries"
	coursecommands "gitflic.ru/lms/backend/internal/application/usecases/course/commands"
	coursequeries "gitflic.ru/lms/backend/internal/application/usecases/course/queries"
	enrollmentcommands "gitflic.ru/lms/backend/internal/application/usecases/enrollment/commands"
	orgcommands "gitflic.ru/lms/backend/internal/application/usecases/organization/commands"
	orgqueries "gitflic.ru/lms/backend/internal/application/usecases/organization/queries"
	personcommands "gitflic.ru/lms/backend/internal/application/usecases/person/commands"
	personqueries "gitflic.ru/lms/backend/internal/application/usecases/person/queries"
	questioncommands "gitflic.ru/lms/backend/internal/application/usecases/question/commands"
	questiongrading "gitflic.ru/lms/backend/internal/application/usecases/question/grading"
	quizcommands "gitflic.ru/lms/backend/internal/application/usecases/quiz/commands"
	domaingrading "gitflic.ru/lms/backend/internal/domain/grading"
	matchinggrading "gitflic.ru/lms/backend/internal/domain/grading/matching"
	"gitflic.ru/lms/backend/internal/domain/grading/registry"
	selectablegrading "gitflic.ru/lms/backend/internal/domain/grading/selectable"
	sequencegrading "gitflic.ru/lms/backend/internal/domain/grading/sequence"
	shortgrading "gitflic.ru/lms/backend/internal/domain/grading/short"
	typedgrading "gitflic.ru/lms/backend/internal/domain/grading/typed"
	"gitflic.ru/lms/backend/internal/domain/question"
	attemptinfra "gitflic.ru/lms/backend/internal/infrastructure/attempt"
	"gitflic.ru/lms/backend/internal/infrastructure/auth"
	"gitflic.ru/lms/backend/internal/infrastructure/postgres"
	transporthttp "gitflic.ru/lms/backend/internal/transport/http"
	"gitflic.ru/lms/backend/internal/transport/http/handlers"
	"gitflic.ru/lms/backend/internal/transport/http/session"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	cfg := configFromEnv()
	db, err := sql.Open("pgx", cfg.databaseURL)
	if err != nil {
		logger.Error("open database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		logger.Error("ping database", "error", err)
		os.Exit(1)
	}

	accountRepo := postgres.NewAccountRepository(db)
	organizationRepo := postgres.NewOrganizationRepository(db)
	personRepo := postgres.NewPersonRepository(db)
	bankRepo := postgres.NewBankRepository(db)
	questionRepo := postgres.NewQuestionRepository(db)
	quizRepo := postgres.NewQuizRepository(db)
	attemptRepo := postgres.NewAttemptRepository(db)
	attemptPolicy := postgres.NewAttemptPolicy(db)
	courseRepo := postgres.NewCourseRepository(db)
	versionRepo := postgres.NewVersionRepository(db)
	blockRepo := postgres.NewBlockRepository(db)
	progressRepo := postgres.NewProgressRepository(db)
	coursePolicy := postgres.NewCoursePolicy(db)
	courseQuery := postgres.NewCourseQueryService(db)
	enrollmentRepo := postgres.NewEnrollmentRepository(db)

	passwordComparer := auth.NewBcryptPasswordComparer()
	sessionManager := session.NewManager([]byte(cfg.sessionSecret), cfg.sessionTTL)
	checkerRegistry, err := registry.New(map[question.Type]domaingrading.Checker{
		question.TypeSelectable: selectablegrading.New(),
		question.TypeSequence:   sequencegrading.New(),
		question.TypeMatching:   matchinggrading.New(),
		question.TypeTyped:      typedgrading.New(),
		question.TypeShort:      shortgrading.New(),
	})
	if err != nil {
		logger.Error("create checker registry", "error", err)
		os.Exit(1)
	}

	answerValidators := map[question.Type]domaingrading.AnswerValidator{
		question.TypeSelectable: selectablegrading.NewValidator(),
		question.TypeSequence:   sequencegrading.NewValidator(),
		question.TypeMatching:   matchinggrading.NewValidator(),
		question.TypeShort:      shortgrading.NewValidator(),
	}

	authHandler := handlers.NewAuthHandler(
		accountcommands.NewAuthenticateUseCase(accountRepo, passwordComparer),
		sessionManager,
	)
	apiHandlers := &handlers.API{
		Accounts: handlers.AccountUseCases{
			Create:         accountcommands.NewCreateUseCase(accountRepo, passwordComparer, accountRepo),
			ChangeLogin:    accountcommands.NewChangeLoginUseCase(accountRepo, accountRepo),
			ChangePassword: accountcommands.NewChangePasswordUseCase(accountRepo, passwordComparer, accountRepo),
			ChangeRole:     accountcommands.NewChangeRoleUseCase(accountRepo, accountRepo),
			Block:          accountcommands.NewBlockUseCase(accountRepo, accountRepo),
			Activate:       accountcommands.NewActivateUseCase(accountRepo, accountRepo),
			Delete:         accountcommands.NewDeleteUseCase(accountRepo, accountRepo),
			List:           accountqueries.NewListQuery(accountRepo),
			Get:            accountqueries.NewGetByIDQuery(accountRepo),
		},
		Organizations: handlers.OrganizationUseCases{
			Create:    orgcommands.NewCreateUseCase(organizationRepo, organizationRepo),
			Rename:    orgcommands.NewRenameUseCase(organizationRepo, organizationRepo),
			ChangeINN: orgcommands.NewChangeINNUseCase(organizationRepo, organizationRepo),
			Delete:    orgcommands.NewDeleteByIDUseCase(organizationRepo, organizationRepo),
			ListName:  orgqueries.NewListByNameQuery(organizationRepo),
			ListINN:   orgqueries.NewListByINNQuery(organizationRepo),
			Get:       orgqueries.NewGetDetailsByIDQuery(organizationRepo),
		},
		Persons: handlers.PersonUseCases{
			Create:            personcommands.NewCreateUseCase(personRepo, personRepo),
			CreateWithProfile: personcommands.NewCreateWithProfileUseCase(personRepo, personRepo),
			Rename:            personcommands.NewRenameUseCase(personRepo, personRepo),
			AttachProfile:     personcommands.NewAttachProfileUseCase(personRepo, personRepo),
			ReplaceProfile:    personcommands.NewReplaceProfileUseCase(personRepo, personRepo),
			DetachProfile:     personcommands.NewDetachProfileUseCase(personRepo, personRepo),
			ChangeSNILS:       personcommands.NewChangeSNILSUseCase(personRepo, personRepo),
			ChangeDOB:         personcommands.NewChangeDateOfBirthUseCase(personRepo, personRepo),
			ChangeJobTitle:    personcommands.NewChangeJobTitleUseCase(personRepo, personRepo),
			ChangeEducation:   personcommands.NewChangeEducationUseCase(personRepo, personRepo),
			AssignOrg:         personcommands.NewAssignOrganizationUseCase(personRepo, personRepo),
			RemoveOrg:         personcommands.NewRemoveOrganizationUseCase(personRepo, personRepo),
			Delete:            personcommands.NewDeleteByIDUseCase(personRepo, personRepo),
			Get:               personqueries.NewGetDetailsByIdQuery(personRepo),
			ListLastName:      personqueries.NewListByLastNameQuery(personRepo),
			ListOrg:           personqueries.NewListByOrganizationIDQuery(personRepo),
			ListSNILS:         personqueries.NewListBySnilsQuery(personRepo),
		},
		Banks: handlers.BankUseCases{
			Create: bankcommands.NewCreateUseCase(bankRepo, bankRepo),
			Rename: bankcommands.NewRenameUseCase(bankRepo, bankRepo),
			Add:    bankcommands.NewAddQuestionsUseCase(bankRepo, bankRepo),
			Remove: bankcommands.NewRemoveQuestionsUseCase(bankRepo, bankRepo),
			Clear:  bankcommands.NewClearQuestionsUseCase(bankRepo, bankRepo),
			Delete: bankcommands.NewDeleteUseCase(bankRepo, bankRepo),
			List:   bankqueries.NewListQuery(bankRepo),
			Get:    bankqueries.NewGetDetailsByIDQuery(bankRepo),
		},
		Questions: handlers.QuestionUseCases{
			Create:           questioncommands.NewCreateUseCase(questionRepo),
			ChangeTitle:      questioncommands.NewChangeTitleUseCase(questionRepo),
			ChangeAttachment: questioncommands.NewChangeAttachmentUseCase(questionRepo),
			RemoveAttachment: questioncommands.NewRemoveAttachmentUseCase(questionRepo),
			Selectable:       questioncommands.NewChangeSelectableOptionsUseCase(questionRepo),
			Sequence:         questioncommands.NewChangeSequenceOptionsUseCase(questionRepo),
			Matching:         questioncommands.NewChangeMatchingPairsUseCase(questionRepo),
			Typed:            questioncommands.NewChangeTypedContentUseCase(questionRepo),
			Short:            questioncommands.NewChangeShortVariantsUseCase(questionRepo),
			Grade:            questiongrading.NewGradeUseCase(questionRepo, checkerRegistry, answerValidators),
			ValidateAnswer:   questiongrading.NewValidateAnswerUseCase(questionRepo, answerValidators),
			Repository:       questionRepo,
		},
		Quizzes: handlers.QuizUseCases{
			Create:        quizcommands.NewCreateUseCase(quizRepo, bankRepo),
			Rename:        quizcommands.NewRenameUseCase(quizRepo),
			ChangeLimits:  quizcommands.NewChangeLimitsUseCase(quizRepo),
			ChangeShuffle: quizcommands.NewChangeShufflePolicyUseCase(quizRepo),
			Replace:       quizcommands.NewReplaceSourcesUseCase(quizRepo, bankRepo),
			Repository:    quizRepo,
		},
		Attempts: handlers.AttemptUseCases{
			Start:      attemptcommands.NewStartUseCase(attemptRepo, quizRepo, attemptPolicy, questionRepo, attemptinfra.NewMathRandQuestionShuffler()),
			Answer:     attemptcommands.NewAddAnswerUseCase(attemptRepo),
			Finish:     attemptcommands.NewFinishUseCase(attemptRepo),
			Cancel:     attemptcommands.NewCancelUseCase(attemptRepo),
			Repository: attemptRepo,
		},
		Courses: handlers.CourseUseCases{
			Create:     coursecommands.NewCreateCourseUseCase(courseRepo),
			Rename:     coursecommands.NewRenameCourseUseCase(courseRepo),
			Version:    coursecommands.NewCreateVersionUseCase(courseRepo, versionRepo),
			Block:      coursecommands.NewAddBlockUseCase(versionRepo, blockRepo),
			Publish:    coursecommands.NewPublishVersionUseCase(versionRepo),
			Progress:   coursecommands.NewMarkProgressUseCase(progressRepo),
			List:       coursequeries.NewListQuery(courseQuery),
			Ratings:    coursequeries.NewRatingsQuery(courseQuery),
			Statistics: coursequeries.NewStudentStatisticsQuery(courseQuery),
			Query:      courseQuery,
		},
		Enrollments: handlers.EnrollmentUseCases{
			Create: enrollmentcommands.NewCreateUseCase(enrollmentRepo, coursePolicy, enrollmentRepo),
		},
	}

	server := transporthttp.NewServer(transporthttp.ServerConfig{
		Logger:  logger,
		Session: sessionManager,
		Handlers: transporthttp.Handlers{
			Auth: authHandler,
			API:  apiHandlers,
		},
	})

	httpServer := &http.Server{
		Addr:              cfg.httpAddr,
		Handler:           server.Handler(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	logger.Info("starting lms api", "addr", cfg.httpAddr)
	if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("serve http", "error", err)
		os.Exit(1)
	}
}

type config struct {
	httpAddr      string
	databaseURL   string
	sessionSecret string
	sessionTTL    time.Duration
}

func configFromEnv() config {
	return config{
		httpAddr:      env("HTTP_ADDR", ":8080"),
		databaseURL:   env("DATABASE_URL", "postgres://lms:lms@localhost:5433/lms?sslmode=disable"),
		sessionSecret: env("SESSION_SECRET", "local-development-secret-change-me"),
		sessionTTL:    durationEnv("SESSION_TTL", 8*time.Hour),
	}
}

func env(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func durationEnv(key string, fallback time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	parsed, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}
	return parsed
}
