# Финальный рефакторинг: question / grading / attempt — итоги

Дата: 2026-05-28

## 1. Что было сделано (общий итог двух итераций)

### Итерация 1 (001-task)
- Введён `AnswerValidator` interface + 4 реализации
- Введён `CheckerRegistry` (map-based, O(1) lookup)
- `GradeUseCase` обновлён: использует Registry + validators
- `ValidateAnswerUseCase` обновлён: использует AnswerValidator напрямую (не GradeUseCase)
- Унифицированы sentinel errors в `grading/errors.go`
- Удалены дублирующиеся `errors.go` в 4 пакетах
- Добавлены тесты: `validator_test.go` (4 шт), `registry_test.go`
- Обновлены диаграммы

### Итерация 2 (002-task — текущая)
- **Удалён `Supports()`** из `Checker` interface и всех 4 реализаций
- Удалены `TestChecker_Supports` из всех 4 `checker_test.go`
- Обновлён `registry_test.go` — убран `Supports()` из `stubChecker`
- Обновлён `fakeChecker` в `grading_test.go` — убран `Supports()`
- Исправлен баг в class-diagram (дублирующийся Checker)
- Обновлены Mermaid-диаграммы

## 2. Исправленные проблемы

| # | Проблема | Решение | Статус |
|---|----------|---------|--------|
| P1 | GradeUseCase перебирал checker'ы через `Supports()` | `Registry.Get()` — O(1) map lookup | ✅ |
| P2 | `Supports()` создавал лишний runtime dispatch | Удалён из интерфейса и реализаций | ✅ |
| P3 | Validation и grading смешаны | `AnswerValidator` отделён от `Checker` | ✅ |
| P4 | `ValidateAnswerUseCase` использовал `GradeUseCase` | Использует `AnswerValidator` напрямую | ✅ |
| P5 | Type assertions в application-слое | Локализованы в domain (Validator, Checker) | ✅ |
| P6 | Дублирование sentinel errors | Унифицированы в `grading/errors.go` | ✅ |
| P7 | Дублирующийся Checker в class-diagram | Удалён | ✅ |

## 3. Оставленные компромиссы

### Структурная валидация ID (P5 из анализа)
**Решение:** оставлено на будущее. Проверка существования option/prompt/match ID в ответе студента не реализована. Текущее поведение: несуществующие ID молча не матчатся с правильными ответами → score 0. Это не баг, а conscious design choice. Добавление потребует:
- `AnswerValidator.Validate()` должен принимать конкретные типы и сверять ID
- Каждый validator должен иметь доступ к структуре вопроса (options/pairs)
- Либо validator должен получать конкретные типы (а не интерфейсы)

### Пакетная структура question/*/answer
**Решение:** оставлена без изменений. Структура `question/{type}/answer/` оправдана:
- Высокая cohesion (2-5 файлов на пакет)
- Предотвращает циклические импорты
- Соответствует Go-идиоме «маленькие пакеты»

### Checker interface без type identification
**Решение:** после удаления `Supports()`, Checker не имеет метода идентификации типа. Это осознанное решение:
- Тип задаётся явно через ключ map при создании Registry
- Не нужен ни `Supports()`, ни `Type()` — идентификация вынесена в composition root
- Минимальный интерфейс: только `Check()`

## 4. Архитектурный flow (после рефакторинга)

```
ValidateAnswerUseCase
    → Repository.FindByID()
    → validators[Type].Validate(q, a)    // только валидация
    → return error | nil

GradeUseCase
    → Repository.FindByID()
    → validators[Type].Validate(q, a)    // валидация
    → Registry.Get(Type)                 // O(1) lookup checker'а
    → Checker.Check(q, a)                // вычисление Score
    → return GradeOutput{Score}
```

```
Registry
    → New(map[Type]Checker)             // composition root
    → Get(Type) (Checker, error)        // O(1), детерминированный
```

```
Checker (минимальный интерфейс)
    → Check(Question, Answer) (Score, error)  // только grading
                                              // + defensive type assertions
```

## 5. Что НЕ было сделано (и почему)

| Идея | Почему отклонено |
|------|-----------------|
| Добавить `Type()` на Checker | Избыточно: тип задаётся ключом map в composition root |
| ValidatorRegistry | Избыточно: map[string]Validator достаточно, нет отдельной логики |
| Visitor pattern для grading | Overengineering: 4 типа, явный switch/match проще |
| Объединить Validator + Checker в один интерфейс | Нарушает разделение ответственности (Single Responsibility) |
| Вынести answer из question/*/answer | Создаст циклические импорты или раздует пакеты |
| Reflection-based dispatch | Запрещено явно: недетерминированно, медленно, нечитаемо |

## 6. Результат тестов

```
25 пакетов → PASS
- domain/grading/matching
- domain/grading/registry  
- domain/grading/score
- domain/grading/selectable
- domain/grading/sequence
- domain/grading/short
- domain/question (и все подпакеты)
- domain/attempt (и все подпакеты)
- application/usecases/question/grading
```
