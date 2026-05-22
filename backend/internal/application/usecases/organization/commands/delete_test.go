package commands_test

import (
	"context"
	"errors"
	"testing"

	orgports "gitflic.ru/lms/backend/internal/application/ports/organization"
	"gitflic.ru/lms/backend/internal/application/usecases/organization/commands"
	"gitflic.ru/lms/backend/internal/application/usecases/organization/common"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeleteUseCase(t *testing.T) {
	t.Run("некорректный идентификатор организации", func(t *testing.T) {
		// Arrange
		in := commands.DeleteByIDInput{
			ActorRole:      adminRole(),
			OrganizationID: "xxx",
		}

		r := mockRepository{}
		a := mockAuditRecorder{}
		uc := commands.NewDeleteByIDUseCase(&r, &a)

		// Act
		err := uc.Execute(context.Background(), in)

		// Assert
		assert.Error(t, err)
		r.AssertExpectations(t)
	})

	t.Run("запрещает доступ не администратору", func(t *testing.T) {
		id := uuid.New()
		in := commands.DeleteByIDInput{
			ActorRole:      studentRole(),
			OrganizationID: id.String(),
		}

		r := mockRepository{}
		a := mockAuditRecorder{}
		uc := commands.NewDeleteByIDUseCase(&r, &a)

		err := uc.Execute(context.Background(), in)

		assert.ErrorIs(t, err, common.ErrForbidden)
		r.AssertNotCalled(t, "DeleteByID", mock.Anything, mock.Anything)
		a.AssertNotCalled(t, "RecordOrganizationAudit", mock.Anything, mock.Anything)
	})

	t.Run("операция удаления с ошибкой", func(t *testing.T) {
		// Arrange
		org := organizationFixture()
		in := commands.DeleteByIDInput{
			ActorRole:      adminRole(),
			OrganizationID: org.ID().String(),
		}

		r := mockRepository{}
		a := mockAuditRecorder{}
		uc := commands.NewDeleteByIDUseCase(&r, &a)
		r.On("FindByID", mock.Anything, org.ID()).Return(org, nil)
		r.On("DeleteByID", mock.Anything, org.ID()).Return(errors.New(""))

		// Act
		err := uc.Execute(context.Background(), in)

		// Assert
		assert.Error(t, err)
		r.AssertExpectations(t)
	})

	t.Run("поиск организации завершился с ошибкой", func(t *testing.T) {
		id := uuid.New()
		in := commands.DeleteByIDInput{
			ActorRole:      adminRole(),
			OrganizationID: id.String(),
		}

		r := mockRepository{}
		a := mockAuditRecorder{}
		uc := commands.NewDeleteByIDUseCase(&r, &a)
		r.On("FindByID", mock.Anything, id).Return(nil, errors.New(""))

		err := uc.Execute(context.Background(), in)

		assert.Error(t, err)
		r.AssertExpectations(t)
		r.AssertNotCalled(t, "DeleteByID", mock.Anything, mock.Anything)
		a.AssertNotCalled(t, "RecordOrganizationAudit", mock.Anything, mock.Anything)
	})

	t.Run("успех", func(t *testing.T) {
		// Arrange
		org := organizationFixture()
		in := commands.DeleteByIDInput{
			ActorRole:      adminRole(),
			OrganizationID: org.ID().String(),
		}

		r := mockRepository{}
		a := mockAuditRecorder{}
		uc := commands.NewDeleteByIDUseCase(&r, &a)
		r.On("FindByID", mock.Anything, org.ID()).Return(org, nil)
		r.On("DeleteByID", mock.Anything, org.ID()).Return(nil)
		a.On("RecordOrganizationAudit", mock.Anything, mock.MatchedBy(func(event orgports.AuditEvent) bool {
			return event.Action == orgports.AuditActionDelete &&
				event.OrganizationID == org.ID().String() &&
				event.ActorRole == adminRole().Kind().String()
		})).Return(nil)

		// Act
		err := uc.Execute(context.Background(), in)

		// Assert
		assert.NoError(t, err)
		r.AssertExpectations(t)
		a.AssertExpectations(t)
	})
}
