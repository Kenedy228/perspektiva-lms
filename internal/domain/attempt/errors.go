package attempt

import "errors"

var (
	ErrInvalidItem      = errors.New("ошибка элемента попытки")
	ErrInvalidAttempt   = errors.New("ошибка теста курса")
	ErrInfiniteDeadline = errors.New("попытка пройти тест не может быть просречена, так как не имеет дедлайна")
	ErrNotExpiredYet    = errors.New("попытка пройти тест еще не просрочена")
	ErrNotModifiable    = errors.New("попытка пройти тест неизменяема (завершена или истекла)")
	ErrUnexistingItem   = errors.New("вопроса с таким идентификатором в попытке нет")
	ErrNotFinishedYet   = errors.New("попытка пройти тест еще не завершена")
	ErrInactive         = errors.New("попытка пройти тест не активна")
)
