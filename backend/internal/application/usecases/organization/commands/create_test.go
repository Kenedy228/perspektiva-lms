package commands_test

import (
	"context"
	"errors"
	"testing"

	orgports "gitflic.ru/lms/backend/internal/application/ports/organization"
	"gitflic.ru/lms/backend/internal/application/usecases/organization/commands"
	"gitflic.ru/lms/backend/internal/application/usecases/organization/common"
	"gitflic.ru/lms/backend/internal/domain/organization/inn"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateUseCase(t *testing.T) {
	t.Run("некорректный ИНН", func(t *testing.T) {
		// Arrange
		in := commands.CreateInput{
			ActorRole: adminRole(),
			INN:       "xxxx",
			INNType:   "org",
			Name:      "ООО 'Ромашка'",
		}

		r := mockRepository{}
		a := mockAuditRecorder{}
		uc := commands.NewCreateUseCase(&r, &a)

		// Act
		out, err := uc.Execute(context.Background(), in)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, out)
		r.AssertExpectations(t)
	})

	t.Run("запрещает доступ не администратору", func(t *testing.T) {
		in := commands.CreateInput{
			ActorRole: studentRole(),
			INN:       "1030000000",
			INNType:   inn.TypeLegalEntity.String(),
			Name:      "ООО 'Ромашка'",
		}

		r := mockRepository{}
		a := mockAuditRecorder{}
		uc := commands.NewCreateUseCase(&r, &a)

		out, err := uc.Execute(context.Background(), in)

		assert.ErrorIs(t, err, common.ErrForbidden)
		assert.Nil(t, out)
		r.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
		a.AssertNotCalled(t, "RecordOrganizationAudit", mock.Anything, mock.Anything)
	})

	t.Run("некорректное название организации", func(t *testing.T) {
		// Arrange
		in := commands.CreateInput{
			ActorRole: adminRole(),
			INN:       "1030000000",
			INNType:   inn.TypeLegalEntity.String(),
			Name:      "",
		}

		r := mockRepository{}
		a := mockAuditRecorder{}
		uc := commands.NewCreateUseCase(&r, &a)

		// Act
		out, err := uc.Execute(context.Background(), in)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, out)
		r.AssertExpectations(t)
	})

	t.Run("ошибка операции сохранения", func(t *testing.T) {
		// Arrange
		in := commands.CreateInput{
			ActorRole: adminRole(),
			INN:       "1030000000",
			INNType:   inn.TypeLegalEntity.String(),
			Name:      "ООО 'Ромашка'",
		}

		r := mockRepository{}
		a := mockAuditRecorder{}
		uc := commands.NewCreateUseCase(&r, &a)
		r.On("Save", mock.Anything, mock.AnythingOfType("*organization.Organization")).Return(errors.New(""))

		// Act
		out, err := uc.Execute(context.Background(), in)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, out)
		r.AssertExpectations(t)
	})

	t.Run("успех", func(t *testing.T) {
		// Arrange
		in := commands.CreateInput{
			ActorRole: adminRole(),
			INN:       "1030000000",
			INNType:   inn.TypeLegalEntity.String(),
			Name:      "ООО 'Ромашка'",
		}

		r := mockRepository{}
		a := mockAuditRecorder{}
		uc := commands.NewCreateUseCase(&r, &a)
		r.On("Save", mock.Anything, mock.AnythingOfType("*organization.Organization")).Return(nil)
		a.On("RecordOrganizationAudit", mock.Anything, mock.MatchedBy(func(event orgports.AuditEvent) bool {
			return event.Action == orgports.AuditActionCreate &&
				event.OrganizationID != "" &&
				event.ActorRole == adminRole().Kind().String()
		})).Return(nil)

		// Act
		out, err := uc.Execute(context.Background(), in)

		// Assert
		assert.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, out.ID)
		r.AssertExpectations(t)
		a.AssertExpectations(t)
	})
}
