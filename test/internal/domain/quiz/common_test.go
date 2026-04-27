package quiz_test

import "gitflic.ru/lms/internal/domain/shared/limit"

func makeLimit(val int) limit.Limit {
	l, _ := limit.New(val)
	return l
}

