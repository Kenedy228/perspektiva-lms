package block

import "gitflic.ru/lms/internal/domain/element"

func copyElements(original []*element.Element) []*element.Element {
	elements := make([]*element.Element, 0, len(original))

	for i := range original {
		elements = append(elements, original[i].Clone())
	}

	return elements
}
