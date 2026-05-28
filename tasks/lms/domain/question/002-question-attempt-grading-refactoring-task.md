# Задача для opencode: финальная доработка question/grading/attempt архитектуры

## Цель

Завершить рефакторинг пакетов вопросов, ответов, проверки ответов и прохождения тестирования в LMS backend.

Текущая кодовая база уже стала значительно лучше:
- улучшен application layer;
- улучшены ошибки;
- улучшены тесты;
- use case'ы стали чище;
- структура домена стала понятнее.

Но архитектурный рефакторинг ещё не завершён.

Основная цель этой задачи:
- убрать оставшиеся хрупкие места;
- завершить разделение validation и grading;
- убрать ручной runtime dispatch checker'ов;
- сохранить простоту Go-кода;
- не уйти в overengineering.

## Контекст проекта

Проект:
```txt
perspektiva-lms
```

Основные пакеты:

```txt
backend/internal/domain/question
backend/internal/domain/attempt
backend/internal/domain/grading

backend/internal/application/usecases/question/grading
backend/internal/application/usecases/attempt
```

## Текущее состояние

После предыдущего рефакторинга:

### Уже хорошо

- `GradeUseCase` стал чище;
- validation input улучшен;
- ошибки стали понятнее;
- домен `Attempt` выглядит сильным;
- snapshot strategy работает хорошо;
- тестов стало больше;
- package boundaries в целом хорошие.

### Оставшиеся проблемы

#### 1. GradeUseCase всё ещё вручную перебирает checker'ы

Сейчас используется:

```go
for i := range uc.checkers {
    if uc.checkers[i].Supports(q.Type()) {
        return q, uc.checkers[i], nil
    }
}
```

Проблемы:
- runtime dispatch;
- скрытая конфигурация;
- нет защиты от duplicate checker;
- use case знает слишком много о механике выбора checker;
- lookup O(n).

#### 2. Validation и grading всё ещё связаны

Сейчас:

```go
ValidateAnswerUseCase
    -> GradeUseCase.Execute()
```

Это означает:
- validation фактически использует grading;
- checker участвует в validation;
- score computation и answer compatibility смешаны.

#### 3. Checker abstraction всё ещё слишком перегружен

Checker сейчас отвечает сразу за:
- type support;
- type assertions;
- validation;
- grading.

Это создаёт лишнюю сложность.

## Что нужно сделать

==================================================
ЧАСТЬ 1 — Ввести CheckerRegistry
==================================================

Нужно убрать ручной перебор checker'ов.

Создать registry:

```go
type CheckerRegistry struct {
    checkers map[question.Type]grading.Checker
}
```

Пример API:

```go
func NewCheckerRegistry(
    checkers ...grading.Checker,
) (*CheckerRegistry, error)

func (r *CheckerRegistry) Get(
    t question.Type,
) (grading.Checker, error)
```

Требования:

### Registry должен:

- валидировать duplicate checker для одного question.Type;
- возвращать понятную ошибку;
- иметь deterministic behavior;
- использовать map lookup вместо перебора;
- регистрироваться в composition root.

### GradeUseCase НЕ должен:

- перебирать checker'ы;
- вызывать Supports();
- знать механику выбора checker.

После рефакторинга:

```go
checker, err := registry.Get(q.Type())
```

и всё.

==================================================
ЧАСТЬ 2 — Убрать Supports()
==================================================

После внедрения registry:
- удалить `Supports()` из checker interface;
- убрать runtime-dispatch по Supports();
- checker должен только выполнять grading.

Было:

```go
type Checker interface {
    Check(question.Question, question.Answer) score.Score
    Supports(question.Type) bool
}
```

Должно стать проще.

==================================================
ЧАСТЬ 3 — Разделить Validation и Grading
==================================================

Нужно ввести отдельный validator.

Пример:

```go
type AnswerValidator interface {
    Validate(
        q question.Question,
        ans question.Answer,
    ) error
}
```

Validator должен отвечать только за:

- nil checks;
- совместимость типа вопроса и ответа;
- корректность структуры ответа;
- существование referenced IDs;
- ограничения бизнес-правил ответа.

Validator НЕ должен:
- считать score;
- заниматься grading;
- выбирать checker.

==================================================
ЧАСТЬ 4 — Обновить ValidateAnswerUseCase
==================================================

Сейчас:

```go
ValidateAnswerUseCase
    -> GradeUseCase
```

Нужно:
- разорвать эту связь;
- сделать отдельный validation flow.

Должно быть примерно так:

```txt
ValidateAnswerUseCase
    -> QuestionRepository
    -> AnswerValidator

GradeUseCase
    -> QuestionRepository
    -> AnswerValidator
    -> CheckerRegistry
    -> Checker
```

==================================================
ЧАСТЬ 5 — Локализовать runtime type assertions
==================================================

Type assertions допустимы:
- внутри checker;
- внутри validator.

Они НЕ должны:
- появляться в use case;
- появляться в HTTP handler;
- дублироваться по проекту.

Нужно проверить:
- не размазались ли assertions;
- не дублируется ли логика compatibility.

==================================================
ЧАСТЬ 6 — НЕ УСЛОЖНИТЬ КОД
==================================================

Критически важно:

НЕ делать:
- Java-style architecture;
- interface ради interface;
- generic abstraction hell;
- сложные factories;
- service locator;
- event bus;
- visitor pattern;
- reflection.

Нужно:
- explicit Go code;
- простые зависимости;
- минимальные abstraction;
- читаемый orchestration layer.

==================================================
ЧАСТЬ 7 — Проверить Attempt Aggregate
==================================================

Attempt уже выглядит хорошо.

Нужно:
- не сломать snapshot strategy;
- не убрать Clone();
- не ослабить инварианты;
- не дать утечку mutable state;
- проверить restore methods;
- проверить lifecycle transitions.

Можно улучшить:
- читаемость;
- ошибки;
- edge-case тесты.

Но НЕ переписывать aggregate заново.

==================================================
ЧАСТЬ 8 — Добавить тесты
==================================================

Нужно покрыть:

### CheckerRegistry

- duplicate checker;
- missing checker;
- valid registry;
- deterministic lookup.

### AnswerValidator

- nil question;
- nil answer;
- wrong answer type;
- wrong question type;
- invalid IDs;
- invalid structure;
- empty answers;
- valid answers.

### GradeUseCase

- checker найден;
- checker не найден;
- validator error;
- grading success;
- repository error.

### ValidateAnswerUseCase

- validation success;
- validation failure;
- repository error.

### Attempt

- immutable snapshots;
- restore invalid state;
- deadline edge cases;
- state transitions.

==================================================
ЧАСТЬ 9 — Документация
==================================================

Обновить:

```txt
docs/diagrams/domain-question-class-diagram.md
docs/diagrams/domain-question-sequence-diagram.md
```

Если архитектура изменилась — диаграммы должны отражать:
- CheckerRegistry;
- AnswerValidator;
- новый validation flow.

Также создать:

```txt
docs/refactoring/final-question-grading-refactoring.md
```

Внутри:
- какие проблемы были;
- что было исправлено;
- какие компромиссы остались;
- почему были приняты такие решения;
- почему НЕ был сделан overengineering.

==================================================
ОГРАНИЧЕНИЯ
==================================================

- Не ломать публичное API без необходимости.
- Не делать rewrite.
- Не менять бизнес-логику.
- Не усложнять package structure.
- Не плодить интерфейсы.
- Не тащить лишние зависимости.
- Не превращать Go-код в enterprise Java.

==================================================
КРИТЕРИИ ГОТОВНОСТИ
==================================================

Задача считается выполненной, если:

1. `GradeUseCase` больше не перебирает checker'ы.
2. `Supports()` удалён.
3. Появился `CheckerRegistry`.
4. Validation отделён от grading.
5. `ValidateAnswerUseCase` больше не использует `GradeUseCase`.
6. Runtime type assertions локализованы.
7. Все тесты проходят.
8. Добавлены новые тесты.
9. Код стал проще, а не сложнее.
10. Attempt aggregate сохранил инварианты.
11. Mermaid diagrams обновлены.
12. Документация обновлена.

==================================================
ФОРМАТ РЕЗУЛЬТАТА
==================================================

После выполнения показать:

1. список изменённых файлов;
2. архитектурный overview;
3. какие проблемы были исправлены;
4. какие компромиссы остались;
5. результат test suite;
6. обновлённые Mermaid diagrams;
7. краткое объяснение новой architecture flow.
