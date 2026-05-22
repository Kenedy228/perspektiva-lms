package commands_test

import (
	"context"
	"errors"
	"testing"

	orgports "gitflic.ru/lms/backend/internal/application/ports/organization"
	"gitflic.ru/lms/backend/internal/application/usecases/organization/commands"
	"gitflic.ru/lms/backend/internal/application/usecases/organization/common"
	domainorg "gitflic.ru/lms/backend/internal/domain/organization"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRenameUseCase(t *testing.T) {
	t.Run("некорректный идентификатор организации", func(t *testing.T) {
		// Arrange
		in := commands.RenameInput{
			ActorRole:      adminRole(),
			OrganizationID: "xxx",
			Name:           "ООО 'Ромашка'",
		}

		r := mockRepository{}
		a := mockAuditRecorder{}
		uc := commands.NewRenameUseCase(&r, &a)

		// Act
		out, err := uc.Execute(context.Background(), in)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, out)
		r.AssertExpectations(t)
	})

	t.Run("запрещает доступ не администратору", func(t *testing.T) {
		in := commands.RenameInput{
			ActorRole:      studentRole(),
			OrganizationID: uuid.NewString(),
			Name:           "ООО 'Ромашка'",
		}

		r := mockRepository{}
		a := mockAuditRecorder{}
		uc := commands.NewRenameUseCase(&r, &a)

		out, err := uc.Execute(context.Background(), in)

		assert.ErrorIs(t, err, common.ErrForbidden)
		assert.Nil(t, out)
		r.AssertNotCalled(t, "FindByID", mock.Anything, mock.Anything)
		r.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
		a.AssertNotCalled(t, "RecordOrganizationAudit", mock.Anything, mock.Anything)
	})

	t.Run("некорректный название организации", func(t *testing.T) {
		// Arrange
		in := commands.RenameInput{
			ActorRole:      adminRole(),
			OrganizationID: uuid.NewString(),
			Name:           "",
		}

		r := mockRepository{}
		a := mockAuditRecorder{}
		uc := commands.NewRenameUseCase(&r, &a)

		// Act
		out, err := uc.Execute(context.Background(), in)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, out)
		r.AssertExpectations(t)
	})

	t.Run("поиск организации завершился с ошибкой", func(t *testing.T) {
		// Arrange
		org := organizationFixture()
		in := commands.RenameInput{
			ActorRole:      adminRole(),
			OrganizationID: org.ID().String(),
			Name:           "ООО 'Ромашка'",
		}

		r := mockRepository{}
		a := mockAuditRecorder{}
		uc := commands.NewRenameUseCase(&r, &a)
		r.On("FindByID", mock.Anything, org.ID()).Return(nil, errors.New(""))

		// Act
		out, err := uc.Execute(context.Background(), in)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, out)
		r.AssertExpectations(t)
	})

	t.Run("сохранение изменений с ошибкой", func(t *testing.T) {
		// Arrange
		org := organizationFixture()
		in := commands.RenameInput{
			ActorRole:      adminRole(),
			OrganizationID: org.ID().String(),
			Name:           "ООО 'Ромашка'",
		}

		r := mockRepository{}
		a := mockAuditRecorder{}
		uc := commands.NewRenameUseCase(&r, &a)
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
		in := commands.RenameInput{
			ActorRole:      adminRole(),
			OrganizationID: org.ID().String(),
			Name:           "ООО 'Спартак'",
		}

		r := mockRepository{}
		a := mockAuditRecorder{}
		uc := commands.NewRenameUseCase(&r, &a)
		r.On("FindByID", mock.Anything, org.ID()).Return(org, nil)
		r.On("Save", mock.Anything, mock.MatchedBy(func(saved *domainorg.Organization) bool {
			return saved.ID() == org.ID() && saved.Name().Value() == "ООО 'Спартак'"
		})).Return(nil)
		a.On("RecordOrganizationAudit", mock.Anything, mock.MatchedBy(func(event orgports.AuditEvent) bool {
			return event.Action == orgports.AuditActionRename &&
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
