# PostgreSQL Schema Overview

This document describes the planned PostgreSQL objects for the LMS backend. Migrations are expected to be managed with `goose`.

## Goose

Installed version checked locally:

```bash
goose version: v3.26.0
```

Recommended migration layout:

```text
migrations/
  00001_create_core_identity.sql
  00002_create_question_banks.sql
  00003_create_questions.sql
  00004_create_quizzes.sql
  00005_create_courses.sql
  00006_create_enrollments_and_progress.sql
  00007_create_attempts.sql
  00008_create_audit_logs.sql
```

Recommended commands:

```bash
goose -dir migrations postgres "$DATABASE_URL" status
goose -dir migrations postgres "$DATABASE_URL" up
goose -dir migrations postgres "$DATABASE_URL" down
goose -dir migrations postgres "$DATABASE_URL" create create_courses sql
```

## Extensions And Conventions

Extensions:

- `pgcrypto` for `gen_random_uuid()` if IDs are generated in PostgreSQL.

Common columns:

- `id uuid primary key`
- `created_at timestamptz not null default now()`
- `updated_at timestamptz not null default now()`

Use enums or checked text columns for status/type fields. Use `jsonb` for polymorphic question content and answer payloads unless the project later chooses fully normalized per-question-type tables.

## Identity And Access

### organizations

Columns:

- `id uuid primary key`
- `name text not null`
- `inn text not null`
- timestamps

Constraints and indexes:

- `unique (inn)`
- `check (length(trim(name)) > 0)`
- `index organizations_name_idx on organizations using gin/to_tsvector or trigram`

### persons

Columns:

- `id uuid primary key`
- `first_name text not null`
- `last_name text not null`
- `middle_name text`
- `organization_id uuid null references organizations(id)`
- timestamps

Indexes:

- `index persons_last_name_idx`
- `index persons_organization_id_idx`

### person_profiles

Columns:

- `person_id uuid primary key references persons(id) on delete cascade`
- `snils text not null`
- `date_of_birth date not null`
- `job_title text`
- `education text`

Constraints and indexes:

- `unique (snils)`
- `check (date_of_birth <= current_date)`

### accounts

Columns:

- `id uuid primary key`
- `person_id uuid not null unique references persons(id)`
- `login text not null`
- `password_hash text not null`
- `role text not null`
- `status text not null`
- timestamps

Constraints and indexes:

- `unique (login)`
- `check (role in ('admin', 'creator', 'student', 'organization'))`
- `check (status in ('active', 'blocked', 'deleted'))`
- `index accounts_person_id_idx`

## Question Banks And Questions

### question_banks

Columns:

- `id uuid primary key`
- `title text not null`
- timestamps

Indexes:

- `index question_banks_title_idx`

### question_bank_questions

Columns:

- `bank_id uuid not null references question_banks(id) on delete cascade`
- `question_id uuid not null references questions(id) on delete restrict`
- `position int null`

Constraints and indexes:

- `primary key (bank_id, question_id)`
- `index question_bank_questions_question_id_idx`
- optional `unique (bank_id, position)` when ordering is needed

### questions

Columns:

- `id uuid primary key`
- `type text not null`
- `title text not null`
- `instruction text not null`
- `attachment jsonb null`
- `content jsonb not null`
- timestamps

Constraints and indexes:

- `check (type in ('selectable', 'matching', 'sequence', 'typed', 'short'))`
- `index questions_type_idx`
- `index questions_title_idx`

Question content examples:

- selectable: options with IDs and correctness flags
- matching: prompt/match pairs
- sequence: ordered options
- typed: text with blanks and accepted variants
- short: accepted answer variants

## Quizzes And Attempts

### quizzes

Columns:

- `id uuid primary key`
- `title text not null`
- `max_attempts int not null default 0`
- `time_limit_seconds int not null default 0`
- `shuffle_questions boolean not null default false`
- timestamps

Constraints:

- `check (max_attempts >= 0)`
- `check (time_limit_seconds >= 0)`

### quiz_sources

Columns:

- `id uuid primary key`
- `quiz_id uuid not null references quizzes(id) on delete cascade`
- `bank_id uuid not null references question_banks(id) on delete restrict`
- `criteria_type text not null`
- `question_count int not null`
- `question_ids uuid[] null`

Constraints and indexes:

- `unique (quiz_id, bank_id)`
- `check (criteria_type in ('random', 'manual'))`
- `check (question_count > 0)`
- `index quiz_sources_quiz_id_idx`
- `index quiz_sources_bank_id_idx`

### attempts

Columns:

- `id uuid primary key`
- `enrollment_id uuid not null references enrollments(id) on delete cascade`
- `quiz_id uuid not null references quizzes(id) on delete restrict`
- `status text not null`
- `started_at timestamptz not null`
- `deadline_at timestamptz null`
- `finished_at timestamptz null`
- timestamps

Constraints and indexes:

- `check (status in ('in_progress', 'finished', 'expired', 'interrupted', 'cancelled'))`
- `check (deadline_at is null or deadline_at >= started_at)`
- `check (finished_at is null or finished_at >= started_at)`
- `index attempts_enrollment_quiz_idx on attempts(enrollment_id, quiz_id)`
- `index attempts_status_idx`

### attempt_items

Columns:

- `attempt_id uuid not null references attempts(id) on delete cascade`
- `question_id uuid not null`
- `position int not null`
- `question_snapshot jsonb not null`

Constraints:

- `primary key (attempt_id, question_id)`
- `unique (attempt_id, position)`

### attempt_answers

Columns:

- `attempt_id uuid not null references attempts(id) on delete cascade`
- `question_id uuid not null`
- `answer jsonb not null`
- `answered_at timestamptz not null`

Constraints and indexes:

- `primary key (attempt_id, question_id)`
- `foreign key (attempt_id, question_id) references attempt_items(attempt_id, question_id) on delete cascade`

## Courses, Versions, Content

### courses

Columns:

- `id uuid primary key`
- `title text not null`
- timestamps

Indexes:

- `index courses_title_idx`

### course_versions

Columns:

- `id uuid primary key`
- `course_id uuid not null references courses(id) on delete cascade`
- `title text not null`
- `status text not null`
- `position int not null`
- timestamps

Constraints and indexes:

- `check (status in ('draft', 'published', 'deleted'))`
- `unique (course_id, position)`
- `index course_versions_course_status_idx on course_versions(course_id, status)`

### course_blocks

Columns:

- `id uuid primary key`
- `version_id uuid not null references course_versions(id) on delete cascade`
- `title text not null`
- `position int not null`
- timestamps

Constraints and indexes:

- `unique (version_id, position)`
- `index course_blocks_version_id_idx`

### course_elements

Columns:

- `id uuid primary key`
- `block_id uuid not null references course_blocks(id) on delete cascade`
- `title text not null`
- `content_type text not null`
- `content jsonb not null`
- `position int not null`
- timestamps

Constraints and indexes:

- `check (content_type in ('text', 'slides', 'document', 'video', 'quiz'))`
- `unique (block_id, position)`
- `index course_elements_block_id_idx`
- `index course_elements_content_type_idx`

For quiz elements, `content` should include `quiz_id`, and a database-level FK can be added with a generated column if strict SQL-level enforcement is required.

## Enrollment, Progress, Ratings

### enrollments

Columns:

- `id uuid primary key`
- `course_id uuid not null references courses(id) on delete restrict`
- `version_id uuid not null references course_versions(id) on delete restrict`
- `account_id uuid not null references accounts(id) on delete restrict`
- `activated_at date not null`
- `deactivated_at date not null`
- timestamps

Constraints and indexes:

- `unique (account_id, course_id, version_id)`
- `check (deactivated_at >= activated_at)`
- `index enrollments_account_id_idx`
- `index enrollments_course_version_idx on enrollments(course_id, version_id)`
- application rule: enrollment is rejected when the target version status is `deleted`

### course_progress

Columns:

- `id uuid primary key`
- `enrollment_id uuid not null unique references enrollments(id) on delete cascade`
- `version_id uuid not null references course_versions(id) on delete restrict`
- timestamps

Indexes:

- `index course_progress_version_id_idx`

### course_progress_markers

Columns:

- `progress_id uuid not null references course_progress(id) on delete cascade`
- `element_id uuid not null references course_elements(id) on delete cascade`
- `marker_type text not null`
- `completed_at timestamptz not null`

Constraints:

- `primary key (progress_id, element_id)`
- `check (marker_type in ('read', 'watched', 'download', 'quiz'))`

### student_ratings view

Suggested view columns:

- `account_id`
- `enrollment_id`
- `course_id`
- `version_id`
- `completed_items`
- `total_items`
- `completion_percent`

The view should count tracked elements, including quizzes, PDFs/readable documents, and videos. Static downloadable files may be excluded from required progress unless the business decides download should count.

Indexes that support the view:

- `course_elements(block_id, content_type)`
- `course_progress_markers(progress_id, element_id)`
- `enrollments(account_id, course_id, version_id)`

## Object Storage Metadata

### stored_objects

Columns:

- `id uuid primary key`
- `storage_key text not null`
- `file_name text not null`
- `content_type text not null`
- `size_bytes bigint not null`
- `bucket text not null`
- timestamps

Constraints and indexes:

- `unique (bucket, storage_key)`
- `check (size_bytes > 0)`
- `index stored_objects_content_type_idx`

Course element `content` JSON can reference `stored_object_id` or `storage_key`.

## Audit

### audit_logs

Columns:

- `id uuid primary key`
- `actor_account_id uuid null references accounts(id)`
- `actor_role text not null`
- `action text not null`
- `entity_type text not null`
- `entity_id uuid null`
- `payload jsonb not null default '{}'::jsonb`
- `created_at timestamptz not null default now()`

Indexes:

- `index audit_logs_actor_account_id_idx`
- `index audit_logs_entity_idx on audit_logs(entity_type, entity_id)`
- `index audit_logs_action_idx`
- `index audit_logs_created_at_idx`
