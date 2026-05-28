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

# Задача 5 — Обновить ports и repositories

## Цель

Убрать version-centric flow из ports и PostgreSQL repositories.

## Что менять

```txt
backend/internal/application/ports
backend/internal/infrastructure/postgres
```

## Что удалить

- CanEnrollVersion
- FindVersion
- SaveVersion
- version_id в ProgressRepository
- course_versions joins
- course_version_links joins

## Новая модель

```txt
courses
  → blocks
      → elements
```

## Progress

```txt
course/enrollment
  → progress
      → completed elements
```

## Read models

Убрать:
- VersionID
- Versions

Добавить:
- Blocks
- Elements
- Progress

## Тесты

Проверить:
- save/load course with blocks
- preserve order
- progress without version

## Критерии готовности

- Repository active flow больше не version-centric
