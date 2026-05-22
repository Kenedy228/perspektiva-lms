-- +goose Up
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE organizations (
    id uuid PRIMARY KEY,
    inn text NOT NULL UNIQUE,
    inn_type text NOT NULL CHECK (inn_type IN ('ip', 'physical', 'organization')),
    name text NOT NULL,
    deleted_at timestamptz,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX organizations_name_idx ON organizations (lower(name));

CREATE TABLE persons (
    id uuid PRIMARY KEY,
    first_name text NOT NULL,
    last_name text NOT NULL,
    middle_name text NOT NULL DEFAULT '',
    snils text UNIQUE,
    date_of_birth date,
    job_title text NOT NULL DEFAULT '',
    education text NOT NULL DEFAULT '',
    organization_id uuid REFERENCES organizations(id) ON DELETE SET NULL,
    deleted_at timestamptz,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT persons_profile_complete CHECK (
        (snils IS NULL AND date_of_birth IS NULL) OR
        (snils IS NOT NULL AND date_of_birth IS NOT NULL)
    )
);

CREATE INDEX persons_last_name_idx ON persons (lower(last_name));
CREATE INDEX persons_organization_id_idx ON persons (organization_id);

CREATE TABLE accounts (
    id uuid PRIMARY KEY,
    person_id uuid NOT NULL UNIQUE REFERENCES persons(id) ON DELETE RESTRICT,
    login text NOT NULL UNIQUE,
    password_hash text NOT NULL,
    role text NOT NULL CHECK (role IN ('admin', 'creator', 'student', 'organization')),
    status text NOT NULL CHECK (status IN ('active', 'blocked', 'deleted')),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX accounts_role_idx ON accounts (role);
CREATE INDEX accounts_status_idx ON accounts (status);

CREATE TABLE questions (
    id uuid PRIMARY KEY,
    type text NOT NULL,
    payload jsonb NOT NULL,
    deleted_at timestamptz,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX questions_type_idx ON questions (type);

CREATE TABLE question_banks (
    id uuid PRIMARY KEY,
    title text NOT NULL,
    deleted_at timestamptz,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX question_banks_title_idx ON question_banks (lower(title));

CREATE TABLE question_bank_questions (
    bank_id uuid NOT NULL REFERENCES question_banks(id) ON DELETE CASCADE,
    question_id uuid NOT NULL REFERENCES questions(id) ON DELETE RESTRICT,
    position integer NOT NULL,
    PRIMARY KEY (bank_id, question_id),
    UNIQUE (bank_id, position)
);

CREATE INDEX question_bank_questions_question_id_idx ON question_bank_questions (question_id);

CREATE TABLE quizzes (
    id uuid PRIMARY KEY,
    title text NOT NULL,
    time_limit_seconds integer NOT NULL,
    attempts_limit integer NOT NULL,
    shuffle_questions boolean NOT NULL,
    deleted_at timestamptz,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE quiz_sources (
    quiz_id uuid NOT NULL REFERENCES quizzes(id) ON DELETE CASCADE,
    bank_id uuid NOT NULL REFERENCES question_banks(id) ON DELETE RESTRICT,
    criteria_type text NOT NULL CHECK (criteria_type IN ('manual', 'random')),
    question_count integer NOT NULL CHECK (question_count > 0),
    question_ids uuid[] NOT NULL DEFAULT '{}',
    position integer NOT NULL,
    PRIMARY KEY (quiz_id, bank_id),
    UNIQUE (quiz_id, position)
);

CREATE INDEX quiz_sources_bank_id_idx ON quiz_sources (bank_id);

CREATE TABLE courses (
    id uuid PRIMARY KEY,
    title text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX courses_title_idx ON courses (lower(title));

CREATE TABLE course_versions (
    id uuid PRIMARY KEY,
    title text NOT NULL,
    status text NOT NULL CHECK (status IN ('draft', 'published', 'deleted')),
    enrollment_opens_at timestamptz,
    enrollment_closes_at timestamptz,
    copied_from_version_id uuid REFERENCES course_versions(id) ON DELETE SET NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX course_versions_status_idx ON course_versions (status);

CREATE TABLE course_version_links (
    course_id uuid NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    version_id uuid NOT NULL REFERENCES course_versions(id) ON DELETE RESTRICT,
    position integer NOT NULL,
    PRIMARY KEY (course_id, version_id),
    UNIQUE (course_id, position)
);

CREATE INDEX course_version_links_version_id_idx ON course_version_links (version_id);

CREATE TABLE course_blocks (
    id uuid PRIMARY KEY,
    title text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE course_version_blocks (
    version_id uuid NOT NULL REFERENCES course_versions(id) ON DELETE CASCADE,
    block_id uuid NOT NULL REFERENCES course_blocks(id) ON DELETE RESTRICT,
    position integer NOT NULL,
    PRIMARY KEY (version_id, block_id),
    UNIQUE (version_id, position)
);

CREATE INDEX course_version_blocks_block_id_idx ON course_version_blocks (block_id);

CREATE TABLE course_elements (
    id uuid PRIMARY KEY,
    type text NOT NULL CHECK (type IN ('text', 'slides', 'document', 'video', 'quiz')),
    title text NOT NULL,
    object_key text,
    quiz_id uuid REFERENCES quizzes(id) ON DELETE RESTRICT,
    payload jsonb NOT NULL DEFAULT '{}',
    requires_read_marker boolean NOT NULL DEFAULT false,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT course_elements_content_check CHECK (
        (type IN ('text', 'slides', 'document', 'video') AND object_key IS NOT NULL AND quiz_id IS NULL) OR
        (type = 'quiz' AND quiz_id IS NOT NULL AND object_key IS NULL)
    )
);

CREATE INDEX course_elements_quiz_id_idx ON course_elements (quiz_id);

CREATE TABLE course_block_elements (
    block_id uuid NOT NULL REFERENCES course_blocks(id) ON DELETE CASCADE,
    element_id uuid NOT NULL REFERENCES course_elements(id) ON DELETE RESTRICT,
    position integer NOT NULL,
    PRIMARY KEY (block_id, element_id),
    UNIQUE (block_id, position)
);

CREATE INDEX course_block_elements_element_id_idx ON course_block_elements (element_id);

CREATE TABLE enrollments (
    id uuid PRIMARY KEY,
    account_id uuid NOT NULL REFERENCES accounts(id) ON DELETE RESTRICT,
    course_id uuid NOT NULL REFERENCES courses(id) ON DELETE RESTRICT,
    version_id uuid NOT NULL REFERENCES course_versions(id) ON DELETE RESTRICT,
    status text NOT NULL CHECK (status IN ('inactive', 'active', 'expired')),
    enrolled_at timestamptz NOT NULL,
    completed_at timestamptz,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    UNIQUE (account_id, course_id, version_id)
);

CREATE INDEX enrollments_account_id_idx ON enrollments (account_id);
CREATE INDEX enrollments_course_id_idx ON enrollments (course_id);
CREATE INDEX enrollments_version_id_idx ON enrollments (version_id);

CREATE TABLE course_progress (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    enrollment_id uuid NOT NULL UNIQUE REFERENCES enrollments(id) ON DELETE CASCADE,
    version_id uuid NOT NULL REFERENCES course_versions(id) ON DELETE RESTRICT,
    completed_elements integer NOT NULL DEFAULT 0,
    total_elements integer NOT NULL DEFAULT 0,
    score numeric(6,2) NOT NULL DEFAULT 0,
    updated_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT course_progress_counts_check CHECK (
        completed_elements >= 0 AND total_elements >= 0 AND completed_elements <= total_elements
    )
);

CREATE INDEX course_progress_version_id_idx ON course_progress (version_id);

CREATE TABLE course_progress_markers (
    progress_id uuid NOT NULL REFERENCES course_progress(id) ON DELETE CASCADE,
    element_id uuid NOT NULL REFERENCES course_elements(id) ON DELETE RESTRICT,
    marker_type text NOT NULL CHECK (marker_type IN ('read', 'watched', 'download', 'quiz')),
    completed_at timestamptz NOT NULL,
    PRIMARY KEY (progress_id, element_id)
);

CREATE TABLE attempts (
    id uuid PRIMARY KEY,
    enrollment_id uuid NOT NULL REFERENCES enrollments(id) ON DELETE CASCADE,
    quiz_id uuid NOT NULL REFERENCES quizzes(id) ON DELETE RESTRICT,
    account_id uuid NOT NULL REFERENCES accounts(id) ON DELETE RESTRICT,
    status text NOT NULL CHECK (status IN ('finished', 'expired', 'in_progress', 'interrupted', 'cancelled')),
    started_at timestamptz NOT NULL,
    deadline_at timestamptz,
    submitted_at timestamptz,
    score numeric(6,2),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX attempts_enrollment_quiz_idx ON attempts (enrollment_id, quiz_id);
CREATE INDEX attempts_account_id_idx ON attempts (account_id);

CREATE TABLE attempt_items (
    attempt_id uuid NOT NULL REFERENCES attempts(id) ON DELETE CASCADE,
    question_id uuid NOT NULL REFERENCES questions(id) ON DELETE RESTRICT,
    position integer NOT NULL,
    PRIMARY KEY (attempt_id, question_id),
    UNIQUE (attempt_id, position)
);

CREATE TABLE attempt_answers (
    attempt_id uuid NOT NULL REFERENCES attempts(id) ON DELETE CASCADE,
    question_id uuid NOT NULL REFERENCES questions(id) ON DELETE RESTRICT,
    answer_payload jsonb NOT NULL,
    answered_at timestamptz NOT NULL,
    PRIMARY KEY (attempt_id, question_id)
);

CREATE TABLE stored_objects (
    object_key text PRIMARY KEY,
    bucket text NOT NULL,
    filename text NOT NULL,
    content_type text NOT NULL,
    size_bytes bigint NOT NULL CHECK (size_bytes >= 0),
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE audit_events (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    action text NOT NULL,
    entity_id uuid,
    actor_role text NOT NULL,
    payload jsonb NOT NULL DEFAULT '{}',
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX audit_events_action_idx ON audit_events (action);
CREATE INDEX audit_events_entity_id_idx ON audit_events (entity_id);
CREATE INDEX audit_events_created_at_idx ON audit_events (created_at);

-- +goose Down
DROP TABLE IF EXISTS audit_events;
DROP TABLE IF EXISTS stored_objects;
DROP TABLE IF EXISTS attempt_answers;
DROP TABLE IF EXISTS attempt_items;
DROP TABLE IF EXISTS attempts;
DROP TABLE IF EXISTS course_progress_markers;
DROP TABLE IF EXISTS course_progress;
DROP TABLE IF EXISTS enrollments;
DROP TABLE IF EXISTS course_block_elements;
DROP TABLE IF EXISTS course_elements;
DROP TABLE IF EXISTS course_version_blocks;
DROP TABLE IF EXISTS course_blocks;
DROP TABLE IF EXISTS course_version_links;
DROP TABLE IF EXISTS course_versions;
DROP TABLE IF EXISTS courses;
DROP TABLE IF EXISTS quiz_sources;
DROP TABLE IF EXISTS quizzes;
DROP TABLE IF EXISTS question_bank_questions;
DROP TABLE IF EXISTS question_banks;
DROP TABLE IF EXISTS questions;
DROP TABLE IF EXISTS accounts;
DROP TABLE IF EXISTS persons;
DROP TABLE IF EXISTS organizations;
