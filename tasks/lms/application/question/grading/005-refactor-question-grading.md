# Рефакторинг application-модуля grading под новый домен grading

## Контекст

Необходимо отрефакторить application usecase-модуль проверки и оценки ответов.

Пакет:

```text
internal/application/usecases/question/grading
```

Доменный пакет `internal/domain/grading` был изменён, поэтому application-код необходимо привести в соответствие с новым API домена.

Текущий модуль содержит:

* `GradeUseCase`;
* `ValidateAnswerUseCase`;
* `GradeInput`;
* `GradeOutput`;
* загрузку вопроса;
* выбор checker-а;
* вызов доменного grading checker-а.

---

# Цель 1. Адаптировать GradeUseCase под новый grading domain

Необходимо проанализировать актуальный API домена grading и обновить:

```go
GradeUseCase
NewGradeUseCase
GradeInput
GradeOutput
Execute
```

Требования:

* использовать актуальный интерфейс/тип checker-а из `internal/domain/grading`;
* использовать актуальный тип `score.Score`;
* сохранить сценарий: загрузить вопрос → выбрать checker → проверить ответ → вернуть score;
* не дублировать доменную логику проверки ответа в application-слое;
* application-слой должен только оркестрировать сценарий.

---

# Цель 2. Убрать двойную загрузку вопроса в ValidateAnswerUseCase

В текущей реализации `ValidateAnswerUseCase.Execute` сначала вызывает `uc.grade.load`, а затем вызывает `uc.grade.Execute`, из-за чего вопрос загружается дважды.

Необходимо исправить это поведение.

Требования:

* вопрос не должен загружаться два раза;
* если `ValidateAnswerUseCase` остаётся, он должен использовать `GradeUseCase.Execute` напрямую;
* если после анализа окажется, что `ValidateAnswerUseCase` больше не нужен из-за нового домена — удалить его вместе с input-типом и тестами;
* решение должно соответствовать архитектуре проекта.

Предпочтительный вариант, если usecase остаётся:

```go
func (uc *ValidateAnswerUseCase) Execute(ctx context.Context, in ValidateAnswerInput) error {
	_, err := uc.grade.Execute(ctx, GradeInput{
		QuestionID: in.QuestionID,
		Answer:     in.Answer,
	})
	return err
}
```

---

# Цель 3. Добавить валидацию входных данных application-слоя

Необходимо проверить входные данные до выполнения сценария.

Минимальные проверки:

* `QuestionID` не должен быть пустым;
* `QuestionID` должен быть корректным UUID;
* UUID не должен быть `uuid.Nil`;
* `Answer` не должен быть `nil`, если `question.Answer` является интерфейсом.

Требования:

* не дублировать доменные инварианты;
* проверять только application-level input;
* использовать sentinel error `ErrInvalidInput`;
* ошибки должны содержать понятный контекст.

---

# Цель 4. Улучшить выбор checker-а

Необходимо проверить текущую логику выбора checker-а:

```go
for i := range uc.checkers {
	if uc.checkers[i].Supports(q.Type()) {
		return q, uc.checkers[i], nil
	}
}
```

Требования:

* сохранить стратегию выбора checker-а по типу вопроса;
* адаптировать под новый API, если интерфейс checker-а изменился;
* не добавлять лишних абстракций без необходимости;
* при отсутствии checker-а возвращать `ErrUnsupportedChecker`;
* ошибка должна содержать тип вопроса.

---

# Цель 5. Улучшить именование внутренних методов

Текущий метод:

```go
load
```

имеет слишком общее название.

Необходимо переименовать его в более понятное, например:

```go
loadQuestionAndChecker
```

или другое имя, лучше отражающее смысл.

Требования:

* имя должно отражать, что метод загружает вопрос и подбирает checker;
* обновить все вызовы и тесты.

---

# Цель 6. Русифицировать ошибки

Необходимо перевести ошибки application-модуля grading на русский язык.

Текущие ошибки-кандидаты:

```go
ErrInvalidInput       = errors.New("question grading invalid input")
ErrUnsupportedChecker = errors.New("question grading checker not found")
```

Требования:

* ошибки должны быть на русском языке;
* сообщения должны быть понятными;
* использовать wrapping через `%w`;
* не терять sentinel errors;
* не оставлять английские сообщения вида:

    * `check answer`;
    * `parse question id`;
    * `question id is required`;
    * `find question`.

Пример ожидаемого стиля:

```go
var (
	ErrInvalidInput       = errors.New("некорректные данные для оценки ответа")
	ErrUnsupportedChecker = errors.New("проверяющий для типа вопроса не найден")
)
```

```go
fmt.Errorf("проверка ответа: %w", err)
fmt.Errorf("разбор идентификатора вопроса: %w", err)
fmt.Errorf("%w: идентификатор вопроса обязателен", ErrInvalidInput)
fmt.Errorf("поиск вопроса: %w", err)
fmt.Errorf("%w: тип вопроса %s", ErrUnsupportedChecker, q.Type())
```

---

# Цель 7. Обновить тесты

Необходимо добавить или обновить тесты для application-модуля grading.

Тесты должны покрывать:

* успешную оценку ответа;
* ошибку при пустом `QuestionID`;
* ошибку при некорректном UUID;
* ошибку при `uuid.Nil`;
* ошибку при `Answer == nil`, если применимо;
* ошибку repository при загрузке вопроса;
* ошибку checker-а при проверке ответа;
* отсутствие подходящего checker-а;
* корректный выбор checker-а по типу вопроса;
* корректный score в `GradeOutput`;
* отсутствие двойной загрузки вопроса в `ValidateAnswerUseCase`.

Требования:

* использовать table-driven tests;
* использовать fake/mock repository;
* использовать fake checker;
* проверять ошибки через `errors.Is`;
* не тестировать повторно доменную grading-логику;
* тестировать orchestration application-слоя.

---

# Цель 8. Добавить godoc

Необходимо добавить или обновить godoc-комментарии на русском языке для публичного API пакета.

Документация нужна для:

* `GradeUseCase`;
* `NewGradeUseCase`;
* `GradeInput`;
* `GradeOutput`;
* `GradeUseCase.Execute`;
* `ValidateAnswerUseCase`, если он остаётся;
* `NewValidateAnswerUseCase`, если он остаётся;
* `ValidateAnswerInput`, если он остаётся;
* `ValidateAnswerUseCase.Execute`, если он остаётся;
* экспортируемых ошибок.

Требования:

* комментарий должен начинаться с имени документируемой сущности;
* комментарий должен объяснять назначение usecase-а;
* избегать формальных комментариев без смысла.

---

# Цель 9. Проверить архитектурные проблемы

После рефакторинга необходимо проверить архитектуру модуля.

Проверить:

* не дублируется ли доменная grading-логика в application-слое;
* не нарушены ли границы application/domain;
* не протекает ли infrastructure;
* не используется ли доменный checker неправильно;
* нужен ли `ValidateAnswerUseCase`;
* не создаёт ли `GradeUseCase` слишком много ответственности;
* не нужно ли разделить input validation в отдельную приватную функцию;
* нет ли мёртвого кода;
* нет ли лишних алиасов импортов.

Если проблема относится к текущей задаче — исправить.
Если проблема крупнее текущего рефакторинга — описать её в финальном отчёте без unrelated-изменений.

---

# Ограничения

* Не менять domain/grading без необходимости.
* Не менять domain/question без необходимости.
* Не менять application ports без необходимости.
* Не менять infrastructure и transport без необходимости.
* Не дублировать доменную grading-логику.
* Не добавлять новые абстракции без явной необходимости.
* Не добавлять лишние интерфейсы.
* Не выполнять unrelated-рефакторинг.
* Сохранять существующий стиль проекта.
* Panic в конструкторах можно оставить, если это текущий стиль проекта.

---

# Проверка результата

После выполнения задачи необходимо выполнить:

```bash
go test ./internal/application/usecases/question/grading/...
go test ./...
go vet ./...
```

---

# Ожидаемый результат

В результате выполнения задачи:

* application-модуль grading соответствует новому domain/grading;
* `GradeUseCase` использует актуальный API домена;
* `ValidateAnswerUseCase` больше не загружает вопрос дважды или удалён как ненужный;
* входные данные валидируются на application-уровне;
* ошибки переведены на русский язык;
* тесты обновлены;
* godoc добавлен;
* архитектурные проблемы проверены;
* проект успешно проходит `go test ./...` и `go vet ./...`.
