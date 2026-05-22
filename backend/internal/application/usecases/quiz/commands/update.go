package commands

import (
	"context"
	"fmt"

	quizports "gitflic.ru/lms/backend/internal/application/ports/quiz"
	"gitflic.ru/lms/backend/internal/application/usecases/quiz/common"
	quizdomain "gitflic.ru/lms/backend/internal/domain/quiz"
)

type RenameUseCase struct {
	r quizports.Repository
}

func NewRenameUseCase(r quizports.Repository) *RenameUseCase {
	if r == nil {
		panic("quiz rename usecase requires repository")
	}
	return &RenameUseCase{r: r}
}

func (uc *RenameUseCase) Execute(ctx context.Context, in RenameInput) error {
	if err := common.RequireManager(in.ActorRole); err != nil {
		return err
	}

	q, err := loadQuiz(ctx, uc.r, in.QuizID)
	if err != nil {
		return err
	}

	t, err := buildTitle(in.Title)
	if err != nil {
		return err
	}

	if err := q.Rename(t); err != nil {
		return fmt.Errorf("rename quiz: %w", err)
	}

	return saveQuiz(ctx, uc.r, q)
}

type ChangeLimitsUseCase struct {
	r quizports.Repository
}

func NewChangeLimitsUseCase(r quizports.Repository) *ChangeLimitsUseCase {
	if r == nil {
		panic("quiz change limits usecase requires repository")
	}
	return &ChangeLimitsUseCase{r: r}
}

func (uc *ChangeLimitsUseCase) Execute(ctx context.Context, in ChangeLimitsInput) error {
	if err := common.RequireManager(in.ActorRole); err != nil {
		return err
	}

	q, err := loadQuiz(ctx, uc.r, in.QuizID)
	if err != nil {
		return err
	}

	attempts, err := buildAttempts(in.MaxAttempts)
	if err != nil {
		return err
	}
	timeLimit, err := buildTimeLimit(in.TimeLimitSeconds)
	if err != nil {
		return err
	}

	if err := q.ChangeMaxAttempts(attempts); err != nil {
		return fmt.Errorf("change quiz attempts limit: %w", err)
	}
	if err := q.ChangeTimeLimit(timeLimit); err != nil {
		return fmt.Errorf("change quiz time limit: %w", err)
	}

	return saveQuiz(ctx, uc.r, q)
}

type ChangeShufflePolicyUseCase struct {
	r quizports.Repository
}

func NewChangeShufflePolicyUseCase(r quizports.Repository) *ChangeShufflePolicyUseCase {
	if r == nil {
		panic("quiz change shuffle policy usecase requires repository")
	}
	return &ChangeShufflePolicyUseCase{r: r}
}

func (uc *ChangeShufflePolicyUseCase) Execute(ctx context.Context, in ChangeShufflePolicyInput) error {
	if err := common.RequireManager(in.ActorRole); err != nil {
		return err
	}

	q, err := loadQuiz(ctx, uc.r, in.QuizID)
	if err != nil {
		return err
	}

	q.ChangeShuffleQuestions(in.ShuffleQuestions)
	return saveQuiz(ctx, uc.r, q)
}

type AddSourceUseCase struct {
	r         quizports.Repository
	inspector quizports.QuestionBankInspector
}

func NewAddSourceUseCase(r quizports.Repository, inspector quizports.QuestionBankInspector) *AddSourceUseCase {
	if r == nil {
		panic("quiz add source usecase requires repository")
	}
	if inspector == nil {
		panic("quiz add source usecase requires question bank inspector")
	}
	return &AddSourceUseCase{r: r, inspector: inspector}
}

func (uc *AddSourceUseCase) Execute(ctx context.Context, in ChangeSourceInput) error {
	if err := common.RequireManager(in.ActorRole); err != nil {
		return err
	}

	q, err := loadQuiz(ctx, uc.r, in.QuizID)
	if err != nil {
		return err
	}

	s, err := buildSource(ctx, uc.inspector, in.Source)
	if err != nil {
		return err
	}

	if err := q.AddSource(s); err != nil {
		return fmt.Errorf("add quiz source: %w", err)
	}

	return saveQuiz(ctx, uc.r, q)
}

type ChangeSourceUseCase struct {
	r         quizports.Repository
	inspector quizports.QuestionBankInspector
}

func NewChangeSourceUseCase(r quizports.Repository, inspector quizports.QuestionBankInspector) *ChangeSourceUseCase {
	if r == nil {
		panic("quiz change source usecase requires repository")
	}
	if inspector == nil {
		panic("quiz change source usecase requires question bank inspector")
	}
	return &ChangeSourceUseCase{r: r, inspector: inspector}
}

func (uc *ChangeSourceUseCase) Execute(ctx context.Context, in ChangeSourceInput) error {
	if err := common.RequireManager(in.ActorRole); err != nil {
		return err
	}

	q, err := loadQuiz(ctx, uc.r, in.QuizID)
	if err != nil {
		return err
	}

	s, err := buildSource(ctx, uc.inspector, in.Source)
	if err != nil {
		return err
	}

	if err := q.ChangeSource(s); err != nil {
		return fmt.Errorf("change quiz source: %w", err)
	}

	return saveQuiz(ctx, uc.r, q)
}

type RemoveSourceUseCase struct {
	r quizports.Repository
}

func NewRemoveSourceUseCase(r quizports.Repository) *RemoveSourceUseCase {
	if r == nil {
		panic("quiz remove source usecase requires repository")
	}
	return &RemoveSourceUseCase{r: r}
}

func (uc *RemoveSourceUseCase) Execute(ctx context.Context, in ChangeSourceInput) error {
	if err := common.RequireManager(in.ActorRole); err != nil {
		return err
	}

	q, err := loadQuiz(ctx, uc.r, in.QuizID)
	if err != nil {
		return err
	}

	bankID, err := parseRequiredUUID(in.Source.BankID, "question bank id")
	if err != nil {
		return err
	}

	if err := q.RemoveSourceByBankID(bankID); err != nil {
		return fmt.Errorf("remove quiz source: %w", err)
	}

	return saveQuiz(ctx, uc.r, q)
}

type ReplaceSourcesUseCase struct {
	r         quizports.Repository
	inspector quizports.QuestionBankInspector
}

func NewReplaceSourcesUseCase(r quizports.Repository, inspector quizports.QuestionBankInspector) *ReplaceSourcesUseCase {
	if r == nil {
		panic("quiz replace sources usecase requires repository")
	}
	if inspector == nil {
		panic("quiz replace sources usecase requires question bank inspector")
	}
	return &ReplaceSourcesUseCase{r: r, inspector: inspector}
}

func (uc *ReplaceSourcesUseCase) Execute(ctx context.Context, in ReplaceSourcesInput) error {
	if err := common.RequireManager(in.ActorRole); err != nil {
		return err
	}

	q, err := loadQuiz(ctx, uc.r, in.QuizID)
	if err != nil {
		return err
	}

	sources, err := buildSources(ctx, uc.inspector, in.Sources)
	if err != nil {
		return err
	}

	if err := q.ReplaceSources(sources); err != nil {
		return fmt.Errorf("replace quiz sources: %w", err)
	}

	return saveQuiz(ctx, uc.r, q)
}

func saveQuiz(ctx context.Context, r quizports.Repository, q *quizdomain.Quiz) error {
	if err := r.Save(ctx, q); err != nil {
		return fmt.Errorf("save quiz: %w", err)
	}
	return nil
}
