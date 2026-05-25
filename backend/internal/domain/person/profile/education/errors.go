package education

import "fmt"

// ErrInvalid базовый sentinel-маркер для всех ошибок валидации сведений об образовании.
// Используется как корневая ошибка, оборачиваемая через fmt.Errorf("%w: ...").
var ErrInvalid = fmt.Errorf("некорректные сведения об образовании")
