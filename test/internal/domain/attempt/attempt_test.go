package attempt_test

import (
	"testing"
	"time"

	"gitflic.ru/lms/internal/domain/attempt"
	"gitflic.ru/lms/internal/domain/question"
	"gitflic.ru/lms/internal/domain/quiz/limit"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func validParams() attempt.Params {
	t, _ := limit.NewTime(0)

	return attempt.Params{
		EnrollmentID: uuid.New(),
		QuizID:       uuid.New(),
		Questions:    []question.Question{mockQuestion{id: uuid.New()}},
		TimeLimit:    t,
	}
}

func TestNew(t *testing.T) {
	now := time.Now()

	t.Run("пустой ID зачисления возвращает ошибку", func(t *testing.T) {
		p := validParams()
		p.EnrollmentID = uuid.Nil
		_, err := attempt.New(p, now)
		assert.ErrorIs(t, err, attempt.ErrInvalid)
	})

	t.Run("пустой ID теста возвращает ошибку", func(t *testing.T) {
		p := validParams()
		p.QuizID = uuid.Nil
		_, err := attempt.New(p, now)
		assert.ErrorIs(t, err, attempt.ErrInvalid)
	})

	t.Run("пустой список вопросов возвращает ошибку", func(t *testing.T) {
		p := validParams()
		p.Questions = []question.Question{}
		_, err := attempt.New(p, now)
		assert.ErrorIs(t, err, attempt.ErrInvalid)
	})

	t.Run("успешное создание попытки", func(t *testing.T) {
		p := validParams()
		a, err := attempt.New(p, now)

		assert.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, a.ID())
		assert.Equal(t, attempt.StatusInProgress, a.Status())
		assert.Equal(t, 1, a.CountItems())
		assert.Equal(t, 0, a.CountAnswers())
		assert.Equal(t, now, a.StartedAt())
		assert.True(t, a.CanModify())
	})
}

func TestAddAnswer(t *testing.T) {
	now := time.Now()
	p := validParams()
	qID := p.Questions[0].ID()

	t.Run("успешное добавление ответа", func(t *testing.T) {
		a, _ := attempt.New(p, now)
		err := a.AddAnswer(qID, mockAnswer{}, now)

		assert.NoError(t, err)
		assert.Equal(t, 1, a.CountAnswers())
		assert.NotNil(t, a.Answers()[qID])
	})

	t.Run("добавление ответа к несуществующему вопросу возвращает ошибку", func(t *testing.T) {
		a, _ := attempt.New(p, now)
		err := a.AddAnswer(uuid.New(), mockAnswer{}, now)

		assert.ErrorIs(t, err, attempt.ErrNotFound)
	})

	t.Run("добавление ответа в завершенную попытку возвращает конфликт", func(t *testing.T) {
		a, _ := attempt.New(p, now)
		a.Finish(now)

		err := a.AddAnswer(qID, mockAnswer{}, now)
		assert.ErrorIs(t, err, attempt.ErrStateConflict)
	})
}

func TestFinish(t *testing.T) {
	now := time.Now()

	t.Run("успешное завершение активной попытки", func(t *testing.T) {
		a, _ := attempt.New(validParams(), now)
		finishTime := now.Add(time.Minute)

		err := a.Finish(finishTime)

		assert.NoError(t, err)
		assert.Equal(t, attempt.StatusFinished, a.Status())
		assert.Equal(t, finishTime, a.FinishedAt())
		assert.False(t, a.CanModify())
	})

	t.Run("повторное завершение возвращает конфликт", func(t *testing.T) {
		a, _ := attempt.New(validParams(), now)
		a.Finish(now)

		err := a.Finish(now)
		assert.ErrorIs(t, err, attempt.ErrStateConflict)
	})
}

func TestSetExpired(t *testing.T) {
	now := time.Now()

	t.Run("безлимитная попытка возвращает ошибку", func(t *testing.T) {
		p := validParams() // безлимитная по умолчанию
		a, _ := attempt.New(p, now)

		err := a.SetExpired(now.Add(time.Hour))
		assert.ErrorIs(t, err, attempt.ErrStateConflict) // попытка не имеет дедлайна
	})

	// Примечание: для проверки успешного просрочивания тебе нужно передать в validParams
	// лимит (например 10 минут) и вызвать SetExpired(now.Add(15 * time.Minute))
}

func TestInterruptAndCancel(t *testing.T) {
	now := time.Now()

	t.Run("успешное прерывание", func(t *testing.T) {
		a, _ := attempt.New(validParams(), now)
		err := a.Interrupt(now)

		assert.NoError(t, err)
		assert.Equal(t, attempt.StatusInterrupted, a.Status())
	})

	t.Run("успешная отмена", func(t *testing.T) {
		a, _ := attempt.New(validParams(), now)
		err := a.Cancel()

		assert.NoError(t, err)
		assert.Equal(t, attempt.StatusCancelled, a.Status())
	})

	t.Run("отмена завершенной попытки возвращает ошибку", func(t *testing.T) {
		a, _ := attempt.New(validParams(), now)
		a.Finish(now)

		err := a.Cancel()
		assert.ErrorIs(t, err, attempt.ErrStateConflict)
	})
}

func TestStatus(t *testing.T) {
	t.Run("проверка Title и String", func(t *testing.T) {
		s := attempt.StatusFinished
		assert.Equal(t, "завершен", s.Title())
		assert.Equal(t, "finished", s.String())

		s = attempt.StatusInProgress
		assert.Equal(t, "в процессе", s.Title())
		assert.Equal(t, "in_progress", s.String())
	})
}
