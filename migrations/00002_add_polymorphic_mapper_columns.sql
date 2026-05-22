-- +goose Up
ALTER TABLE course_elements
    ADD COLUMN IF NOT EXISTS payload jsonb NOT NULL DEFAULT '{}';

ALTER TABLE attempts
    ADD COLUMN IF NOT EXISTS deadline_at timestamptz;

ALTER TABLE attempt_answers
    ADD COLUMN IF NOT EXISTS answered_at timestamptz NOT NULL DEFAULT now();

-- +goose Down
ALTER TABLE attempt_answers
    DROP COLUMN IF EXISTS answered_at;

ALTER TABLE attempts
    DROP COLUMN IF EXISTS deadline_at;

ALTER TABLE course_elements
    DROP COLUMN IF EXISTS payload;
