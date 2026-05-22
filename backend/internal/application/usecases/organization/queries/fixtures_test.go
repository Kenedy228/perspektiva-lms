package queries_test

import (
	"gitflic.ru/lms/backend/internal/application/ports/organization"
	"gitflic.ru/lms/backend/internal/domain/account"
	"github.com/google/uuid"
)

func organizationShortSliceViewFixture() []organization.OrganizationShortView {
	view := organization.OrganizationShortView{
		ID:               uuid.NewString(),
		OrganizationName: "ООО 'Ромашка'",
	}

	return []organization.OrganizationShortView{view, view, view}
}

func organizationDetailedViewFixture() organization.OrganizationDetailedView {
	view := organization.OrganizationDetailedView{
		ID:               uuid.NewString(),
		INN:              "0000000000",
		INNTitle:         "организация",
		OrganizationName: "ООО 'Ромашка'",
	}

	return view
}

func adminRole() account.Role {
	return account.NewAdminRole()
}

func studentRole() account.Role {
	return account.NewStudentRole()
}
