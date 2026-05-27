package education

import "strings"

// normalizeValue принимает сырые сведения об образовании, удаляет незначащие пробелы
// и заменяет последовательности двух и более пробельных символов на одиночный пробел.
func normalizeValue(value string) string {
	trimmed := strings.TrimSpace(value)
	return strings.Join(strings.Fields(trimmed), " ")
}
