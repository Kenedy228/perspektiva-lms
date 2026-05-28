# Задача 1 — Только Course Domain + Content Model (без глобального рефакторинга)

## ВАЖНО

Это НЕ задача на рефакторинг всего проекта.

Это НЕ задача на migration всей LMS.

Это локальная задача только для:

```txt
backend/internal/domain/course
```

и минимально необходимого кода вокруг него.

==================================================
ЖЁСТКИЙ SCOPE
=============

РАЗРЕШЕНО менять ТОЛЬКО:

```txt
backend/internal/domain/course
backend/internal/domain/course/block
backend/internal/domain/course/element
backend/internal/domain/course/progress
backend/internal/domain/course/title
```

Дополнительно РАЗРЕШЕНО менять ТОЛЬКО если проект не компилируется:

```txt
backend/internal/application/usecases/course/commands/helpers.go
```

И всё.

==================================================
ЗАПРЕЩЕНО МЕНЯТЬ
================

НЕ трогать:

```txt
application/usecases/*
кроме helpers.go

postgres repositories
transport/http
migrations
ports
bank
quiz
question
grading
attempt
auth
user
enrollment
session
main.go
frontend
docs
diagrams
DTO
routes
```

НЕ создавать:

* новые сервисы;
* новые repository;
* новые подсистемы;
* новую архитектуру;
* compatibility layers;
* generic abstractions.

==================================================
ЦЕЛЬ ЗАДАЧИ
===========

Нужно привести ТОЛЬКО domain/course к консистентному состоянию.

После выполнения должно быть:

```txt
Course → Block → Element
```

без Version-flow внутри домена.

==================================================
ЧТО ИМЕННО НУЖНО СДЕЛАТЬ
========================

# 1. Course aggregate

Файл:

```txt
backend/internal/domain/course/course.go
```

## Нужно:

Удалить методы:

```go
VersionIDs()
AddVersionID()
RemoveVersionID()
HasVersion()
```

И удалить поле:

```go
versionIDs
```

## Должно остаться:

```go
blockIDs
```

И методы:

```go
BlockIDs()
AddBlockID()
RemoveBlockID()
MoveBlock()
```

## Важно

НЕ менять public API других модулей.

НЕ добавлять новые зависимости.

==================================================

# 2. Progress domain

Папка:

```txt
backend/internal/domain/course/progress
```

## Нужно:

Удалить из домена:

```go
versionID
VersionID()
```

И убрать `versionID` из:

```go
New(...)
Restore(...)
```

## Целевая модель:

```txt
Progress
  - enrollmentID
  - completed element markers
```

## Реализовать:

```go
MarkCompleted(elementID, at)
UnmarkCompleted(elementID)
IsCompleted(elementID)
CompletedCount()
Percent()
```

## Важно

Progress:

* не должен знать Course;
* не должен знать Block;
* не должен знать repository;
* работает только через elementID.

==================================================

# 3. Element content model

Папка:

```txt
backend/internal/domain/course/element
```

## Цель

Оставить только 3 content type:

```txt
test
lecture_material
download_file
```

==================================================

## 3.1 Test content

Создать/исправить пакет:

```txt
content/test
```

Content должен:

* хранить `quizID`;
* запрещать `uuid.Nil`;
* реализовывать общий `element.Content`;
* поддерживать Clone().

==================================================

## 3.2 LectureMaterial content

Создать/исправить пакет:

```txt
content/lecturematerial
```

Поддерживаемые kind:

```txt
video
pdf
```

Content должен:

* хранить file.File;
* запрещать unsupported kind;
* поддерживать Clone().

==================================================

## 3.3 DownloadFile content

Создать/исправить пакет:

```txt
content/downloadfile
```

Content должен:

* хранить file.File;
* поддерживать Clone().

==================================================

# 4. Старые content packages

Найти:

```txt
attachment
quiz
slides
text
video
```

## Разрешено:

* либо удалить;
* либо оставить deprecated;
* но НЕ использовать в active domain flow.

==================================================

# 5. helpers.go

Файл:

```txt
backend/internal/application/usecases/course/commands/helpers.go
```

## Разрешено менять ТОЛЬКО:

* imports;
* вызовы новых content packages;

## Запрещено:

* менять архитектуру usecases;
* создавать новые сценарии;
* менять unrelated logic.

==================================================

# 6. Тесты

Разрешено менять ТОЛЬКО тесты внутри:

```txt
backend/internal/domain/course
```

Нужно покрыть:

```txt
Course block flow
Progress without versionID
Test content
LectureMaterial content
DownloadFile content
Clone()
MoveBlock()
MoveElement()
Completion tracking
```

==================================================
ЧТО НЕ НУЖНО ДЕЛАТЬ
===================

НЕ делать:

```txt
repository refactoring
postgres refactoring
HTTP routes
DTO updates
migration changes
application architecture changes
new services
new modules
global cleanup
```

НЕ пытаться:

* закончить весь Course refactoring;
* обновить весь проект;
* мигрировать всю БД.

==================================================
КРИТЕРИЙ УСПЕХА
===============

Задача выполнена, если:

1. domain/course больше не содержит Version methods.
2. Progress больше не содержит versionID.
3. Content model консистентна.
4. helpers.go компилируется.
5. Course domain tests проходят.
6. Никакие unrelated modules не изменены.

==================================================
ФИНАЛЬНЫЙ ОТЧЁТ
===============

После выполнения показать:

1. список изменённых файлов;
2. какие файлы НЕ трогались;
3. какие Version-зависимости удалены;
4. какие старые content packages остались;
5. результат domain tests;
6. что осталось на будущее.
