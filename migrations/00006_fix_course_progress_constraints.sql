-- +goose Up
ALTER TABLE course_progress
    DROP CONSTRAINT IF EXISTS course_progress_counts_check;

ALTER TABLE course_progress
    DROP COLUMN IF EXISTS score;

-- +goose Down
ALTER TABLE course_progress
    ADD COLUMN IF NOT EXISTS score numeric(6,2) NOT NULL DEFAULT 0;

ALTER TABLE course_progress
    ADD CONSTRAINT course_progress_counts_check CHECK (
        completed_elements >= 0 AND total_elements >= 0 AND completed_elements <= total_elements
    ) NOT VALID;

ALTER TABLE course_progress
    VALIDATE CONSTRAINT course_progress_counts_check;
