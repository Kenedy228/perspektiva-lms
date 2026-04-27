package person

import (
	"gitflic.ru/lms/internal/domain/person/name"
	"gitflic.ru/lms/internal/domain/person/profile"
)

type Params struct {
	Name    name.Name
	Profile *profile.Profile
}
