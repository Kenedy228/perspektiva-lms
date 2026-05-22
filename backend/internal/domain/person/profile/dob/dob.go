package dob

import "time"

// DateOfBirth объект-значение нормализованная дата рождения человека.
type DateOfBirth struct {
	date time.Time
}

// New создает объект DateOfBirth.
func New(date, asOf time.Time) (DateOfBirth, error) {
	date = normalize(date)
	asOf = normalize(asOf)

	if err := validateAdultDateOfBirth(date, asOf); err != nil {
		return DateOfBirth{}, err
	}

	return DateOfBirth{
		date: date,
	}, nil
}

// Date возвращает нормализованную дату рождения.
func (db DateOfBirth) Date() time.Time {
	return db.date
}

// IsZero сигнализирует, был ли инициализирован объект.
func (db DateOfBirth) IsZero() bool {
	return db.date.IsZero()
}
