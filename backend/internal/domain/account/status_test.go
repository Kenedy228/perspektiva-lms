package account

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStatus_IsValid(t *testing.T) {
	tests := []struct {
		name string
		s    Status
		want bool
	}{
		{name: "active", s: StatusActive, want: true},
		{name: "blocked", s: StatusBlocked, want: true},
		{name: "deleted", s: StatusDeleted, want: true},
		{
			name: "неопределенное значение через приведение типов",
			s:    Status("status"),
			want: false,
		},
		{
			name: "неопределенное значение без приведения типов",
			s:    "status",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tt.s.IsValid(), "IsValid()")
		})
	}
}

func TestStatus_String(t *testing.T) {
	tests := []struct {
		name string
		s    Status
		want string
	}{
		{
			name: "определенное значение",
			s:    StatusActive,
			want: "active",
		},
		{
			name: "неопределенное значение",
			s:    Status("status"),
			want: "status",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tt.s.String(), "String()")
		})
	}
}

func TestStatus_Title(t *testing.T) {
	tests := []struct {
		name string
		s    Status
		want string
	}{
		{name: "active", s: StatusActive, want: "активный"},
		{name: "blocked", s: StatusBlocked, want: "заблокирован"},
		{name: "deleted", s: StatusDeleted, want: "удалён"},
		{
			name: "для неопределенных значений возвращает пустую строку",
			s:    Status("undefined"),
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tt.s.Title(), "Value()")
		})
	}
}
