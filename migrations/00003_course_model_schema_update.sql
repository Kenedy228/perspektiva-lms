-- +goose Up
CREATE TABLE IF NOT EXISTS course_blocks_links (
    course_id uuid NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    block_id uuid NOT NULL REFERENCES course_blocks(id) ON DELETE RESTRICT,
    position integer NOT NULL,
    PRIMARY KEY (course_id, block_id),
    UNIQUE (course_id, position)
);

CREATE INDEX IF NOT EXISTS course_blocks_links_block_id_idx ON course_blocks_links (block_id);

WITH source_links AS (
    SELECT
        cvl.course_id,
        cvb.block_id,
        MIN(cvl.position * 100000 + cvb.position) AS sort_order
    FROM course_version_links cvl
    JOIN course_version_blocks cvb ON cvb.version_id = cvl.version_id
    GROUP BY cvl.course_id, cvb.block_id
),
ranked_links AS (
    SELECT
        course_id,
        block_id,
        ROW_NUMBER() OVER (PARTITION BY course_id ORDER BY sort_order, block_id) - 1 AS position
    FROM source_links
)
INSERT INTO course_blocks_links (course_id, block_id, position)
SELECT course_id, block_id, position
FROM ranked_links
ON CONFLICT (course_id, block_id) DO UPDATE SET position = EXCLUDED.position;

ALTER TABLE course_elements
    ADD COLUMN IF NOT EXISTS completion_mode text NOT NULL DEFAULT 'none'
    CHECK (completion_mode IN ('none', 'manual'));

UPDATE course_elements
SET completion_mode = CASE WHEN requires_read_marker THEN 'manual' ELSE 'none' END
WHERE completion_mode IS NULL OR completion_mode = '';

DROP INDEX IF EXISTS course_progress_version_id_idx;
ALTER TABLE course_progress
    DROP COLUMN IF EXISTS version_id;

-- +goose Down
ALTER TABLE course_progress
    ADD COLUMN IF NOT EXISTS version_id uuid REFERENCES course_versions(id) ON DELETE RESTRICT;

UPDATE course_progress cp
SET version_id = e.version_id
FROM enrollments e
WHERE e.id = cp.enrollment_id
  AND cp.version_id IS NULL;

ALTER TABLE course_progress
    ALTER COLUMN version_id SET NOT NULL;

CREATE INDEX IF NOT EXISTS course_progress_version_id_idx ON course_progress (version_id);

ALTER TABLE course_elements
    DROP COLUMN IF EXISTS completion_mode;

DROP INDEX IF EXISTS course_blocks_links_block_id_idx;
DROP TABLE IF EXISTS course_blocks_links;
