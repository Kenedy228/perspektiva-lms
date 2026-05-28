-- +goose Up
UPDATE course_elements
SET type = 'document'
WHERE type IN ('text', 'slides');

ALTER TABLE course_elements
    DROP CONSTRAINT IF EXISTS course_elements_content_check;

ALTER TABLE course_elements
    ADD CONSTRAINT course_elements_content_check CHECK (
        (type IN ('document', 'video') AND object_key IS NOT NULL AND quiz_id IS NULL) OR
        (type = 'quiz' AND quiz_id IS NOT NULL AND object_key IS NULL)
    ) NOT VALID;

ALTER TABLE course_elements
    VALIDATE CONSTRAINT course_elements_content_check;

DO $$
DECLARE
    con_name text;
BEGIN
    SELECT conname INTO con_name
    FROM pg_constraint
    WHERE conrelid = 'course_elements'::regclass
      AND contype = 'c'
      AND conname <> 'course_elements_content_check'
      AND conname <> 'course_elements_completion_mode_check'
      AND pg_get_constraintdef(oid) LIKE '%type%IN%text%';

    IF con_name IS NOT NULL THEN
        EXECUTE 'ALTER TABLE course_elements DROP CONSTRAINT ' || con_name;
    END IF;
END $$;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conrelid = 'course_elements'::regclass
          AND contype = 'c'
          AND conname = 'course_elements_type_check'
    ) THEN
        ALTER TABLE course_elements
            ADD CONSTRAINT course_elements_type_check
            CHECK (type IN ('document', 'video', 'quiz')) NOT VALID;
        ALTER TABLE course_elements
            VALIDATE CONSTRAINT course_elements_type_check;
    END IF;
END $$;

-- +goose Down
ALTER TABLE course_elements
    DROP CONSTRAINT IF EXISTS course_elements_content_check;

ALTER TABLE course_elements
    ADD CONSTRAINT course_elements_content_check CHECK (
        (type IN ('text', 'slides', 'document', 'video') AND object_key IS NOT NULL AND quiz_id IS NULL) OR
        (type = 'quiz' AND quiz_id IS NOT NULL AND object_key IS NULL)
    ) NOT VALID;

ALTER TABLE course_elements
    VALIDATE CONSTRAINT course_elements_content_check;

ALTER TABLE course_elements
    DROP CONSTRAINT IF EXISTS course_elements_type_check;

ALTER TABLE course_elements
    ADD CONSTRAINT course_elements_type_check
    CHECK (type IN ('text', 'slides', 'document', 'video', 'quiz')) NOT VALID;

ALTER TABLE course_elements
    VALIDATE CONSTRAINT course_elements_type_check;
