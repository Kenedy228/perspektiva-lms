package typed

import "gitflic.ru/lms/internal/domain/question"

type BlankParams struct {
	Placeholder string
	Variants    []question.Content
}
