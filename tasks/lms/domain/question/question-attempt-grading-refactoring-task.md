# Задача для opencode: улучшение пакетов тестирования и вопросов

## Цель

Провести точечное улучшение архитектуры пакетов, связанных с вопросами, ответами, проверкой ответов и прохождением тестирования в LMS backend.

Основная цель — снизить хрупкость текущей реализации, уменьшить зависимость от runtime type assertions, разделить ответственность между валидацией и оцениванием, а также сделать код проще для сопровождения и расширения новыми типами вопросов.

## Контекст проекта

Проект: `perspektiva-lms`

Backend написан на Go и использует DDD/чистую архитектуру.

Основные интересующие области:

```txt
backend/internal/domain/question
backend/internal/domain/question/selectable
backend/internal/domain/question/selectable/answer
backend/internal/domain/question/matching
backend/internal/domain/question/matching/answer
backend/internal/domain/question/sequence
backend/internal/domain/question/sequence/answer
backend/internal/domain/question/short
backend/internal/domain/question/short/answer

backend/internal/domain/attempt
backend/internal/domain/grading

backend/internal/application/usecases/question/grading
backend/internal/application/usecases/attempt
```

Если в кодовой базе есть дополнительные пакеты, связанные с тестами, quiz, attempts, grading, validation или answers, их тоже нужно учитывать.

## Текущие проблемы

Перед изменениями нужно подтвердить проблемы по коду.

Ожидаемые архитектурные слабые места:

1. Проверка ответов сильно завязана на runtime type assertions.
2. `Checker` принимает общий `question.Question` и общий `question.Answer`, после чего приводит их к конкретным типам.
3. Валидация совместимости ответа и вопроса смешана с оцениванием.
4. `ValidateAnswerUseCase` может быть слишком сильно завязан на `GradeUseCase`.
5. Выбор checker'а через перебор `Supports(q.Type())` может быть заменён на более явный registry.
6. Пакеты `question/*/answer` могут быть избыточно раздроблены.
7. Логика вопроса, ответа, проверки и попытки может быть распределена неравномерно.

Важно: не считать эти пункты автоматически истинными. Сначала проверить реальный код.

## Что нужно сделать

### 1. Провести анализ текущей реализации

Проанализировать:

- интерфейс `question.Question`;
- интерфейс `question.Answer`;
- конкретные типы вопросов;
- конкретные типы ответов;
- value objects вопросов;
- `attempt.Attempt`;
- `attempt.Item`;
- `attempt.Entry`;
- `grading.Checker`;
- конкретные checker'ы;
- use case'ы проверки и оценки;
- создание попытки;
- добавление ответа в попытку;
- завершение попытки;
- тесты доменного и application слоя.

После анализа кратко зафиксировать найденные проблемы в комментариях к PR или в отдельном файле:

```txt
docs/refactoring/question-attempt-grading-analysis.md
```

Если папки `docs/refactoring` нет — создать.

### 2. Разделить валидацию ответа и оценивание

Необходимо отделить проверку совместимости ответа с вопросом от подсчёта score.

Сейчас оценивание может выполнять сразу несколько задач:

- найти нужный checker;
- проверить тип вопроса;
- проверить тип ответа;
- посчитать score.

Нужно выделить отдельную ответственность:

```go
type AnswerValidator interface {
    Validate(q question.Question, ans question.Answer) error
}
```

Валидация должна проверять:

- вопрос не nil;
- ответ не nil;
- тип ответа соответствует типу вопроса;
- структура ответа допустима для данного вопроса;
- ответ не содержит ссылок на несуществующие option/prompt/match ID, если это применимо;
- пустой ответ обрабатывается согласно текущим бизнес-правилам проекта.

Оценивание должно заниматься только расчётом `Score`.

Пример желаемого разделения:

```txt
ValidateAnswerUseCase
    -> AnswerValidator.Validate(question, answer)

GradeUseCase
    -> AnswerValidator.Validate(question, answer)
    -> CheckerRegistry.Get(question.Type())
    -> Checker.Check(question, answer)
```

Важно:

- не ломать публичные контракты без необходимости;
- если текущие DTO/API завязаны на старое поведение, сохранить совместимость;
- если нужно изменить имена, делать это аккуратно.

### 3. Ввести registry для checker'ов

Заменить неявный перебор checker'ов через `Supports(q.Type())` на явный registry.

Пример идеи:

```go
type CheckerRegistry struct {
    checkers map[question.Type]grading.Checker
}

func NewCheckerRegistry(checkers ...grading.Checker) (*CheckerRegistry, error)

func (r *CheckerRegistry) Get(t question.Type) (grading.Checker, error)
```

Требования:

- registry должен валидировать дубли checker'ов для одного типа вопроса;
- registry должен возвращать понятную ошибку, если checker для типа вопроса не найден;
- `GradeUseCase` не должен сам перебирать checker'ы;
- регистрация checker'ов должна происходить в composition root/bootstrap.

Если текущая архитектура делает registry избыточным, предложить более простой вариант, но не оставлять хрупкий перебор без анализа.

### 4. Локализовать runtime type assertions

Runtime type assertions в Go допустимы, но они должны быть строго локализованы.

Нужно добиться, чтобы приведение типов:

```go
q.(*selectable.Question)
ans.(selectableanswer.Answer)
```

не было размазано по application/use case слою.

Допустимые места:

- конкретный checker;
- специализированный validator;
- фабрика/mapper ответа;
- registry/dispatcher, если он берёт на себя эту ответственность.

Недопустимо:

- делать type assertions в HTTP handler'ах;
- делать type assertions в разных use case'ах без единого правила;
- дублировать проверку совместимости в нескольких местах.

### 5. Проверить пакетную структуру question/answer

Проанализировать текущую структуру:

```txt
question/selectable
question/selectable/answer
question/matching
question/matching/answer
question/sequence
question/sequence/answer
question/short
question/short/answer
```

Нужно оценить, оправдано ли отдельное вложение `answer`.

Если оно создаёт лишнюю сложность, предложить и реализовать более простую структуру:

```txt
question/selectable
    question.go
    option.go
    answer.go
    checker.go или validator.go

question/matching
    question.go
    pair.go
    answer.go
    checker.go или validator.go
```

Но реорганизацию пакетов делать только если она реально упрощает код.

Если перенос пакетов приведёт к большому количеству изменений без пользы — оставить текущую структуру и только улучшить границы ответственности.

### 6. Улучшить attempt-пакет

Проверить пакет `domain/attempt`.

Особенно:

- создание попытки;
- snapshot вопросов;
- clone ответов;
- добавление ответа;
- проверка дедлайна;
- смена статусов;
- завершение попытки;
- восстановление из хранилища;
- инварианты агрегата.

Нужно сохранить сильные стороны:

- snapshot question через `Clone()`;
- snapshot answer через `Clone()`;
- контроль lifecycle через `Status`;
- запрет изменений после завершения/истечения/отмены;
- защиту внутренних коллекций через копирование.

Возможные улучшения:

- убрать дублирующиеся проверки;
- сделать ошибки более точными;
- улучшить читаемость методов;
- добавить тесты на edge cases;
- проверить, не нарушаются ли инварианты через `Restore`.

Не переписывать `Attempt`, если он уже корректен. Улучшения должны быть точечными.

### 7. Улучшить тесты

Добавить или обновить тесты для:

#### question

- создание валидных вопросов;
- ошибки при пустом/некорректном title;
- ограничения количества options/pairs/variants;
- `Clone()` не должен возвращать ссылку на изменяемое внутреннее состояние;
- `ChangeTitle`;
- `ChangeOptions`;
- `ChangePairs`;
- `ChangeVariants`.

#### answer

- пустой ответ;
- clone ответа;
- защита внутренних слайсов;
- некорректные ID;
- дубликаты ID, если это запрещено бизнес-правилами.

#### grading

- correct answer -> `Score(1.0)`;
- incorrect answer -> `Score(0.0)`;
- answer не того типа -> понятная ошибка;
- question не того типа -> понятная ошибка;
- отсутствующий checker -> ошибка registry;
- дубли checker'ов -> ошибка registry;
- short answer с нормализацией.

#### attempt

- создание попытки;
- создание snapshot вопросов;
- добавление ответа;
- перезапись ответа;
- добавление ответа после finish/expired/cancelled/interrupted;
- добавление ответа после deadline;
- добавление ответа на несуществующий questionID;
- finish before start;
- finish after deadline;
- restore валидного состояния;
- restore невалидного состояния.

#### application use cases

- `ValidateAnswerUseCase` должен только валидировать ответ;
- `GradeUseCase` должен валидировать и затем считать score;
- не найден вопрос;
- нет checker'а для типа вопроса;
- некорректный input;
- корректный сценарий.

### 8. Обновить документацию

Если архитектура изменилась, обновить или создать:

```txt
docs/refactoring/question-attempt-grading-analysis.md
docs/diagrams/domain-question-class-diagram.md
docs/diagrams/domain-question-sequence-diagram.md
```

Диаграммы должны быть Mermaid.

Если диаграммы уже существуют — обновить их под новую архитектуру.

## Ограничения

- Не менять бизнес-логику без необходимости.
- Не менять API контракты без явной причины.
- Не делать большой rewrite.
- Не тащить внешние библиотеки без необходимости.
- Не превращать Go-код в Java-style hierarchy.
- Не плодить интерфейсы ради интерфейсов.
- Не нарушать dependency direction.
- Домен не должен зависеть от transport, infrastructure, SQL, HTTP, JSON DTO.
- Application слой может оркестрировать, но не должен содержать бизнес-инварианты домена.
- Infrastructure не должна диктовать структуру домена.

## Критерии готовности

Задача считается выполненной, если:

1. Код компилируется.
2. Все существующие тесты проходят.
3. Добавлены новые тесты на критичные сценарии.
4. Валидация ответа отделена от оценивания.
5. `GradeUseCase` не занимается ручным перебором checker'ов.
6. Runtime type assertions локализованы.
7. Ошибки стали понятнее.
8. Пакетная структура не стала сложнее.
9. `Attempt` сохранил свои инварианты.
10. Mermaid-диаграммы обновлены, если архитектура изменилась.
11. Документ с кратким анализом изменений добавлен в `docs/refactoring`.

## Рекомендуемый порядок выполнения

1. Запустить тесты и зафиксировать текущее состояние.
2. Проанализировать `question`, `answer`, `grading`, `attempt`.
3. Написать короткий analysis document.
4. Добавить тесты на текущее поведение.
5. Ввести `AnswerValidator`.
6. Ввести `CheckerRegistry`.
7. Обновить `GradeUseCase` и `ValidateAnswerUseCase`.
8. Локализовать type assertions.
9. Улучшить ошибки.
10. Обновить тесты.
11. Обновить Mermaid-диаграммы.
12. Запустить полный test suite.
13. Проверить, что публичное поведение не сломано.

## Формат результата

После выполнения нужно показать:

1. список изменённых файлов;
2. краткое описание архитектурных изменений;
3. какие проблемы были найдены;
4. какие проблемы были исправлены;
5. какие проблемы стоит оставить на будущее;
6. результат запуска тестов;
7. обновлённые Mermaid-диаграммы или ссылки на них.

## Важное замечание

Цель задачи — не переписать проект ради красоты, а сделать кодовую базу проще, надёжнее и понятнее.

Если в каком-то месте текущая реализация уже хорошая, её нужно оставить и просто покрыть тестами.

