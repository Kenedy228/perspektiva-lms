// Пакет timelimit предоставляет функционал для создания лимитов времени в секундах.
package timelimit

import "time"

// Максимальный лимит (в секундах).
const maxSeconds = 24 * 60 * 60

// Структура Limit является контейнером секунд лимита.
type Limit struct {
	seconds int
}

// New возвращает новый лимит если значение принадлежит [0, maxSeconds].
func New(seconds int) (Limit, error) {
	if err := validateSeconds(seconds); err != nil {
		return Limit{}, err
	}

	return Limit{
		seconds: seconds,
	}, nil
}

// IsInfinite сигнализирует, что временной лимит отсутствует (бесконечность)
func (l Limit) IsInfinite() bool {
	return l.seconds == 0
}

// Seconds геттер для значения секунд лимита времени
func (l Limit) Seconds() int {
	return l.seconds
}

// TryDuration хелпер для конвертации секунд в time.Duration.
// Второе значение сигнализирует, что лимит бесконечный (значение false), или конечный (значение true).
func (l Limit) TryDuration() (time.Duration, bool) {
	if l.seconds == 0 {
		return 0, false
	}

	return time.Second * time.Duration(l.seconds), true
}
