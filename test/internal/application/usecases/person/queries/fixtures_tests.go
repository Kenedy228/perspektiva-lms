package queries_test

import (
	personports "gitflic.ru/lms/internal/application/ports/person"
	"github.com/google/uuid"
)

func personShortViewSlicesFixture() []personports.PersonShortView {
	view := personports.PersonShortView{
		ID:               uuid.NewString(),
		FullName:         "Иванов Иван Иванович",
		OrganizationName: "ООО 'Ромашка'",
	}

	return []personports.PersonShortView{view, view, view}
}

func personDetailedFixture() personports.PersonDetailedView {
	view := personports.PersonDetailedView{
		ID:               uuid.NewString(),
		FirstName:        "Иван",
		LastName:         "Иванов",
		MiddleName:       "Иванович",
		Snils:            "",
		JobTitle:         "",
		Education:        "",
		DateOfBirth:      "2020-02-02",
		OrganizationName: "",
	}

	return view
}
