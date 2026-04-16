package organization

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateInn(t *testing.T) {
	tests := []struct {
		name    string
		inn     string
		wantErr error
	}{
		{name: "valid 10 digits (company)", inn: "1234567890", wantErr: nil},
		{name: "valid 12 digits (individual)", inn: "123456789012", wantErr: nil},
		{name: "error empty", inn: "", wantErr: ErrEmptyInn},
		{name: "error spaces", inn: "   ", wantErr: ErrEmptyInn},
		{name: "error 9 digits (too short)", inn: "123456789", wantErr: ErrInvalidInn},
		{name: "error 11 digits", inn: "12345678901", wantErr: ErrInvalidInn},
		{name: "error 13 digits (too long)", inn: "1234567890123", wantErr: ErrInvalidInn},
		{name: "error letters", inn: "123456789A", wantErr: ErrInvalidInn},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateInn(tt.inn)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestNew(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		params := Params{
			Name: "ООО Ромашка",
			Inn:  "7711223344",
		}
		org, err := New(params)
		require.NoError(t, err)
		require.NotNil(t, org)

		assert.NotEqual(t, uuid.Nil, org.Id())
		assert.Equal(t, "ООО Ромашка", org.Name())
		assert.Equal(t, "7711223344", org.Inn())

		assert.False(t, org.CreatedAt().IsZero())
		assert.False(t, org.UpdatedAt().IsZero())
	})

	t.Run("error invalid name", func(t *testing.T) {
		params := Params{
			Name: "",
			Inn:  "7711223344",
		}
		org, err := New(params)
		assert.ErrorIs(t, err, ErrEmptyName)
		assert.Nil(t, org)
	})
}

func TestOrganization_Mutators(t *testing.T) {
	params := Params{
		Name: "Старое Имя",
		Inn:  "7711223344",
	}
	org, _ := New(params)
	originalUpdatedAt := org.UpdatedAt()

	time.Sleep(time.Millisecond * 10)

	t.Run("Rename success", func(t *testing.T) {
		err := org.Rename("Новое Имя")
		require.NoError(t, err)
		assert.Equal(t, "Новое Имя", org.Name())
		assert.True(t, org.UpdatedAt().After(originalUpdatedAt))
	})

	t.Run("ChangeInn success", func(t *testing.T) {
		err := org.ChangeInn("112233445566")
		require.NoError(t, err)
		assert.Equal(t, "112233445566", org.Inn())
	})
}

func TestOrganization_Delete(t *testing.T) {
	params := Params{
		Name: "ООО Ромашка",
		Inn:  "7711223344",
	}
	org, _ := New(params)

	t.Run("initial state", func(t *testing.T) {
		assert.False(t, org.IsDeleted())
		assert.True(t, org.DeletedAt().IsZero())
	})

	t.Run("after delete", func(t *testing.T) {
		org.Delete()
		assert.True(t, org.IsDeleted())
		assert.False(t, org.DeletedAt().IsZero())
		assert.Equal(t, org.DeletedAt(), org.UpdatedAt())
	})

	t.Run("idempotent delete", func(t *testing.T) {
		firstDeletedAt := org.DeletedAt()
		time.Sleep(time.Millisecond * 10)

		org.Delete()
		assert.Equal(t, firstDeletedAt, org.DeletedAt())
	})
}
