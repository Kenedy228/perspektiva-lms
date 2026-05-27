package dob

import "errors"

// ErrInvalid базовый sentinel-маркер для всех ошибок валидации даты рождения.
// Используется как корневая ошибка, оборачиваемая через fmt.Errorf("%w: ...").
var ErrInvalid = errors.New("некорректная дата рождения")
