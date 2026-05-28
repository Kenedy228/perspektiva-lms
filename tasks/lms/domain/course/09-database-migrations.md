# Подзадача 9 — Миграции БД

## Цель
Обновить schema под новую course-модель.

## Что менять
- backend/migrations

## Нужно
- убрать active usage course_versions
- добавить positions
- добавить completion_mode
- убрать version_id из progress
