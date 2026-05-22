//go:build legacy
// +build legacy

package name_test

import (
	"strings"
	"testing"

	name2 "gitflic.ru/lms/backend/internal/domain/person/name"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Run("успех", func(t *testing.T) {
		tc := []struct {
			name           string
			firstname      string
			lastname       string
			middlename     string
			wantFirstname  string
			wantLastname   string
			wantMiddlename string
		}{
			{
				name:           "без отчества (пробелы по краям удаляются)",
				firstname:      " Иван ",
				lastname:       " Петров ",
				middlename:     "",
				wantFirstname:  "Иван",
				wantLastname:   "Петров",
				wantMiddlename: "",
			},
			{
				name:           "полное ФИО",
				firstname:      "Анна-Мария",
				lastname:       "Иванова",
				middlename:     "Сергеевна",
				wantFirstname:  "Анна-Мария",
				wantLastname:   "Иванова",
				wantMiddlename: "Сергеевна",
			},
			{
				name:           "с пробелами между словами (пробелы удаляются)",
				firstname:      "Анна   Мария",
				lastname:       "Иванова",
				middlename:     "Сергеевна",
				wantFirstname:  "Анна Мария",
				wantLastname:   "Иванова",
				wantMiddlename: "Сергеевна",
			},
			{
				name:           "при создании объекта регистр нормализуется",
				firstname:      "анна-мария",
				lastname:       "иванова",
				middlename:     "сергеевна",
				wantFirstname:  "Анна-Мария",
				wantLastname:   "Иванова",
				wantMiddlename: "Сергеевна",
			},
		}

		for _, tt := range tc {
			t.Run(tt.name, func(t *testing.T) {
				// Arrange
				n, err := name2.New(tt.firstname, tt.lastname, tt.middlename)

				// Assert
				assert.NoError(t, err)
				assert.Equal(t, tt.wantFirstname, n.FirstName())
				assert.Equal(t, tt.wantLastname, n.LastName())
				assert.Equal(t, tt.wantMiddlename, n.MiddleName())
			})
		}
	})

	t.Run("ошибка", func(t *testing.T) {
		tc := []struct {
			name       string
			firstname  string
			lastname   string
			middlename string
			wantErr    error
		}{
			{
				name:       "пустой firstname",
				firstname:  "",
				lastname:   "Иванова",
				middlename: "",
				wantErr:    name2.ErrInvalid,
			},
			{
				name:       "пустой lastname",
				firstname:  "Мария",
				lastname:   "",
				middlename: "",
				wantErr:    name2.ErrInvalid,
			},
			{
				name:       "firstname из пробелов",
				firstname:  "   \t\t\t\t\t ",
				lastname:   "Иванова",
				middlename: "",
				wantErr:    name2.ErrInvalid,
			},
			{
				name:       "lastname из пробелов",
				firstname:  "Мария",
				lastname:   "\t\t\t\t\t\t    ",
				middlename: "",
				wantErr:    name2.ErrInvalid,
			},
			{
				name:       "firstname с цифрами",
				firstname:  "Мария123",
				lastname:   "Иванова",
				middlename: "",
				wantErr:    name2.ErrInvalid,
			},
			{
				name:       "lastname с цифрами",
				firstname:  "Мария",
				lastname:   "Иванова123",
				middlename: "",
				wantErr:    name2.ErrInvalid,
			},
			{
				name:       "middlename с цифрами",
				firstname:  "Мария",
				lastname:   "Иванова",
				middlename: "Петро123вна",
				wantErr:    name2.ErrInvalid,
			},
			{
				name:       "firstname с недопустимыми символами",
				firstname:  "Мария...",
				lastname:   "Иванова",
				middlename: "",
				wantErr:    name2.ErrInvalid,
			},
			{
				name:       "lastname с недопустимыми символами",
				firstname:  "Мария",
				lastname:   "Иванова......[]]",
				middlename: "",
				wantErr:    name2.ErrInvalid,
			},
			{
				name:       "middlename с недопустимыми символами",
				firstname:  "Мария",
				lastname:   "Иванова",
				middlename: "Петровна[[[[",
				wantErr:    name2.ErrInvalid,
			},
			{
				name:       "firstname с превышением лимита по символам",
				firstname:  strings.Repeat("A", name2.PartCharsLimit+1),
				lastname:   "Иванова",
				middlename: "",
				wantErr:    name2.ErrInvalid,
			},
			{
				name:       "lastname с превышением лимита по символам",
				firstname:  "Мария",
				lastname:   strings.Repeat("A", name2.PartCharsLimit+1),
				middlename: "",
				wantErr:    name2.ErrInvalid,
			},
			{
				name:       "middlename с превышением лимита по символам",
				firstname:  "Мария",
				lastname:   "Иванова",
				middlename: strings.Repeat("A", name2.PartCharsLimit+1),
				wantErr:    name2.ErrInvalid,
			},
			{
				name:       "middlename с разделителем и пробелами вокруг него",
				firstname:  "Мария",
				lastname:   "Иванова",
				middlename: "Иванович - Петрович",
				wantErr:    name2.ErrInvalid,
			},
		}

		for _, tt := range tc {
			t.Run(tt.name, func(t *testing.T) {
				// Arrange
				_, err := name2.New(tt.firstname, tt.lastname, tt.middlename)

				// Assert
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
			})
		}
	})
}

func TestFullname(t *testing.T) {
	tc := []struct {
		name         string
		firstname    string
		lastname     string
		middlename   string
		wantFullname string
	}{
		{
			name:         "без отчества",
			firstname:    "Иван",
			lastname:     "Петров",
			middlename:   "",
			wantFullname: "Петров Иван",
		},
		{
			name:         "с отчеством",
			firstname:    "Иван",
			lastname:     "Петров",
			middlename:   "Петрович",
			wantFullname: "Петров Иван Петрович",
		},
		{
			name:         "с отчеством несколько слов",
			firstname:    "Иван",
			lastname:     "Петров",
			middlename:   "Петрович Иванович",
			wantFullname: "Петров Иван Петрович Иванович",
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			n, err := name2.New(tt.firstname, tt.lastname, tt.middlename)
			require.NoError(t, err)

			// Assert
			assert.Equal(t, tt.wantFullname, n.Fullname())
		})
	}
}

func TestWithInitials(t *testing.T) {
	tc := []struct {
		name             string
		firstname        string
		lastname         string
		middlename       string
		wantWithInitials string
	}{
		{
			name:             "без отчества",
			firstname:        "Иван",
			lastname:         "Петров",
			middlename:       "",
			wantWithInitials: "Петров И.",
		},
		{
			name:             "с отчеством",
			firstname:        "Иван",
			lastname:         "Петров",
			middlename:       "Петрович",
			wantWithInitials: "Петров И.П.",
		},
		{
			name:             "отчество в несколько слов",
			firstname:        "Иван",
			lastname:         "Петров",
			middlename:       "Петрович Иванович",
			wantWithInitials: "Петров И.П.И.",
		},
		{
			name:             "имя в несколько слов",
			firstname:        "Иван Сергей",
			lastname:         "Петров",
			middlename:       "Петрович",
			wantWithInitials: "Петров И.С.П.",
		},
		{
			name:             "имя в несколько слов с разделителем",
			firstname:        "Иван-Сергей",
			lastname:         "Петров",
			middlename:       "Петрович",
			wantWithInitials: "Петров И.-С.П.",
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			n, err := name2.New(tt.firstname, tt.lastname, tt.middlename)
			require.NoError(t, err)

			// Assert
			assert.Equal(t, tt.wantWithInitials, n.WithInitials())
		})
	}
}
