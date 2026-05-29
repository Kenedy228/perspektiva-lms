---
name: project-critical-fixes-done
description: Критические задачи бэкенда исправлены — MinIO аутентификация и AttemptPolicy SQL-баг
metadata:
  type: project
---

Оба критических бага закрыты (2026-05-29).

**Исправление 1 — AttemptPolicy.CanStartQuiz**
Файл: `backend/internal/infrastructure/postgres/attempt.go:131`
Старый SQL ссылался на `e.version_id` из `enrollments` (удалён миграцией 00007). Заменён на JOIN через `course_blocks_links` → `course_block_elements` → `course_elements`.
**Why:** Приводило к SQL-ошибке при любой попытке студента начать тест.

**Исправление 2 — MinIO аутентификация**
Файл: `backend/internal/infrastructure/storage/minio/`
Заменён самописный HTTP-адаптер без авторизации на официальный SDK `github.com/minio/minio-go/v7`.
Config расширен: `AccessKey`, `SecretKey`, `UseSSL` (env: `MINIO_ACCESS_KEY`, `MINIO_SECRET_KEY`, `MINIO_USE_SSL`).
`GetDownloadURL` теперь генерирует presigned URL через `PresignedGetObject` с TTL.
**Why:** Без подписи S3 все PUT/DELETE в production MinIO возвращали 403.

**Как применять:**
- Поставить env-переменные `MINIO_ACCESS_KEY`, `MINIO_SECRET_KEY` в deploy-конфиге.
- `MINIO_ENDPOINT` теперь без схемы (`localhost:9000`, не `http://localhost:9000`) — minio-go принимает только хост:порт.
