package dob

import "time"

type DateOfBirth struct {
	date time.Time
}

func New(date, asOf time.Time) (DateOfBirth, error) {
	nDate := normalize(date)
	nAsOf := normalize(asOf)

	if err := validateAdultDateOfBirth(nDate, nAsOf); err != nil {
		return DateOfBirth{}, err
	}

	return DateOfBirth{
		date: nDate,
	}, nil
}

func (db DateOfBirth) Date() time.Time {
	return db.date
}
