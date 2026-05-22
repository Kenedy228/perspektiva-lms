package attempt

import (
	"testing"
	"time"

	"gitflic.ru/lms/backend/internal/domain/attempt/answer"
	"gitflic.ru/lms/backend/internal/domain/question"
	"gitflic.ru/lms/backend/internal/domain/question/attachment"
	"gitflic.ru/lms/backend/internal/domain/quiz/limit"
	"gitflic.ru/lms/backend/internal/domain/shared/title"
	"github.com/google/uuid"
)

type mockQuestion struct {
	id    uuid.UUID
	qType question.Type
}

func (q mockQuestion) ID() uuid.UUID { return q.id }

func (q mockQuestion) Title() title.Title {
	t, _ := title.New("Question")
	return t
}

func (q mockQuestion) Attachment() (attachment.Attachment, bool) {
	return attachment.Attachment{}, false
}

func (q mockQuestion) Instruction() string { return q.qType.DefaultInstruction() }
func (q mockQuestion) Type() question.Type { return q.qType }
func (q mockQuestion) Clone() question.Question {
	return mockQuestion{id: q.id, qType: q.qType}
}
func (q mockQuestion) ChangeTitle(title.Title)                {}
func (q mockQuestion) ChangeAttachment(attachment.Attachment) {}
func (q mockQuestion) RemoveAttachment()                      {}
func (q mockQuestion) HasAttachment() bool                    { return false }

func asQuestions(items ...mockQuestion) []question.Question {
	questions := make([]question.Question, 0, len(items))
	for i := range items {
		questions = append(questions, items[i])
	}
	return questions
}

func answerEntry(id uuid.UUID, ans question.Answer, at time.Time) (answer.Entry, error) {
	return answer.New(id, ans, at)
}

func validParams(t *testing.T, questions ...mockQuestion) Params {
	t.Helper()

	timeLimit, err := limit.NewTime(0)
	if err != nil {
		t.Fatalf("create time limit: %v", err)
	}

	return Params{
		EnrollmentID: uuid.New(),
		QuizID:       uuid.New(),
		TimeLimit:    timeLimit,
		Questions:    asQuestions(questions...),
	}
}

func mustTimeLimit(t *testing.T, seconds int) limit.Time {
	t.Helper()

	timeLimit, err := limit.NewTime(seconds)
	if err != nil {
		t.Fatalf("create time limit: %v", err)
	}
	return timeLimit
}
