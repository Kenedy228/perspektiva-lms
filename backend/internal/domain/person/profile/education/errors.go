package education

import "errors"

// ErrInvalid базовый sentinel-маркер для всех ошибок валидации сведений об образовании.
// Используется как корневая ошибка, оборачиваемая через fmt.Errorf("%w: ...").
var ErrInvalid = errors.New("некорректные сведения об образовании")
