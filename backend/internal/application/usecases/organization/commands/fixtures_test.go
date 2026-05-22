package commands_test

import (
	"gitflic.ru/lms/backend/internal/domain/account"
	"gitflic.ru/lms/backend/internal/domain/organization"
	inn2 "gitflic.ru/lms/backend/internal/domain/organization/inn"
	"gitflic.ru/lms/backend/internal/domain/organization/name"
)

func organizationFixture() *organization.Organization {
	i, _ := inn2.New("1030000000", inn2.TypeLegalEntity)
	n, _ := name.New("ООО 'Ромашка'")
	org, _ := organization.New(i, n)
	return org
}

func adminRole() account.Role {
	return account.NewAdminRole()
}

func studentRole() account.Role {
	return account.NewStudentRole()
}
