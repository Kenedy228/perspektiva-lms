package commands_test

import (
	"context"
	"testing"

	accountports "gitflic.ru/lms/backend/internal/application/ports/account"
	"gitflic.ru/lms/backend/internal/application/usecases/account/commands"
	"gitflic.ru/lms/backend/internal/application/usecases/account/common"
	accountdomain "gitflic.ru/lms/backend/internal/domain/account"
	"gitflic.ru/lms/backend/internal/domain/role"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateUseCase_AdminUniquenessAndAudit(t *testing.T) {
	t.Run("forbidden for non-admin", func(t *testing.T) {
		r := mockRepository{}
		h := mockPasswordHasher{}
		a := mockAuditRecorder{}
		uc := commands.NewCreateUseCase(&r, &h, &a)

		out, err := uc.Execute(context.Background(), commands.CreateInput{
			ActorRole: role.NewStudent(),
			Login:     "stud2026",
			Password:  "secret",
			Role:      role.NewStudent(),
			PersonID:  uuid.NewString(),
		})

		assert.ErrorIs(t, err, common.ErrForbidden)
		assert.Nil(t, out)
		r.AssertNotCalled(t, "Save")
	})

	t.Run("rejects person with existing account", func(t *testing.T) {
		r := mockRepository{}
		h := mockPasswordHasher{}
		a := mockAuditRecorder{}
		uc := commands.NewCreateUseCase(&r, &h, &a)
		personID := uuid.New()

		r.On("LoginExists", mock.Anything, mock.Anything, uuid.Nil).Return(false, nil)
		r.On("PersonHasAccount", mock.Anything, personID, uuid.Nil).Return(true, nil)

		out, err := uc.Execute(context.Background(), commands.CreateInput{
			ActorRole: role.NewAdmin(),
			Login:     "stud2026",
			Password:  "secret",
			Role:      role.NewStudent(),
			PersonID:  personID.String(),
		})

		assert.ErrorIs(t, err, common.ErrConflict)
		assert.Nil(t, out)
		h.AssertNotCalled(t, "Hash")
	})

	t.Run("saves active account and audits", func(t *testing.T) {
		r := mockRepository{}
		h := mockPasswordHasher{}
		a := mockAuditRecorder{}
		uc := commands.NewCreateUseCase(&r, &h, &a)
		personID := uuid.New()

		r.On("LoginExists", mock.Anything, mock.Anything, uuid.Nil).Return(false, nil)
		r.On("PersonHasAccount", mock.Anything, personID, uuid.Nil).Return(false, nil)
		h.On("Hash", "secret").Return(hashFixture(), nil)
		r.On("Save", mock.Anything, mock.MatchedBy(func(acc *accountdomain.Account) bool {
			return acc.PersonID() == personID && acc.IsActive()
		})).Return(nil)
		a.On("RecordAccountAudit", mock.Anything, mock.MatchedBy(func(event accountports.AuditEvent) bool {
			return event.Action == accountports.AuditActionCreate && event.ActorRole == role.TypeAdmin.String()
		})).Return(nil)

		out, err := uc.Execute(context.Background(), commands.CreateInput{
			ActorRole: role.NewAdmin(),
			Login:     "stud2026",
			Password:  "secret",
			Role:      role.NewStudent(),
			PersonID:  personID.String(),
		})

		require.NoError(t, err)
		assert.NotEmpty(t, out.ID)
	})
}

func TestChangeUseCases_AdminOnly(t *testing.T) {
	t.Run("change login checks uniqueness", func(t *testing.T) {
		r := mockRepository{}
		a := mockAuditRecorder{}
		uc := commands.NewChangeLoginUseCase(&r, &a)
		acc := accountFixture()

		r.On("LoginExists", mock.Anything, mock.Anything, acc.ID()).Return(false, nil)
		r.On("FindByID", mock.Anything, acc.ID()).Return(acc, nil)
		r.On("Save", mock.Anything, acc).Return(nil)
		a.On("RecordAccountAudit", mock.Anything, mock.MatchedBy(func(event accountports.AuditEvent) bool {
			return event.Action == accountports.AuditActionChangeLogin
		})).Return(nil)

		out, err := uc.Execute(context.Background(), commands.ChangeLoginInput{
			ActorRole: role.NewAdmin(),
			AccountID: acc.ID().String(),
			Login:     "newlogin",
		})

		require.NoError(t, err)
		assert.Equal(t, acc.ID().String(), out.ID)
	})

	t.Run("change password hashes by admin", func(t *testing.T) {
		r := mockRepository{}
		h := mockPasswordHasher{}
		a := mockAuditRecorder{}
		uc := commands.NewChangePasswordUseCase(&r, &h, &a)
		acc := accountFixture()

		r.On("FindByID", mock.Anything, acc.ID()).Return(acc, nil)
		h.On("Hash", "new-secret").Return(hashFixture(), nil)
		r.On("Save", mock.Anything, acc).Return(nil)
		a.On("RecordAccountAudit", mock.Anything, mock.MatchedBy(func(event accountports.AuditEvent) bool {
			return event.Action == accountports.AuditActionChangePassword
		})).Return(nil)

		out, err := uc.Execute(context.Background(), commands.ChangePasswordInput{
			ActorRole: role.NewAdmin(),
			AccountID: acc.ID().String(),
			Password:  "new-secret",
		})

		require.NoError(t, err)
		assert.Equal(t, acc.ID().String(), out.ID)
	})

	t.Run("change role", func(t *testing.T) {
		r := mockRepository{}
		a := mockAuditRecorder{}
		uc := commands.NewChangeRoleUseCase(&r, &a)
		acc := accountFixture()

		r.On("FindByID", mock.Anything, acc.ID()).Return(acc, nil)
		r.On("Save", mock.Anything, acc).Return(nil)
		a.On("RecordAccountAudit", mock.Anything, mock.MatchedBy(func(event accountports.AuditEvent) bool {
			return event.Action == accountports.AuditActionChangeRole
		})).Return(nil)

		out, err := uc.Execute(context.Background(), commands.ChangeRoleInput{
			ActorRole:   role.NewAdmin(),
			AccountID:   acc.ID().String(),
			AccountRole: role.NewCreator(),
		})

		require.NoError(t, err)
		assert.Equal(t, acc.ID().String(), out.ID)
	})
}

func TestLifecycleUseCases(t *testing.T) {
	t.Run("block and activate", func(t *testing.T) {
		r := mockRepository{}
		a := mockAuditRecorder{}
		acc := accountFixture()

		block := commands.NewBlockUseCase(&r, &a)
		r.On("FindByID", mock.Anything, acc.ID()).Return(acc, nil).Once()
		r.On("Save", mock.Anything, acc).Return(nil).Once()
		a.On("RecordAccountAudit", mock.Anything, mock.MatchedBy(func(event accountports.AuditEvent) bool {
			return event.Action == accountports.AuditActionBlock
		})).Return(nil).Once()

		out, err := block.Execute(context.Background(), commands.BlockInput{
			ActorRole: role.NewAdmin(),
			AccountID: acc.ID().String(),
		})
		require.NoError(t, err)
		assert.Equal(t, acc.ID().String(), out.ID)
		assert.True(t, acc.IsBlocked())

		activate := commands.NewActivateUseCase(&r, &a)
		r.On("FindByID", mock.Anything, acc.ID()).Return(acc, nil).Once()
		r.On("Save", mock.Anything, acc).Return(nil).Once()
		a.On("RecordAccountAudit", mock.Anything, mock.MatchedBy(func(event accountports.AuditEvent) bool {
			return event.Action == accountports.AuditActionActivate
		})).Return(nil).Once()

		_, err = activate.Execute(context.Background(), commands.ActivateInput{
			ActorRole: role.NewAdmin(),
			AccountID: acc.ID().String(),
		})
		require.NoError(t, err)
		assert.True(t, acc.IsActive())
	})

	t.Run("delete is soft lifecycle delete", func(t *testing.T) {
		r := mockRepository{}
		a := mockAuditRecorder{}
		uc := commands.NewDeleteUseCase(&r, &a)
		acc := accountFixture()

		r.On("FindByID", mock.Anything, acc.ID()).Return(acc, nil)
		r.On("Save", mock.Anything, acc).Return(nil)
		a.On("RecordAccountAudit", mock.Anything, mock.MatchedBy(func(event accountports.AuditEvent) bool {
			return event.Action == accountports.AuditActionDelete
		})).Return(nil)

		err := uc.Execute(context.Background(), commands.DeleteInput{
			ActorRole: role.NewAdmin(),
			AccountID: acc.ID().String(),
		})

		require.NoError(t, err)
		assert.True(t, acc.IsDeleted())
	})
}

func TestAuthenticateUseCase(t *testing.T) {
	t.Run("active account with matching password succeeds", func(t *testing.T) {
		r := mockRepository{}
		c := mockPasswordComparer{}
		uc := commands.NewAuthenticateUseCase(&r, &c)
		acc := accountFixture()

		r.On("FindByLogin", mock.Anything, acc.Login()).Return(acc, nil)
		c.On("Compare", acc.PasswordHash(), "secret").Return(true)

		out, err := uc.Execute(context.Background(), commands.AuthenticateInput{
			Login:    acc.Login().Value(),
			Password: "secret",
		})

		require.NoError(t, err)
		assert.Equal(t, acc.ID().String(), out.AccountID)
	})

	t.Run("blocked account is rejected without password compare", func(t *testing.T) {
		r := mockRepository{}
		c := mockPasswordComparer{}
		uc := commands.NewAuthenticateUseCase(&r, &c)
		acc := accountFixture()
		require.NoError(t, acc.Block())

		r.On("FindByLogin", mock.Anything, acc.Login()).Return(acc, nil)

		out, err := uc.Execute(context.Background(), commands.AuthenticateInput{
			Login:    acc.Login().Value(),
			Password: "secret",
		})

		assert.ErrorIs(t, err, common.ErrInvalidCredentials)
		assert.Nil(t, out)
		c.AssertNotCalled(t, "Compare")
	})
}

func TestConstructorsPanicOnNilDependencies(t *testing.T) {
	assert.Panics(t, func() { commands.NewCreateUseCase(nil, &mockPasswordHasher{}, &mockAuditRecorder{}) })
	assert.Panics(t, func() { commands.NewChangePasswordUseCase(&mockRepository{}, nil, &mockAuditRecorder{}) })
	assert.Panics(t, func() { commands.NewBlockUseCase(&mockRepository{}, nil) })
	assert.Panics(t, func() { commands.NewAuthenticateUseCase(nil, &mockPasswordComparer{}) })
}
