-- +goose Up

-- stored_objects не используется: метаданные файлов хранятся в course_elements.
-- Удаляем таблицу, чтобы не создавать ложное ощущение наличия инвентаря.
DROP TABLE IF EXISTS stored_objects;

-- Восстанавливаем CHECK-ограничение на course_progress,
-- снятое в миграции 00006 (инвариант защищается доменом, но лучше иметь и DB-уровень).
ALTER TABLE course_progress
    ADD CONSTRAINT course_progress_counts_check
    CHECK (completed_elements >= 0 AND total_elements >= 0 AND completed_elements <= total_elements)
    NOT VALID;

ALTER TABLE course_progress VALIDATE CONSTRAINT course_progress_counts_check;

-- +goose Down

ALTER TABLE course_progress
    DROP CONSTRAINT IF EXISTS course_progress_counts_check;

CREATE TABLE IF NOT EXISTS stored_objects (
    object_key   text PRIMARY KEY,
    bucket       text NOT NULL,
    filename     text NOT NULL,
    content_type text NOT NULL,
    size_bytes   bigint NOT NULL CHECK (size_bytes >= 0),
    created_at   timestamptz NOT NULL DEFAULT now()
);
