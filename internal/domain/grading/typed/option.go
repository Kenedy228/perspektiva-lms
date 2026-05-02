package typed

import "strings"

type Option func(*Checker)

func TrimSpace() Option {
	return func(c *Checker) {
		c.normalizers = append(c.normalizers, strings.TrimSpace)
	}
}

func ToLower() Option {
	return func(c *Checker) {
		c.normalizers = append(c.normalizers, strings.ToLower)
	}
}
