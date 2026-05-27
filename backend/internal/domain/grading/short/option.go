package short

import "strings"

// Option настраивает поведение Checker коротких ответов.
type Option func(*Checker)

// TrimSpace добавляет нормализацию, удаляющую пробелы по краям.
func TrimSpace() Option {
	return func(c *Checker) {
		c.normalizers = append(c.normalizers, strings.TrimSpace)
	}
}

// ToLower добавляет нормализацию, приводящую строку к нижнему регистру.
func ToLower() Option {
	return func(c *Checker) {
		c.normalizers = append(c.normalizers, strings.ToLower)
	}
}
