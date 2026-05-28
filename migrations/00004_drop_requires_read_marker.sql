-- +goose Up
ALTER TABLE course_elements
    DROP COLUMN IF EXISTS requires_read_marker;

-- +goose Down
ALTER TABLE course_elements
    ADD COLUMN IF NOT EXISTS requires_read_marker boolean NOT NULL DEFAULT false;

UPDATE course_elements
SET requires_read_marker = CASE WHEN completion_mode = 'manual' THEN true ELSE false END
WHERE requires_read_marker IS NULL OR requires_read_marker = false;
