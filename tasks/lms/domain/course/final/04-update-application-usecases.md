Прочитай задачу и выполни её строго в ограниченном scope.

ВАЖНО: НЕ рефакторить весь проект.
ВАЖНО: НЕ создавать новые глобальные подсистемы.
ВАЖНО: НЕ трогать unrelated modules.

==================================================
SCOPE ЗАДАЧИ
==================================================

Работать только с модулем Course и минимально необходимыми связанными слоями:

РАЗРЕШЕНО менять:

1. Domain:
- backend/internal/domain/course
- backend/internal/domain/course/block
- backend/internal/domain/course/element
- backend/internal/domain/course/progress
- backend/internal/domain/course/title

2. Application слой ТОЛЬКО для Course:
- backend/internal/application/usecases/course
- application ports, если они напрямую связаны с Course

3. Infrastructure ТОЛЬКО repository/mapping для Course:
- postgres course repository/query service/progress repository

4. Transport ТОЛЬКО Course routes/DTO:
- handlers/routes, которые напрямую используют Course/Block/Element/Progress

5. Migrations ТОЛЬКО если без этого проект не собирается или active Course flow невозможен.

==================================================
ЗАПРЕЩЕНО МЕНЯТЬ
==================================================

НЕ трогать:

- question
- grading
- attempt
- bank
- quiz
- enrollment, кроме прямой связи с Course progress
- auth
- user/account
- session
- middleware
- config
- bootstrap/main.go, кроме минимальной правки зависимостей Course
- frontend
- unrelated repositories
- unrelated DTO
- unrelated migrations

==================================================
ЦЕЛЬ
==================================================

Довести Course-модуль до модели:

Course → Block → Element

Без CourseVersion в active domain flow.

==================================================
ЧТО ИМЕННО СДЕЛАТЬ
==================================================

1. В domain/course убрать остатки Version-flow:
- VersionIDs
- AddVersionID
- RemoveVersionID
- HasVersion

2. Оставить только block-based API:
- BlockIDs
- AddBlockID
- RemoveBlockID
- MoveBlock

3. В domain/course/progress убрать versionID.

4. Progress должен работать по completed element IDs.

5. В domain/course/element привести content model к:
- test
- lecture_material
- download_file

6. Проверить, что content-пакеты и imports совпадают.

7. В application/usecases/course убрать только CourseVersion-specific сценарии.

8. В postgres Course mapping убрать active dependency от course_versions только в Course-related коде.

9. Добавить/обновить тесты только для Course module.

==================================================
НЕ СОЗДАВАТЬ
==================================================

Не создавать:

- новые глобальные сервисы;
- универсальный file module;
- generic content management system;
- event bus;
- service locator;
- новые bounded contexts;
- новые unrelated repository;
- новую архитектуру всего проекта.

==================================================
КРИТЕРИЙ УСПЕХА
==================================================

Задача выполнена, если:

1. Course domain больше не содержит Version methods.
2. Progress не содержит versionID.
3. Element content model консистентна.
4. Course usecases не требуют versionID.
5. Course repository active flow не зависит от course_versions.
6. Unrelated modules не изменены.
7. Тесты Course-модуля проходят.
8. Если возможно, проходит go test ./...

==================================================
ФИНАЛЬНЫЙ ОТЧЁТ
==================================================

После выполнения показать:

1. какие файлы изменены;
2. какие файлы НЕ трогались;
3. какие Version-зависимости удалены;
4. какие Course-зависимости остались;
5. результат тестов;
6. что осталось на будущее.

# Задача 4 — Обновить application use cases

## Цель

Перевести application layer на модель:

```txt
Course → Block → Element
```

## Что удалить

- CreateVersion
- PublishVersion
- DeleteVersion
- AddBlockToVersion
- MoveBlockInVersion

## Active use cases

### Course
- CreateCourse
- RenameCourse
- AddBlockToCourse
- RemoveBlockFromCourse
- MoveCourseBlock

### Block
- AddElementToBlock
- RemoveElementFromBlock
- MoveBlockElement

### Element
- CreateTestElement
- CreateLectureMaterialElement
- CreateDownloadFileElement
- ChangeElementCompletionMode

### Progress
- MarkElementCompleted
- UnmarkElementCompleted
- GetCourseProgress

## Правила

Use case:
- не содержит SQL
- не содержит HTTP logic
- использует domain methods
- сохраняет через repository

## Тесты

Добавить:
- AddBlockToCourse
- MoveCourseBlock
- AddElementToBlock
- MoveBlockElement
- MarkElementCompleted

## Критерии готовности

- Active use cases не используют version
# TODO — Доработка задачи 4 (Course application use cases)

## Статус

Задача 4 выполнена частично.

Базовый Course → Block → Element flow уже реализован и работает:

* AddBlockToCourseUseCase
* MoveCourseBlockUseCase
* AddElementToBlockUseCase
* MoveBlockElementUseCase
* MarkProgressUseCase

Architecture boundaries сохранены:

* use cases не содержат SQL;
* use cases не содержат HTTP logic;
* domain orchestration корректная;
* unrelated modules не затронуты.

Но часть planned use cases ещё отсутствует.

---

# Что осталось доделать

## 1. RemoveBlockFromCourseUseCase

### Нужно реализовать

```txt
RemoveBlockFromCourseUseCase
```

### Поведение

Use case должен:

* загрузить Course;
* удалить blockID через domain method;
* сохранить Course.

### Не делать

* Не удалять сам Block aggregate автоматически.
* Не делать cascading delete logic.
* Не добавлять repository cleanup orchestration.

---

## 2. RemoveElementFromBlockUseCase

### Нужно реализовать

```txt
RemoveElementFromBlockUseCase
```

### Поведение

Use case должен:

* загрузить Block;
* удалить elementID через domain method;
* сохранить Block.

### Не делать

* Не удалять сам Element aggregate автоматически.
* Не делать cascade cleanup.

---

## 3. ChangeElementCompletionModeUseCase

### Нужно реализовать

```txt
ChangeElementCompletionModeUseCase
```

### Поведение

Use case должен:

* загрузить Element;
* вызвать:

  ```go
  ChangeCompletionMode(...)
  ```
* сохранить Element.

---

## 4. UnmarkElementCompletedUseCase

### Нужно реализовать

```txt
UnmarkElementCompletedUseCase
```

### Поведение

Use case должен:

* загрузить Progress;
* вызвать:

  ```go
  UnmarkCompleted(...)
  ```
* сохранить Progress.

---

## 5. GetCourseProgress query/use case

### Нужно реализовать

Минимальный read/query flow:

```txt
GetCourseProgress
```

### Должен возвращать

* completed count;
* percent;
* completed element IDs.

### Важно

* Не тащить HTTP DTO в application layer.
* Не добавлять CQRS framework.
* Не делать query bus.

---

# Решение по Create*Element use cases

## Текущее состояние

Сейчас используется:

```txt
AddElementToBlockUseCase
```

через:

```txt
ElementContentInput
```

Это допустимо.

## Решение

НЕ нужно сейчас создавать:

* CreateTestElementUseCase
* CreateLectureMaterialElementUseCase
* CreateDownloadFileElementUseCase

Текущий generic AddElementToBlockUseCase достаточно хороший для текущего проекта.

НЕ усложнять architecture без необходимости.

---

# Решение по LectureMaterial kind

Проверить:

```txt
lecturematerialcontent.New(...)
```

Если kind уже определяется:

* по extension;
* или по file type;

то оставляем как есть.

НЕ нужно тащить explicit enum через весь application layer, если текущая реализация уже стабильна.

---

# Что НЕ делать

НЕ:

* переписывать application architecture;
* вводить query bus;
* вводить command bus;
* вводить service layer;
* плодить interfaces;
* делать отдельные handlers/services/facades;
* трогать unrelated modules.

---

# Scope

Разрешено менять только:

```txt
backend/internal/application/usecases/course
backend/internal/domain/course
```

И минимально:

* application ports;
* repository interfaces;
  если это действительно нужно.

---

# Критерий готовности

Задача считается завершённой, если:

1. RemoveBlockFromCourseUseCase реализован.
2. RemoveElementFromBlockUseCase реализован.
3. ChangeElementCompletionModeUseCase реализован.
4. UnmarkElementCompletedUseCase реализован.
5. GetCourseProgress реализован.
6. Тесты на эти сценарии добавлены.
7. Architecture не усложнилась.
8. Generic AddElementToBlockUseCase сохранён.
9. Unrelated modules не изменены.

