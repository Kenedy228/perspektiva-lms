package duplicate

type Equaler[T any] interface {
	Equal(T) bool
}

func FindAllComparable[T comparable](src []T) []T {
	duplicates := make([]T, 0, len(src))
	visited := make(map[T]struct{}, len(src))

	for i := range src {
		if _, ok := visited[src[i]]; ok {
			duplicates = append(duplicates, src[i])
		}

		visited[src[i]] = struct{}{}
	}

	return duplicates
}

func FindAllNonComparable[T Equaler[T]](src []T) []T {
	duplicates := make([]T, 0, len(src))

	for i := range src {
		for j := range i + 1 {
			if src[i].Equal(src[j]) {
				duplicates = append(duplicates, src[i])
				break
			}
		}
	}

	return duplicates
}
