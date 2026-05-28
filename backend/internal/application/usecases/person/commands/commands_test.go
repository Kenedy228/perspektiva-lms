package commands_test

import (
	"context"
	"errors"
	"testing"

	personports "gitflic.ru/lms/backend/internal/application/ports/person"
	"gitflic.ru/lms/backend/internal/application/usecases/person/commands"
	"gitflic.ru/lms/backend/internal/application/usecases/person/common"
	"gitflic.ru/lms/backend/internal/domain/role"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateUseCase_AdminOnlyAndAudit(t *testing.T) {
	t.Run("forbidden for non-admin", func(t *testing.T) {
		r := mockRepository{}
		a := mockAuditRecorder{}
		uc := commands.NewCreateUseCase(&r, &a)

		out, err := uc.Execute(context.Background(), commands.CreateInput{
			ActorRole: role.NewStudent(),
			FirstName: "Иван",
			LastName:  "Иванов",
		})

		assert.ErrorIs(t, err, common.ErrForbidden)
		assert.Nil(t, out)
		r.AssertNotCalled(t, "Save")
	})

	t.Run("saves and audits", func(t *testing.T) {
		r := mockRepository{}
		a := mockAuditRecorder{}
		uc := commands.NewCreateUseCase(&r, &a)

		r.On("Save", mock.Anything, mock.Anything).Return(nil)
		a.On("RecordPersonAudit", mock.Anything, mock.MatchedBy(func(event personports.AuditEvent) bool {
			return event.Action == personports.AuditActionCreate && event.ActorRole == role.TypeAdmin.String()
		})).Return(nil)

		out, err := uc.Execute(context.Background(), commands.CreateInput{
			ActorRole: role.NewAdmin(),
			FirstName: "Иван",
			LastName:  "Иванов",
		})

		require.NoError(t, err)
		assert.NotEmpty(t, out.ID)
		r.AssertExpectations(t)
		a.AssertExpectations(t)
	})
}

func TestCreateWithProfileUseCase_ChecksSNILSAndAudits(t *testing.T) {
	r := mockRepository{}
	a := mockAuditRecorder{}
	uc := commands.NewCreateWithProfileUseCase(&r, &a)
	orgID := uuid.NewString()

	r.On("SNILSExists", mock.Anything, mock.Anything, uuid.Nil).Return(false, nil)
	r.On("Save", mock.Anything, mock.Anything).Return(nil)
	a.On("RecordPersonAudit", mock.Anything, mock.MatchedBy(func(event personports.AuditEvent) bool {
		return event.Action == personports.AuditActionCreateWithProfile && event.OrganizationID == orgID
	})).Return(nil)

	out, err := uc.Execute(context.Background(), commands.CreateWithProfileInput{
		ActorRole:      role.NewAdmin(),
		DateOfBirth:    validDOB(),
		FirstName:      "Иван",
		LastName:       "Иванов",
		Snils:          validSNILS(),
		JobTitle:       "инженер",
		Education:      "высшее",
		OrganizationID: orgID,
	})

	require.NoError(t, err)
	assert.NotEmpty(t, out.ID)
	r.AssertExpectations(t)
	a.AssertExpectations(t)
}

func TestAttachAndReplaceProfileUseCases(t *testing.T) {
	t.Run("attach rejects duplicate snils", func(t *testing.T) {
		r := mockRepository{}
		a := mockAuditRecorder{}
		uc := commands.NewAttachProfileUseCase(&r, &a)
		personID := uuid.New()

		r.On("SNILSExists", mock.Anything, mock.Anything, personID).Return(true, nil)

		out, err := uc.Execute(context.Background(), commands.AttachProfileInput{
			ActorRole:   role.NewAdmin(),
			DateOfBirth: validDOB(),
			PersonID:    personID.String(),
			Snils:       validSNILS(),
			Education:   "высшее",
		})

		assert.ErrorIs(t, err, common.ErrConflict)
		assert.Nil(t, out)
		r.AssertNotCalled(t, "FindByID")
	})

	t.Run("replace requires existing profile", func(t *testing.T) {
		r := mockRepository{}
		a := mockAuditRecorder{}
		uc := commands.NewReplaceProfileUseCase(&r, &a)
		p := personFixture()

		r.On("SNILSExists", mock.Anything, mock.Anything, p.ID()).Return(false, nil)
		r.On("FindByID", mock.Anything, p.ID()).Return(p, nil)

		out, err := uc.Execute(context.Background(), commands.ReplaceProfileInput{
			ActorRole:   role.NewAdmin(),
			DateOfBirth: validDOB(),
			PersonID:    p.ID().String(),
			Snils:       validSNILS(),
		})

		assert.Error(t, err)
		assert.Nil(t, out)
		r.AssertNotCalled(t, "Save")
	})
}

func TestRenameUseCase_SavesAndAudits(t *testing.T) {
	r := mockRepository{}
	a := mockAuditRecorder{}
	uc := commands.NewRenameUseCase(&r, &a)
	p := personFixture()

	r.On("FindByID", mock.Anything, p.ID()).Return(p, nil)
	r.On("Save", mock.Anything, p).Return(nil)
	a.On("RecordPersonAudit", mock.Anything, mock.MatchedBy(func(event personports.AuditEvent) bool {
		return event.Action == personports.AuditActionRename
	})).Return(nil)

	out, err := uc.Execute(context.Background(), commands.RenameInput{
		ActorRole:  role.NewAdmin(),
		PersonID:   p.ID().String(),
		FirstName:  "Петр",
		LastName:   "Петров",
		MiddleName: "Петрович",
	})

	require.NoError(t, err)
	assert.Equal(t, p.ID().String(), out.ID)
}

func TestProfileFieldUseCases(t *testing.T) {
	t.Run("change snils saves and audits", func(t *testing.T) {
		r := mockRepository{}
		a := mockAuditRecorder{}
		uc := commands.NewChangeSNILSUseCase(&r, &a)
		p := personWithProfileFixture(uuid.Nil)

		r.On("SNILSExists", mock.Anything, mock.Anything, p.ID()).Return(false, nil)
		r.On("FindByID", mock.Anything, p.ID()).Return(p, nil)
		r.On("Save", mock.Anything, p).Return(nil)
		a.On("RecordPersonAudit", mock.Anything, mock.MatchedBy(func(event personports.AuditEvent) bool {
			return event.Action == personports.AuditActionChangeSNILS
		})).Return(nil)

		out, err := uc.Execute(context.Background(), commands.ChangeSNILSInput{
			ActorRole: role.NewAdmin(),
			PersonID:  p.ID().String(),
			Snils:     validSNILS(),
		})

		require.NoError(t, err)
		assert.Equal(t, p.ID().String(), out.ID)
	})

	t.Run("change job title requires profile", func(t *testing.T) {
		r := mockRepository{}
		a := mockAuditRecorder{}
		uc := commands.NewChangeJobTitleUseCase(&r, &a)
		p := personFixture()

		r.On("FindByID", mock.Anything, p.ID()).Return(p, nil)

		out, err := uc.Execute(context.Background(), commands.ChangeJobTitleInput{
			ActorRole: role.NewAdmin(),
			PersonID:  p.ID().String(),
			JobTitle:  "руководитель",
		})

		assert.Error(t, err)
		assert.Nil(t, out)
		r.AssertNotCalled(t, "Save")
	})

	t.Run("change date of birth saves and audits", func(t *testing.T) {
		r := mockRepository{}
		a := mockAuditRecorder{}
		uc := commands.NewChangeDateOfBirthUseCase(&r, &a)
		p := personWithProfileFixture(uuid.Nil)

		r.On("FindByID", mock.Anything, p.ID()).Return(p, nil)
		r.On("Save", mock.Anything, p).Return(nil)
		a.On("RecordPersonAudit", mock.Anything, mock.MatchedBy(func(event personports.AuditEvent) bool {
			return event.Action == personports.AuditActionChangeDateOfBirth
		})).Return(nil)

		out, err := uc.Execute(context.Background(), commands.ChangeDateOfBirthInput{
			ActorRole:   role.NewAdmin(),
			PersonID:    p.ID().String(),
			DateOfBirth: validDOB(),
		})

		require.NoError(t, err)
		assert.Equal(t, p.ID().String(), out.ID)
	})

	t.Run("change education saves and audits", func(t *testing.T) {
		r := mockRepository{}
		a := mockAuditRecorder{}
		uc := commands.NewChangeEducationUseCase(&r, &a)
		p := personWithProfileFixture(uuid.Nil)

		r.On("FindByID", mock.Anything, p.ID()).Return(p, nil)
		r.On("Save", mock.Anything, p).Return(nil)
		a.On("RecordPersonAudit", mock.Anything, mock.MatchedBy(func(event personports.AuditEvent) bool {
			return event.Action == personports.AuditActionChangeEducation
		})).Return(nil)

		out, err := uc.Execute(context.Background(), commands.ChangeEducationInput{
			ActorRole: role.NewAdmin(),
			PersonID:  p.ID().String(),
			Education: "аспирантура",
		})

		require.NoError(t, err)
		assert.Equal(t, p.ID().String(), out.ID)
	})
}

func TestOrganizationAssignmentUseCases(t *testing.T) {
	t.Run("assign organization", func(t *testing.T) {
		r := mockRepository{}
		a := mockAuditRecorder{}
		uc := commands.NewAssignOrganizationUseCase(&r, &a)
		p := personWithProfileFixture(uuid.Nil)
		orgID := uuid.New()

		r.On("FindByID", mock.Anything, p.ID()).Return(p, nil)
		r.On("Save", mock.Anything, p).Return(nil)
		a.On("RecordPersonAudit", mock.Anything, mock.MatchedBy(func(event personports.AuditEvent) bool {
			return event.Action == personports.AuditActionAssignOrganization && event.OrganizationID == orgID.String()
		})).Return(nil)

		out, err := uc.Execute(context.Background(), commands.AssignOrganizationInput{
			ActorRole:      role.NewAdmin(),
			PersonID:       p.ID().String(),
			OrganizationID: orgID.String(),
		})

		require.NoError(t, err)
		assert.Equal(t, p.ID().String(), out.ID)
	})

	t.Run("remove organization", func(t *testing.T) {
		r := mockRepository{}
		a := mockAuditRecorder{}
		uc := commands.NewRemoveOrganizationUseCase(&r, &a)
		p := personWithProfileFixture(uuid.New())

		r.On("FindByID", mock.Anything, p.ID()).Return(p, nil)
		r.On("Save", mock.Anything, p).Return(nil)
		a.On("RecordPersonAudit", mock.Anything, mock.MatchedBy(func(event personports.AuditEvent) bool {
			return event.Action == personports.AuditActionRemoveOrganization
		})).Return(nil)

		out, err := uc.Execute(context.Background(), commands.RemoveOrganizationInput{
			ActorRole: role.NewAdmin(),
			PersonID:  p.ID().String(),
		})

		require.NoError(t, err)
		assert.Equal(t, p.ID().String(), out.ID)
	})
}

func TestDetachProfileUseCase_SavesOnlyWhenProfileExists(t *testing.T) {
	r := mockRepository{}
	a := mockAuditRecorder{}
	uc := commands.NewDetachProfileUseCase(&r, &a)
	p := personWithProfileFixture(uuid.Nil)

	r.On("FindByID", mock.Anything, p.ID()).Return(p, nil)
	r.On("Save", mock.Anything, p).Return(nil)
	a.On("RecordPersonAudit", mock.Anything, mock.MatchedBy(func(event personports.AuditEvent) bool {
		return event.Action == personports.AuditActionDetachProfile
	})).Return(nil)

	out, err := uc.Execute(context.Background(), commands.DetachProfileInput{
		ActorRole: role.NewAdmin(),
		PersonID:  p.ID().String(),
	})

	require.NoError(t, err)
	assert.Equal(t, p.ID().String(), out.PersonID)
}

func TestDeleteByIDUseCase_LoadsBeforeDeleteAndAudits(t *testing.T) {
	r := mockRepository{}
	a := mockAuditRecorder{}
	uc := commands.NewDeleteByIDUseCase(&r, &a)
	p := personFixture()

	r.On("FindByID", mock.Anything, p.ID()).Return(p, nil)
	r.On("DeleteByID", mock.Anything, p.ID()).Return(nil)
	a.On("RecordPersonAudit", mock.Anything, mock.MatchedBy(func(event personports.AuditEvent) bool {
		return event.Action == personports.AuditActionDelete
	})).Return(nil)

	err := uc.Execute(context.Background(), commands.DeleteByIDInput{
		ActorRole: role.NewAdmin(),
		PersonID:  p.ID().String(),
	})

	require.NoError(t, err)
	r.AssertExpectations(t)
	a.AssertExpectations(t)
}

func TestUseCaseConstructorsPanicOnNilDependencies(t *testing.T) {
	assert.Panics(t, func() { commands.NewCreateUseCase(nil, &mockAuditRecorder{}) })
	assert.Panics(t, func() { commands.NewCreateUseCase(&mockRepository{}, nil) })
}

func TestCommandWrapsRepositoryErrors(t *testing.T) {
	r := mockRepository{}
	a := mockAuditRecorder{}
	uc := commands.NewDeleteByIDUseCase(&r, &a)
	p := personFixture()

	r.On("FindByID", mock.Anything, p.ID()).Return(p, errors.New("storage down"))

	err := uc.Execute(context.Background(), commands.DeleteByIDInput{
		ActorRole: role.NewAdmin(),
		PersonID:  p.ID().String(),
	})

	assert.ErrorContains(t, err, "find person")
}
