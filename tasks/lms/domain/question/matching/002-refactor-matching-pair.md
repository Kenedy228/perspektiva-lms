# Рефакторинг пакета pair

## Контекст

Необходимо провести рефакторинг пакета `matching/pair`.

Текущая реализация использует тип `text.Text`.
Необходимо отказаться от него и использовать обычный `string`.

Также необходимо вынести валидацию в отдельный validator-файл и улучшить структуру пакета.

Рефакторинг должен быть локальным, без изменения бизнес-логики, кроме явно описанных изменений.

---

# Цель 1. Избавиться от text.Text

Необходимо полностью отказаться от использования типа:

```go
text.Text
```

Вместо него должен использоваться:

```go
string
```

---

## Задача 1. Изменить структуру Prompt

Текущая структура:

```go
type Prompt struct {
	id   uuid.UUID
	text text.Text
}
```

Должна принять следующий вид:

```go
type Prompt struct {
	id    uuid.UUID
	value string
}
```

Требования:

* поле `text` переименовать в `value`;
* тип `text.Text` заменить на `string`.

---

## Задача 2. Обновить конструкторы

Необходимо обновить:

```go
NewPrompt
RestorePrompt
```

Требования:

* убрать зависимость от `text.Text`;
* принимать `string`;
* использовать новый validator;
* сохранить существующую бизнес-логику.

---

## Задача 3. Обновить методы Prompt

Необходимо обновить:

* `Text()`
* `IsIncomplete()`

Требования:

* `Text()` должен быть переименован в `Value()`;
* `Value()` должен возвращать `string`;
* метод `IsIncomplete` переименовать в `IsZero`;
* логика метода должна сохраниться.

Новая логика `IsZero`:

```go
return p.id == uuid.Nil || len(p.value) == 0
```

---

# Цель 2. Добавить validator

Необходимо создать отдельный файл:

```text
validator.go
```

В validator необходимо вынести валидацию значения.

---

## Задача 1. Добавить validateValue

Необходимо реализовать функцию:

```go
func validateValue(value string) error {
	if err := validateRequired(value); err != nil {
		return err
	}

	if err := validateCharsLimit(value); err != nil {
		return err
	}

	return nil
}
```

---

## Задача 2. Добавить validateRequired

Необходимо реализовать функцию:

```go
func validateRequired(value string) error {
	if value == "" {
		return fmt.Errorf("%w: значение не может быть пустым", ErrInvalid)
	}

	return nil
}
```

---

## Задача 3. Добавить validateCharsLimit

Необходимо реализовать функцию:

```go
func validateCharsLimit(value string) error {
	rc := utf8.RuneCountInString(value)

	if rc > ValueCharsLimit {
		return fmt.Errorf(
			"%w: значение не может прешывать %d символов (текущее количество символов - %d)",
			ErrInvalid,
			ValueCharsLimit,
			rc,
		)
	}

	return nil
}
```

---

## Задача 4. Добавить validateIDRequired

Необходимо реализовать функцию:

```go
func validateIDRequired(id uuid.UUID) error
```

Логика:

* если `id == uuid.Nil`, возвращаем ошибку;
* иначе возвращаем `nil`.

---

# Цель 3. Вынести константы

Необходимо создать отдельный файл:

```text
consts.go
```

Требования:

* вынести константы ограничения количества символов;
* не хранить константы внутри domain-моделей.

Пример:

```go
const (
	ValueCharsLimit  int = 255
	PromptCharsLimit int = 255
)
```

---

# Цель 4. Улучшить ошибки

После выполнения рефакторинга необходимо:

* русифицировать ошибки;
* добавить более подробный контекст;
* использовать wrapping через `%w`.

Пример ожидаемого стиля:

```go
fmt.Errorf(
	"%w: значение не может превышать %d символов (текущее количество символов - %d)",
	ErrInvalid,
	ValueCharsLimit,
	rc,
)
```

---

# Цель 5. Добавить тесты

Необходимо добавить тесты:

* на экспортируемые методы;
* на validator-функции;
* на edge-case сценарии;
* на ошибки;
* на граничные значения.

Особое внимание уделить:

* пустым строкам;
* `uuid.Nil`;
* превышению лимита символов;
* unicode-символам;
* zero-state объектов.

Требования:

* использовать table-driven tests;
* избегать дублирования;
* тестировать не только happy-path.

---

# Цель 6. Добавить godoc

Необходимо добавить godoc-комментарии на русском языке для:

* публичных структур;
* публичных методов;
* публичных функций;
* экспортируемых ошибок.

Требования:

* комментарии должны начинаться с имени сущности;
* комментарии должны описывать поведение и назначение.

---

---

# Цель 7. Переименовать IsIncomplete в Pair и validator-функциях

Необходимо завершить рефакторинг zero-state логики внутри пакета `pair`.

---

## Задача 1. Переименовать Pair.IsIncomplete в Pair.IsZero

Текущий метод:

```go
func (p Pair) IsIncomplete() bool {
	return p.prompt.IsIncomplete() || p.match.IsIncomplete()
}
```

Необходимо:

* переименовать `IsIncomplete` в `IsZero`;
* сохранить текущую бизнес-логику;
* использовать новые методы `Prompt.IsZero` и `Match.IsZero`.

Новая реализация должна иметь следующий вид:

```go
func (p Pair) IsZero() bool {
	return p.prompt.IsZero() || p.match.IsZero()
}
```

---

## Задача 2. Обновить validatePrompt

Текущая реализация:

```go
func validatePrompt(prompt Prompt) error {
	if prompt.IsIncomplete() {
		return fmt.Errorf("%w: invalid value", ErrInvalid)
	}

	return nil
}
```

Необходимо:

* заменить использование `IsIncomplete`;
* использовать новый метод `IsZero`;
* улучшить текст ошибки.

Ожидаемая логика:

```go
func validatePrompt(prompt Prompt) error {
	if prompt.IsZero() {
		return fmt.Errorf("%w: prompt не должен быть zero-value", ErrInvalid)
	}

	return nil
}
```

---

## Задача 3. Обновить validateMatch

Текущая реализация:

```go
func validateMatch(match Match) error {
	if match.IsIncomplete() {
		return fmt.Errorf("%w: invalid value", ErrInvalid)
	}

	return nil
}
```

Необходимо:

* заменить использование `IsIncomplete`;
* использовать новый метод `IsZero`;
* улучшить текст ошибки.

Ожидаемая логика:

```go id="n7d4lb"
func validateMatch(match Match) error {
	if match.IsZero() {
		return fmt.Errorf("%w: match не должен быть zero-value", ErrInvalid)
	}

	return nil
}
```

---

## Требования

* Внутри пакета больше не должно использоваться название `IsIncomplete`.
* Необходимо использовать единообразный подход через `IsZero`.
* Сохранить существующую бизнес-логику.
* Добавить/обновить тесты на новую zero-state логику.

---

# Ограничения

* Не изменять бизнес-логику, кроме явно описанных изменений.
* Не добавлять новые абстракции без необходимости.
* Не использовать panic.
* Не изменять unrelated-код.
* Сохранить существующий стиль проекта.
* Не добавлять лишние интерфейсы.

---

# Проверка результата

После выполнения изменений необходимо выполнить:

```bash
go test ./...
go vet ./...
```

---

# Ожидаемый результат

В результате выполнения задачи:

* пакет больше не использует `text.Text`;
* используется `string`;
* валидация вынесена в validator;
* улучшены ошибки;
* добавлены тесты;
* добавлен godoc;
* проект успешно проходит тесты.
