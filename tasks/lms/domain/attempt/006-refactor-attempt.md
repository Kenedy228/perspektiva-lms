# Рефакторинг доменного модуля attempt под новые пакеты question и grading

## Контекст

Необходимо отрефакторить доменный модуль прохождения теста:

```text
internal/domain/attempt
```

Модуль содержит aggregate `Attempt`, подпакеты `attempt/item` и `attempt/answer`, параметры создания/восстановления попытки, проверки ответов и управление жизненным циклом прохождения теста.

Пакеты `question` и `grading` были изменены, поэтому текущая реализация `attempt` может использовать устаревшие типы, несуществующие импорты, устаревшие методы и лишнюю доменную логику.

---

# Цель 1. Убрать несуществующие и устаревшие импорты answer-пакетов

В текущем коде есть функция:

```go
func validateAnswerForQuestion(q question.Question, ans question.Answer) error
```

Она может использовать устаревшие или уже несуществующие импорты:

```go
matchinganswer "gitflic.ru/lms/backend/internal/domain/question/matching/answer"
selectableanswer "gitflic.ru/lms/backend/internal/domain/question/selectable/answer"
sequenceanswer "gitflic.ru/lms/backend/internal/domain/question/sequence/answer"
shortanswer "gitflic.ru/lms/backend/internal/domain/question/short/answer"
typedanswer "gitflic.ru/lms/backend/internal/domain/question/typed/answer"
```

Необходимо:

* проверить актуальный API пакета `question`;
* проверить актуальный API доменных ответов;
* удалить несуществующие импорты;
* адаптировать проверку соответствия ответа типу вопроса под новый домен;
* если проверка типа ответа больше относится к `grading`, перенести логику в сторону использования доменного `grading.Checker`, либо удалить дублирование из `attempt`;
* не оставлять compile-time dead imports.

---

# Цель 2. Проверить соответствие модуля attempt новому question domain

Необходимо проверить все места, где `attempt` зависит от `question`.

Особое внимание уделить:

```go
question.Question
question.Answer
question.Type
Question.ID()
Question.Type()
Question.Clone()
Answer.Clone()
```

Проверить:

* актуальны ли эти методы;
* существует ли ещё `question.Answer`;
* существует ли ещё `Clone` у вопроса и ответа;
* правильно ли `item.Item` хранит snapshot вопроса;
* правильно ли `answer.Entry` хранит snapshot ответа;
* не используются ли устаревшие value objects;
* не дублируются ли доменные инварианты question-пакета.

Если API изменился — адаптировать `attempt`, `attempt/item` и `attempt/answer`.

---

# Цель 3. Проверить соответствие модуля attempt новому grading domain

Необходимо проанализировать, должен ли `attempt` сам проверять соответствие типа ответа типу вопроса.

Проверить архитектурно:

* должна ли попытка знать конкретные типы ответов selectable/matching/sequence/typed/short;
* не является ли это ответственностью `grading`;
* можно ли оставить в `attempt` только проверку наличия вопроса, статуса и дедлайна;
* не нарушает ли текущая реализация границы домена из-за прямых импортов answer-подпакетов question.

Если новый `grading` домен отвечает за проверку типа ответа, убрать дублирующую type-switch логику из `attempt`.

Если attempt всё ещё должен проверять соответствие — реализовать это через актуальный API, без несуществующих импортов.

---

# Цель 4. Проверить и удалить неиспользуемые методы, функции и типы

Необходимо проверить пакет на наличие мёртвого кода.

Проверить:

* неиспользуемые функции;
* неиспользуемые методы;
* неиспользуемые типы;
* неиспользуемые параметры;
* устаревшие helper-функции;
* устаревшие validator-функции;
* устаревшие поля в `Params` и `RestoreParams`;
* устаревшие методы aggregate `Attempt`.

Особое внимание уделить:

```go
hasQuestion
findItem
validateAnswerForQuestion
Params
RestoreParams
Status
Status.Title
CountItems
CountAnswers
CanModify
```

Требования:

* удалить только реально неиспользуемый или устаревший код;
* не удалять публичное API без анализа использования;
* если публичный метод всё ещё нужен application/infrastructure слою — оставить;
* не выполнять unrelated-рефакторинг.

---

# Цель 5. Русифицировать ошибки и улучшить контекст

Необходимо перевести ошибки пакета `attempt`, `attempt/item`, `attempt/answer` на русский язык.

Текущие ошибки вида:

```go
fmt.Errorf("%w: invalid value", ErrInvalid)
fmt.Errorf("%w: invalid value (%d)", ErrInvalid, i)
fmt.Errorf("%w: invalid value (%s)", ErrInvalid, q.ID())
fmt.Errorf("%w: ID %s", ErrNotFound, questionID)
```

нужно заменить на более понятные.

Примеры ожидаемого стиля:

```go
fmt.Errorf("%w: идентификатор записи на курс обязателен", ErrInvalid)
fmt.Errorf("%w: идентификатор теста обязателен", ErrInvalid)
fmt.Errorf("%w: список вопросов не должен быть пустым", ErrInvalid)
fmt.Errorf("%w: вопрос с индексом %d не должен быть nil", ErrInvalid, i)
fmt.Errorf("%w: идентификатор вопроса с индексом %d обязателен", ErrInvalid, i)
fmt.Errorf("%w: вопрос с идентификатором %s уже добавлен в попытку", ErrInvalid, q.ID())
fmt.Errorf("%w: вопрос с идентификатором %s не найден в попытке", ErrNotFound, questionID)
fmt.Errorf("%w: попытку нельзя изменить в статусе %s", ErrStateConflict, a.status)
fmt.Errorf("%w: время завершения обязательно для статуса %s", ErrInvalid, status)
```

Требования:

* использовать wrapping через `%w`;
* не терять sentinel errors;
* не оставлять английские сообщения;
* добавлять контекст: id, index, status, current/expected value;
* не делать сообщения чрезмерно длинными.

---

# Цель 6. Проверить Status и жизненный цикл Attempt

Необходимо проверить модель статусов:

```go
StatusFinished
StatusExpired
StatusInProgress
StatusInterrupted
StatusCancelled
```

Проверить методы:

```go
Finish
SetExpired
Interrupt
Cancel
CanModify
validateFinishedAt
validateTimeline
ensureBeforeDeadline
```

Требования:

* проверить корректность переходов статусов;
* проверить работу дедлайна;
* проверить восстановление attempt в разных статусах;
* проверить finishedAt для финальных и нефинальных статусов;
* проверить deadlineAt относительно startedAt;
* русифицировать ошибки;
* не менять бизнес-логику без необходимости.

---

# Цель 7. Проверить item.Item

Необходимо проверить подпакет:

```text
internal/domain/attempt/item
```

Проверить:

```go
type Item struct {
	snapshot question.Question
}
```

Требования:

* проверить актуальность `question.Question.Clone`;
* убедиться, что snapshot не отдаётся наружу мутируемым;
* если `Clone` больше не существует — адаптировать snapshot-механику под новый question domain;
* добавить/обновить тесты;
* добавить godoc на русском языке.

---

# Цель 8. Проверить answer.Entry

Необходимо проверить подпакет:

```text
internal/domain/attempt/answer
```

Проверить:

```go
type Entry struct {
	questionID uuid.UUID
	answer     question.Answer
	answeredAt time.Time
}
```

Требования:

* проверить актуальность `question.Answer.Clone`;
* убедиться, что ответ не отдаётся наружу мутируемым;
* если `Clone` больше не существует — адаптировать под новый question domain;
* проверить валидаторы `questionID`, `answer`, `answeredAt`;
* русифицировать ошибки;
* добавить/обновить тесты;
* добавить godoc на русском языке.

---

# Цель 9. Обновить и расширить тесты

Необходимо проверить и обновить тесты пакетов:

```text
internal/domain/attempt
internal/domain/attempt/item
internal/domain/attempt/answer
```

Тесты должны покрывать:

## Attempt.New

* успешное создание попытки;
* пустой enrollmentID;
* пустой quizID;
* пустой список вопросов;
* nil-вопрос в списке;
* вопрос с пустым ID;
* дубликаты вопросов;
* корректное выставление deadline при time limit;
* infinite time limit без deadline;
* ошибки создания item snapshot.

## Attempt.Restore

* успешное восстановление;
* пустой attempt ID;
* некорректный enrollmentID;
* некорректный quizID;
* некорректный status;
* пустой startedAt;
* некорректный finishedAt для статуса;
* некорректный timeline;
* ответы на отсутствующие вопросы;
* несовпадение key questionID и entry.QuestionID;
* ответ неподходящего типа, если такая проверка остаётся в attempt.

## Attempt.AddAnswer

* успешное добавление ответа;
* попытка не в статусе `in_progress`;
* ответ после deadline;
* вопрос не найден;
* некорректный answer;
* замена существующего ответа на тот же вопрос.

## Lifecycle

* `Finish`;
* `SetExpired`;
* `Interrupt`;
* `Cancel`;
* ошибки при некорректных переходах;
* ошибки при пустом времени;
* ошибки при времени раньше `startedAt`;
* ошибки дедлайна.

## item.Item

* успешное создание snapshot;
* nil question;
* snapshot возвращается безопасно;
* ID возвращается корректно.

## answer.Entry

* успешное создание entry;
* пустой questionID;
* nil answer;
* пустое answeredAt;
* answer возвращается безопасно;
* answeredAt возвращается корректно.

Требования:

* использовать table-driven tests;
* проверять errors.Is для sentinel errors;
* не тестировать повторно внутреннюю логику question/grading;
* тестировать именно инварианты attempt;
* избегать хрупких тестов на полный текст ошибки.

---

# Цель 10. Добавить godoc

Необходимо добавить или обновить godoc-комментарии на русском языке для публичного API:

* `Attempt`;
* `New`;
* `Restore`;
* `Params`;
* `RestoreParams`;
* `Status`;
* публичные методы `Attempt`;
* публичные методы `Status`;
* `item.Item`;
* `item.New`;
* публичные методы `item.Item`;
* `answer.Entry`;
* `answer.New`;
* публичные методы `answer.Entry`;
* экспортируемые ошибки.

Требования:

* комментарий должен начинаться с имени документируемой сущности;
* комментарий должен объяснять назначение и поведение;
* избегать формальных комментариев без смысла.

---

# Цель 11. Проверить архитектурные нюансы

После рефакторинга необходимо проверить архитектуру модуля.

Проверить:

* не зависит ли `attempt` от конкретных подпакетов answer вопроса без необходимости;
* не дублирует ли `attempt` responsibility `grading`;
* не дублирует ли `attempt` валидацию question-пакета;
* не протекает ли application/infrastructure в domain;
* не создаёт ли snapshot question/answer проблем с мутабельностью;
* не является ли `validateAnswerForQuestion` нарушением границ;
* не являются ли `Params` и `RestoreParams` слишком широкими;
* нет ли неиспользуемого кода;
* нет ли лишних алиасов импортов.

Если проблема относится к текущей задаче — исправить.
Если проблема крупнее текущего рефакторинга — описать её в финальном отчёте без unrelated-изменений.

---

# Ограничения

* Не менять application-слой без необходимости.
* Не менять infrastructure-слой без необходимости.
* Не менять transport-слой.
* Не менять question domain без необходимости.
* Не менять grading domain без необходимости.
* Не добавлять лишние интерфейсы.
* Не добавлять новые абстракции без явной необходимости.
* Не дублировать question/grading validation.
* Не выполнять unrelated-рефакторинг.
* Сохранять существующий стиль проекта.

---

# Проверка результата

После выполнения задачи необходимо выполнить:

```bash
go test ./internal/domain/attempt/...
go test ./...
go vet ./...
```

---

# Ожидаемый результат

В результате выполнения задачи:

* пакет `attempt` совместим с актуальными `question` и `grading`;
* несуществующие импорты удалены;
* устаревшие методы, функции и типы удалены или адаптированы;
* ошибки переведены на русский язык;
* тесты обновлены и расширены;
* godoc добавлен;
* архитектурные проблемы проверены;
* проект успешно проходит `go test ./...` и `go vet ./...`.
