package block

import "gitflic.ru/lms/internal/domain/element"

type Params struct {
	Title    string
	Elements []*element.Element
}
