package dob

import "time"

// DateOfBirth объект-значение нормализованной даты рождения взрослого человека.
type DateOfBirth struct {
	date time.Time
}

// New создает DateOfBirth из переданной даты рождения, нормализует её и проверяет,
// что на момент asOf человек достиг совершеннолетия.
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

// Date возвращает нормализованную (UTC, без времени) дату рождения.
func (db DateOfBirth) Date() time.Time {
	return db.date
}

// IsZero возвращает true, если DateOfBirth не был инициализирован через New.
func (db DateOfBirth) IsZero() bool {
	return db.date.IsZero()
}
