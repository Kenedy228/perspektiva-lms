package attempt_test

import "gitflic.ru/lms/internal/domain/question"

func newAnswer() *mockAnswer {
	ans := new(mockAnswer)
	ans.On("Clone").Return(ans)
	return ans
}

func makeQuestions(correctCount, incorrectCount int) ([]question.Question, []question.Answer) {
	createQ := func(correct bool, answer question.Answer) question.Question {
		q := new(mockQuestion)
		q.On("Clone").Return(q)
		q.On("CheckAnswer", answer).Return(correct)
		return q
	}

	questions := make([]question.Question, 0, correctCount+incorrectCount)
	answers := make([]question.Answer, 0, correctCount+incorrectCount)

	for range correctCount {
		answer := newAnswer()
		q := createQ(true, answer)
		questions = append(questions, q)
		answers = append(answers, answer)
	}

	for range incorrectCount {
		answer := newAnswer()
		q := createQ(false, answer)
		questions = append(questions, q)
		answers = append(answers, answer)
	}

	return questions, answers
}
