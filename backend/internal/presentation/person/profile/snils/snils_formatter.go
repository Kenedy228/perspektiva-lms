package snils

import (
	"fmt"

	"gitflic.ru/lms/backend/internal/domain/person/profile/snils"
)

type PlainFormatter struct{}

func (f PlainFormatter) Format(s snils.SNILS) string {
	runes := []rune(s.Value())

	return fmt.Sprintf(
		"%c%c%c%c%c%c%c%c%c%c%c",
		runes[0],
		runes[1],
		runes[2],
		runes[3],
		runes[4],
		runes[5],
		runes[6],
		runes[7],
		runes[8],
		runes[9],
		runes[10],
	)
}

type WithDashesFormatter struct{}

func (f WithDashesFormatter) Format(s snils.SNILS) string {
	runes := []rune(s.Value())

	return fmt.Sprintf(
		"%c%c%c-%c%c%c-%c%c%c %c%c",
		runes[0],
		runes[1],
		runes[2],
		runes[3],
		runes[4],
		runes[5],
		runes[6],
		runes[7],
		runes[8],
		runes[9],
		runes[10],
	)
}

type WithSpacesFormatter struct{}

func (f WithSpacesFormatter) Format(s snils.SNILS) string {
	runes := []rune(s.Value())

	return fmt.Sprintf(
		"%c%c%c %c%c%c %c%c%c %c%c",
		runes[0],
		runes[1],
		runes[2],
		runes[3],
		runes[4],
		runes[5],
		runes[6],
		runes[7],
		runes[8],
		runes[9],
		runes[10],
	)
}
