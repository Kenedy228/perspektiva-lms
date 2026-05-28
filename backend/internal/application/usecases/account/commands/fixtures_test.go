package commands_test

import (
	"gitflic.ru/lms/backend/internal/domain/account"
	"gitflic.ru/lms/backend/internal/domain/account/login"
	"gitflic.ru/lms/backend/internal/domain/account/passhash"
	"gitflic.ru/lms/backend/internal/domain/role"
	"github.com/google/uuid"
)

func loginFixture() login.Login {
	l, err := login.New("stud26")
	if err != nil {
		panic(err)
	}

	return l
}

func hashFixture() passhash.Hash {
	h, err := passhash.New("$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy")
	if err != nil {
		panic(err)
	}

	return h
}

func accountFixture() *account.Account {
	personID := uuid.New()
	acc, err := account.New(loginFixture(), hashFixture(), role.NewStudent(), personID)
	if err != nil {
		panic(err)
	}

	return acc
}
