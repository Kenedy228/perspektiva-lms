package commands_test

import (
	"context"
	"errors"
	"testing"

	orgports "gitflic.ru/lms/backend/internal/application/ports/organization"
	"gitflic.ru/lms/backend/internal/application/usecases/organization/commands"
	"gitflic.ru/lms/backend/internal/application/usecases/organization/common"
	domainorg "gitflic.ru/lms/backend/internal/domain/organization"
	"gitflic.ru/lms/backend/internal/domain/organization/inn"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestChangeINNUseCase(t *testing.T) {
	t.Run("некорректный идентификатор организации", func(t *testing.T) {
		// Arrange
		in := commands.ChangeINNInput{
			ActorRole:      adminRole(),
			OrganizationID: "xxx",
			INN:            "1030000000",
			INNType:        inn.TypeLegalEntity.String(),
		}

		r := mockRepository{}
		a := mockAuditRecorder{}
		uc := commands.NewChangeINNUseCase(&r, &a)

		// Act
		out, err := uc.Execute(context.Background(), in)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, out)
		r.AssertExpectations(t)
	})

	t.Run("запрещает доступ не администратору", func(t *testing.T) {
		in := commands.ChangeINNInput{
			ActorRole:      studentRole(),
			OrganizationID: uuid.NewString(),
			INN:            "1030000000",
			INNType:        inn.TypeLegalEntity.String(),
		}

		r := mockRepository{}
		a := mockAuditRecorder{}
		uc := commands.NewChangeINNUseCase(&r, &a)

		out, err := uc.Execute(context.Background(), in)

		assert.ErrorIs(t, err, common.ErrForbidden)
		assert.Nil(t, out)
		r.AssertNotCalled(t, "FindByID", mock.Anything, mock.Anything)
		r.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
		a.AssertNotCalled(t, "RecordOrganizationAudit", mock.Anything, mock.Anything)
	})

	t.Run("некорректный ИНН", func(t *testing.T) {
		// Arrange
		in := commands.ChangeINNInput{
			ActorRole:      adminRole(),
			OrganizationID: uuid.NewString(),
			INN:            "xxx",
			INNType:        "org",
		}

		r := mockRepository{}
		a := mockAuditRecorder{}
		uc := commands.NewChangeINNUseCase(&r, &a)

		// Act
		out, err := uc.Execute(context.Background(), in)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, out)
		r.AssertExpectations(t)
	})

	t.Run("ошибка поиска записи организации", func(t *testing.T) {
		// Arrange
		id := uuid.New()
		in := commands.ChangeINNInput{
			ActorRole:      adminRole(),
			OrganizationID: id.String(),
			INN:            "1030000000",
			INNType:        inn.TypeLegalEntity.String(),
		}

		r := mockRepository{}
		a := mockAuditRecorder{}
		uc := commands.NewChangeINNUseCase(&r, &a)
		r.On("FindByID", mock.Anything, id).Return(nil, errors.New(""))

		// Act
		out, err := uc.Execute(context.Background(), in)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, out)
		r.AssertExpectations(t)
	})

	t.Run("ошибка сохранения изменений", func(t *testing.T) {
		// Arrange
		org := organizationFixture()
		in := commands.ChangeINNInput{
			ActorRole:      adminRole(),
			OrganizationID: org.ID().String(),
			INN:            "1030000000",
			INNType:        inn.TypeLegalEntity.String(),
		}

		r := mockRepository{}
		a := mockAuditRecorder{}
		uc := commands.NewChangeINNUseCase(&r, &a)
		r.On("FindByID", mock.Anything, org.ID()).Return(org, nil)
		r.On("Save", mock.Anything, org).Return(errors.New(""))

		// Act
		out, err := uc.Execute(context.Background(), in)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, out)
		r.AssertExpectations(t)
	})

	t.Run("успех", func(t *testing.T) {
		// Arrange
		org := organizationFixture()
		in := commands.ChangeINNInput{
			ActorRole:      adminRole(),
			OrganizationID: org.ID().String(),
			INN:            "1030000000",
			INNType:        inn.TypeLegalEntity.String(),
		}

		r := mockRepository{}
		a := mockAuditRecorder{}
		uc := commands.NewChangeINNUseCase(&r, &a)
		r.On("FindByID", mock.Anything, org.ID()).Return(org, nil)
		r.On("Save", mock.Anything, mock.MatchedBy(func(saved *domainorg.Organization) bool {
			return saved.ID() == org.ID() && saved.INN().Value() == "1030000000"
		})).Return(nil)
		a.On("RecordOrganizationAudit", mock.Anything, mock.MatchedBy(func(event orgports.AuditEvent) bool {
			return event.Action == orgports.AuditActionChangeINN &&
				event.OrganizationID == org.ID().String() &&
				event.ActorRole == adminRole().Kind().String()
		})).Return(nil)

		// Act
		out, err := uc.Execute(context.Background(), in)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, org.ID().String(), out.ID)
		r.AssertExpectations(t)
		a.AssertExpectations(t)
	})
}
