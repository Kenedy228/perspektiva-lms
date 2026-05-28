-- +goose Up

-- Drop index and column from enrollments
ALTER TABLE enrollments DROP CONSTRAINT IF EXISTS enrollments_account_id_course_id_version_id_key;
DROP INDEX IF EXISTS enrollments_version_id_idx;
ALTER TABLE enrollments DROP COLUMN IF EXISTS version_id;
ALTER TABLE enrollments ADD CONSTRAINT enrollments_account_id_course_id_key UNIQUE (account_id, course_id);

-- Drop index and column from course_progress
DROP INDEX IF EXISTS course_progress_version_id_idx;
ALTER TABLE course_progress DROP COLUMN IF EXISTS version_id;

-- +goose Down

ALTER TABLE course_progress ADD COLUMN version_id uuid REFERENCES course_versions(id) ON DELETE RESTRICT;
CREATE INDEX course_progress_version_id_idx ON course_progress (version_id);

ALTER TABLE enrollments DROP CONSTRAINT IF EXISTS enrollments_account_id_course_id_key;
ALTER TABLE enrollments ADD COLUMN version_id uuid REFERENCES course_versions(id) ON DELETE RESTRICT;
ALTER TABLE enrollments ALTER COLUMN version_id SET NOT NULL;
CREATE INDEX enrollments_version_id_idx ON enrollments (version_id);
ALTER TABLE enrollments ADD CONSTRAINT enrollments_account_id_course_id_version_id_key UNIQUE (account_id, course_id, version_id);
