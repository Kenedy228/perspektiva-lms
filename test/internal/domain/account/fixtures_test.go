package account_test

import (
	"gitflic.ru/lms/internal/domain/account/login"
	"gitflic.ru/lms/internal/domain/account/passhash"
	"gitflic.ru/lms/internal/domain/role"
)

func loginFixture() login.Login {
	l, _ := login.New("student2026")
	return l
}

func hashFixture() passhash.Hash {
	h, _ := passhash.New("$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy")
	return h
}

func roleFixture() role.Role {
	return role.NewAdmin()
}
