# Подзадача 7 — Адаптировать application use cases

## Цель
Перевести use cases на модель Course → Block → Element.

## Что менять
- backend/internal/application/usecases/course

## Удалить flow
- CreateVersion
- PublishVersion
- AddBlockToVersion

## Добавить flow
- AddBlockToCourse
- MoveCourseBlock
- AddElementToBlock
- MoveBlockElement
