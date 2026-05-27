# Анализ: question / attempt / grading — рефакторинг

Дата: 2026-05-28

## 1. Текущее состояние (baseline)

Все существующие тесты проходят: `go test ./internal/domain/question/... ./internal/domain/attempt/... ./internal/domain/grading/... ./internal/application/usecases/question/grading/...` — PASS.

### 1.1. Структура пакетов

```
domain/question/              # Question + Answer интерфейсы, Type enum
  base/                       # Base struct (id, title) — embedding target
    title/                    # Title value object
  matching/                   # MatchingQuestion + Pair/Prompt/Match VO
    pair/                     
    answer/                   # MatchingAnswer
  selectable/                 # SelectableQuestion + Option VO
    option/
    answer/                   # SelectableAnswer + OptionID VO
  sequence/                   # SequenceQuestion + Option VO (без ID)
    option/
    answer/                   # SequenceAnswer + OptionID VO
  short/                      # ShortQuestion + Variant VO
    variant/
    answer/                   # ShortAnswer
domain/attempt/               # Attempt aggregate root
  item/                       # Item (question snapshot)
  answer/                     # Entry (answer snapshot)
domain/grading/               # Checker interface + Score VO
  matching/                   # MatchingChecker
  selectable/                 # SelectableChecker
  sequence/                   # SequenceChecker
  short/                      # ShortChecker
  score/                      # Score value object [0..1]
application/usecases/question/grading/
                              # GradeUseCase, ValidateAnswerUseCase
```

### 1.2. Выявленные проблемы

| # | Проблема | Локация | Серьёзность |
|---|----------|---------|-------------|
| P1 | **Валидация смешана с оцениванием.** `Checker.Check()` одновременно проверяет совместимость типов (type assertions) и вычисляет score. `ValidateAnswerUseCase` вызывает `GradeUseCase` и отбрасывает score — wasteful. | `grading/selectable/checker.go:19-28`, `grading.go:131-133` | Средняя |
| P2 | **Перебор checker'ов через Supports().** `GradeUseCase.loadQuestionAndChecker()` итерирует слайс O(n). Нет защиты от дубликатов checker'ов для одного типа. Ошибка «checker не найден» генерируется в application-слое. | `grading.go:84-88` | Средняя |
| P3 | **Дублирование sentinel errors.** `ErrInvalidQuestionType` и `ErrInvalidAnswerType` определены в 4 пакетах идентично: `grading/{matching,selectable,sequence,short}/errors.go`. | 4 файла | Низкая |
| P4 | **ValidateAnswerUseCase завязан на GradeUseCase.** `ValidateAnswerUseCase` — это тонкая обёртка над `GradeUseCase`, которая вычисляет score и отбрасывает его. | `grading.go:114-133` | Средняя |
| P5 | **Отсутствует структурная валидация ответа.** Checker'ы не проверяют, что ID опций/пар из ответа существуют в вопросе. Несуществующие ID молча трактуются как «неправильно». | Все checker'ы | Низкая (feature gap) |
| P6 | **Runtime type assertions размазаны.** Присутствуют в: а) Checker.Check() (domain), б) GradeUseCase.loadQuestionAndChecker() через Supports() (application), в) buildAnswer() в handler (transport). В (б) они опосредованные через Supports(). | `checker.go`, `grading.go`, `question_api.go` | Низкая |

### 1.3. Что уже сделано хорошо (сохранить)

- **Immutable snapshot в Attempt.** `item.New()` клонирует Question через `Clone()`. `answer.New()` клонирует Answer через `Clone()`. Инвариант: изменение вопроса преподавателем не влияет на активные попытки.
- **Машина состояний Attempt.** 5 статусов, `CanModify()` защищает от изменений после завершения.
- **Защита коллекций.** `Items()` и `Answers()` возвращают копии через `slices.Clone`/`maps.Clone`.
- **Value Objects с инвариантами.** Title, Option, Pair, Variant, Score — все с фабриками и валидацией.
- **Embedding Base.** Четыре типа вопросов разделяют общее состояние через `*base.Base`.
- **Тесты.** Покрыты все ключевые сценарии: создание, restore, lifecycle, проверка ответов.

## 2. План рефакторинга

### 2.1. Ввести AnswerValidator (разделение валидации и оценивания)

```go
// domain/grading/validator.go
type AnswerValidator interface {
    Validate(q question.Question, a question.Answer) error
}
```

Валидатор проверяет:
- `q != nil`
- `a != nil`
- тип ответа соответствует типу вопроса (через type assertion)

Каждый тип вопроса получает свой validator, аналогично checker'ам. Type assertions остаются внутри validator'а и checker'а (локализованы в domain-слое).

**ValidateAnswerUseCase** будет использовать `AnswerValidator` напрямую (не через `GradeUseCase`).
**GradeUseCase** будет использовать и `AnswerValidator`, и `Checker` (через registry).

### 2.2. Ввести CheckerRegistry

```go
// domain/grading/registry/registry.go
type Registry struct {
    checkers map[question.Type]Checker
}
func New(checkers ...Checker) (*Registry, error)
func (r *Registry) Get(t question.Type) (Checker, error)
```

- Замена O(n) перебора `Supports()` на O(1) map lookup.
- Валидация дубликатов checker'ов при создании.
- Понятная ошибка при отсутствии checker'а для типа.
- Инстанцируется в composition root (`main.go`).

### 2.3. Локализовать type assertions

После изменений:
- `Checker.Check()` — type assertions остаются как defensive мера (domain).
- `AnswerValidator.Validate()` — type assertions для проверки совместимости (domain).
- `Registry.Get()` — map lookup, без assertions (domain).
- `GradeUseCase` — вызывает `Validator.Validate()`, затем `Registry.Get()`, затем `Checker.Check()`. Без собственных type assertions (application).
- `buildAnswer()` в handler — ok, это DTO → domain mapping (transport).

### 2.4. Унифицировать sentinel errors

Вынести `ErrInvalidQuestionType` и `ErrInvalidAnswerType` в пакет `domain/grading` (вместо дублирования в 4 пакетах).

### 2.5. Пакетная структура

Оставить без изменений. Текущая структура `question/{type}/answer/` оправдана:
- Каждый answer-пакет содержит типы, специфичные для этого типа вопроса (например, `matching/answer/Pair`).
- Размер пакетов мал (2-5 файлов), cohesion высокий.
- Перенос answer-файлов в родительский пакет создаст циклические импорты (answer импортирует question-тип, question-тип может потенциально ссылаться на answer).

### 2.6. Структурная валидация ответа (P5)

Оставить на будущее. Требует значительного расширения AnswerValidator (проверка существования ID опций/пар). Текущее поведение (silent incorrect) не является багом — это conscious design choice: невалидные ID просто не матчатся с правильными.

## 3. Ожидаемые изменения файлов

| Файл | Действие |
|------|----------|
| `domain/grading/validator.go` | **Новый**: AnswerValidator interface |
| `domain/grading/errors.go` | **Новый**: общие sentinel errors |
| `domain/grading/registry/registry.go` | **Новый**: CheckerRegistry |
| `domain/grading/registry/registry_test.go` | **Новый**: тесты registry |
| `domain/grading/selectable/checker.go` | Изменить: переиспользовать общие errors |
| `domain/grading/selectable/errors.go` | Удалить (перенесено в grading/errors.go) |
| `domain/grading/selectable/validator.go` | **Новый**: SelectableAnswerValidator |
| `domain/grading/selectable/validator_test.go` | **Новый**: тесты validator |
| `domain/grading/matching/checker.go` | Изменить: переиспользовать общие errors |
| `domain/grading/matching/errors.go` | Удалить |
| `domain/grading/matching/validator.go` | **Новый** |
| `domain/grading/matching/validator_test.go` | **Новый** |
| `domain/grading/sequence/checker.go` | Изменить |
| `domain/grading/sequence/errors.go` | Удалить |
| `domain/grading/sequence/validator.go` | **Новый** |
| `domain/grading/sequence/validator_test.go` | **Новый** |
| `domain/grading/short/checker.go` | Изменить |
| `domain/grading/short/errors.go` | Удалить |
| `domain/grading/short/validator.go` | **Новый** |
| `domain/grading/short/validator_test.go` | **Новый** |
| `application/usecases/question/grading/grading.go` | Изменить: использовать Registry + Validator |
| `application/usecases/question/grading/grading_test.go` | Изменить: обновить тесты |
| `cmd/lms-api/main.go` | Изменить: wire Registry + Validators |
| `docs/diagrams/domain-question-class-diagram.md` | Обновить |
| `docs/diagrams/domain-question-sequence-diagram.md` | Обновить |

## 4. Результаты (критерии приёмки)

- [x] Все существующие тесты проходят (25 пакетов, все PASS)
- [x] `ValidateAnswerUseCase` использует `AnswerValidator`, не `GradeUseCase`
- [x] `GradeUseCase` использует `CheckerRegistry.Get()`, не перебор `Supports()`
- [x] Type assertions не выходят за пределы domain-слоя
- [x] Нет дублирования sentinel errors (унифицированы в `grading/errors.go`)
- [x] `Attempt` инварианты сохранены (тесты pass, код не менялся)
- [x] Mermaid-диаграммы обновлены
- [x] Добавлены тесты: 4 validator_test.go + registry_test.go + обновлён grading_test.go

### Добавленные файлы

| Файл | Назначение |
|------|-----------|
| `domain/grading/errors.go` | Унифицированные sentinel errors |
| `domain/grading/validator.go` | `AnswerValidator` interface |
| `domain/grading/registry/registry.go` | `CheckerRegistry` (O(1) map lookup) |
| `domain/grading/registry/registry_test.go` | Тесты registry |
| `domain/grading/registry/doc.go` | Package doc |
| `domain/grading/{selectable,matching,sequence,short}/validator.go` | Конкретные валидаторы |
| `domain/grading/{selectable,matching,sequence,short}/validator_test.go` | Тесты валидаторов |

### Удалённые файлы

| Файл | Причина |
|------|---------|
| `domain/grading/{selectable,matching,sequence,short}/errors.go` | Перенесено в `grading/errors.go` |

### Изменённые файлы

| Файл | Изменение |
|------|-----------|
| `domain/grading/{selectable,matching,sequence,short}/checker.go` | Импорт `grading` + использование общих errors |
| `domain/grading/{selectable,matching,sequence,short}/checker_test.go` | Импорт `grading` + ссылки на общие errors |
| `application/usecases/question/grading/grading.go` | Registry + Validator вместо `[]Checker` + `Supports()` |
| `application/usecases/question/grading/grading_test.go` | Обновлены конструкторы и тесты |
| `cmd/lms-api/main.go` | Wire `CheckerRegistry` + `answerValidators` map |
| `docs/diagrams/domain-question-class-diagram.md` | Добавлены AnswerValidator, Registry, обновлены связи |
| `docs/diagrams/domain-question-sequence-diagram.md` | Разделён validate/grading flow, добавлен Registry |
