package dob

import "time"

type DateOfBirth struct {
	dob time.Time
}

func NewDateOfBirth(dob time.Time) (DateOfBirth, error) {
	if err := validateDateOfBirth(dob); err != nil {
		return DateOfBirth{}, err
	}

	return DateOfBirth{
		dob: dob,
	}, nil
}

func (d DateOfBirth) Date() time.Time {
	return d.dob
}

func (d DateOfBirth) AgeAt(at time.Time) (int, error) {
	if err := validateForAgeAt(d.dob, at); err != nil {
		return -1, err
	}

	rawAge := at.Year() - d.dob.Year()

	if !d.hasHadBirthdayInYear(at) {
		rawAge--
	}

	return rawAge, nil
}

func (d DateOfBirth) hasHadBirthdayInYear(at time.Time) bool {
	leap := func(year int) bool {
		if year%400 == 0 {
			return true
		}

		if year%100 == 0 {
			return false
		}

		return year%4 == 0
	}

	if d.dob.Month() == 2 && d.dob.Day() == 29 && !leap(at.Year()) {
		pseudoDob := time.Date(at.Year(), 2, 28, 0, 0, 0, 0, time.UTC)
		return at.Compare(pseudoDob) != -1
	}

	if at.Month() > d.dob.Month() {
		return true
	}

	if at.Month() == d.dob.Month() && at.Day() >= d.dob.Day() {
		return true
	}

	return false
}
