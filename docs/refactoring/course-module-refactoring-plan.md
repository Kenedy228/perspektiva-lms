# Course Module Refactoring Plan (Analysis Only)

## Scope and constraints
- This document is based strictly on the current `perspektiva-lms` codebase.
- Runtime code is intentionally unchanged for this step.
- Focus areas requested: `course/version`, `versionID`, migrations, repositories, HTTP routes.

## Current course-module architecture map

### Domain layer
- `backend/internal/domain/course/course.go`
  - Aggregate `Course` keeps `versionIDs []uuid.UUID`.
  - Invariants: required ID/title, versions max limit, unique/non-nil `versionID`.
- `backend/internal/domain/course/version/version.go`
  - Aggregate `Version` keeps `status` (`draft/published/deleted`) and `blockIDs`.
  - Editing (title/blocks/move) allowed only in `draft`.
  - `Publish()` requires at least one block.
- `backend/internal/domain/course/block/block.go`
  - Aggregate `Block` stores ordered `elementIDs`.
- `backend/internal/domain/course/element/*`
  - `Element` + typed content (text/slides/document/video/quiz).
- `backend/internal/domain/course/progress/progress.go`
  - `Progress` is enrollment-scoped and stores `versionID` + markers.

### Application ports
- `backend/internal/application/ports/course/repository.go`
  - `CourseRepository`, `VersionRepository`, `BlockRepository`, `ElementRepository`, `ProgressRepository`, `EnrollmentAccess`.
- `backend/internal/application/ports/course/query.go`, `views.go`
  - Read-model contracts for list/details/ratings/statistics.

### Application use cases
- Commands (`backend/internal/application/usecases/course/commands`):
  - `CreateCourseUseCase`
  - `RenameCourseUseCase`
  - `CreateVersionUseCase` (creates `Version`, appends `versionID` to `Course`)
  - `AddBlockUseCase` (adds block to `Version` by `VersionID`)
  - `PublishVersionUseCase` (publishes by `VersionID`)
  - `MarkProgressUseCase` (progress by enrollment; progress internally bound to `versionID`)
- Queries (`backend/internal/application/usecases/course/queries`):
  - `ListQuery`, `RatingsQuery`, `StudentStatisticsQuery`.

### Infrastructure (Postgres)
- `backend/internal/infrastructure/postgres/course.go`
  - `CourseRepository` maps `courses` + `course_version_links`.
  - `VersionRepository` maps `course_versions` + `course_version_blocks`.
  - `BlockRepository`, `ElementRepository`, `ProgressRepository`, `CoursePolicy`, `CourseQueryService`.
- `backend/internal/infrastructure/postgres/enrollment.go`
  - Enrollment persistence uses `version_id`.
- `backend/internal/infrastructure/postgres/attempt.go`
  - Attempt dependency traversal joins via enrollment `version_id` and course structure link tables.

### Transport (HTTP)
- Route registration: `backend/internal/transport/http/server.go`
  - `POST /courses/{id}/versions`
  - `POST /course-versions/{id}/blocks`
  - `POST /course-versions/{id}/publish`
  - `POST /course-progress/{enrollmentID}/elements/{elementID}`
  - `GET /courses`, `POST /courses`, `GET /courses/{id}`, `PATCH /courses/{id}`, `GET /courses/{id}/ratings`
- Handlers: `backend/internal/transport/http/handlers/course_api.go`
  - `CreateCourseVersion` passes path course `{id}` + title into `CreateVersionUseCase`.
  - `AddCourseBlock` and `PublishCourseVersion` are version-centric endpoints.
- DTOs: `backend/internal/transport/http/handlers/dto.go`
  - `CourseVersionRequest`, `CourseBlockRequest`, plus `EnrollmentRequest.VersionID`.

## `versionID` dependency map (cross-module)

### Direct course-module usage
- `Course.versionIDs` domain field in aggregate.
- Version link persistence table: `course_version_links(course_id, version_id, position)`.
- Version-block table: `course_version_blocks(version_id, block_id, position)`.
- Course progress table: `course_progress.version_id`.

### External modules coupled to `versionID`
- Enrollment domain and persistence:
  - `enrollments.version_id` is required and part of unique constraint `(account_id, course_id, version_id)`.
  - Enrollment use case checks `EnrollmentAccess.CanEnrollVersion(versionID)`.
- Attempt infrastructure:
  - Query path traverses enrollment `version_id` -> `course_version_blocks` -> `course_block_elements` -> `course_elements`.
- Course statistics query:
  - Returns `version_id` in rating/stat rows.

## Database and migration state

### Current schema objects tied to version model
- `course_versions`
- `course_version_links`
- `course_version_blocks`
- `enrollments.version_id` FK to `course_versions`
- `course_progress.version_id` FK to `course_versions`
- Related indexes:
  - `course_versions_status_idx`
  - `course_version_links_version_id_idx`
  - `course_version_blocks_block_id_idx`
  - `enrollments_version_id_idx`
  - `course_progress_version_id_idx`

### Migration files observed
- `migrations/00001_init.sql` defines all version-dependent entities above.
- `migrations/00002_add_polymorphic_mapper_columns.sql` is unrelated to course version topology.

## Tests and regression surface

### Existing tests that already encode current behavior
- Domain:
  - `backend/internal/domain/course/version/version_test.go`
- Course use cases:
  - `backend/internal/application/usecases/course/commands/commands_test.go`
  - `backend/internal/application/usecases/course/queries/queries_test.go`
- Enrollment behavior coupled to version access:
  - `backend/internal/application/usecases/enrollment/commands/create_test.go`
- Integration setup includes version-linked enrollment:
  - `backend/internal/infrastructure/postgres/postgres_integration_test.go` (`insertAttemptDependencies`)

### High-risk regression zones for future refactor
- Enrollment creation path and uniqueness logic.
- Attempt question resolution SQL joins.
- Student progress compatibility (`course_progress.version_id`).
- HTTP contract compatibility for version-centric endpoints.

## Refactoring sequence (planned for next tasks)

1. Freeze behavioral baseline with additional targeted tests (before refactor)
- Add/extend tests for:
  - course-version creation/link order preservation;
  - block add/publish draft invariants;
  - enrollment/version policy checks;
  - attempt dependency SQL path assumptions.

2. Isolate version semantics in application boundaries
- Introduce transitional adapters/DTO mapping inside application/infrastructure only.
- Keep HTTP/domain contracts stable while internals move.

3. Repository-level transition
- Update postgres repositories first, behind existing ports where possible.
- Avoid leaking SQL or persistence details into use cases.

4. Migration phase
- Add explicit forward migrations for any schema change.
- Preserve backward compatibility windows where cross-module readers still expect `version_id`.

5. Transport and contract cleanup
- Update HTTP routes/DTOs only after app+repo compatibility is in place.
- Keep error shape and auth behavior unchanged.

6. Final contract hardening
- Remove deprecated fields/routes only after all dependent modules are migrated and tested.

## Concrete coupling checklist (for implementation phases)
- Domain aggregates:
  - `course.Course.versionIDs`
  - `version.Version`
- Use cases:
  - `CreateVersionUseCase`, `AddBlockUseCase`, `PublishVersionUseCase`
- Repositories:
  - `CourseRepository`, `VersionRepository`, `ProgressRepository`, enrollment repository, attempt SQL readers
- Transport:
  - `/courses/{id}/versions`, `/course-versions/{id}/blocks`, `/course-versions/{id}/publish`
  - `EnrollmentRequest.VersionID`
- Migrations:
  - all schema objects listed in “Database and migration state”.

## Assumptions and boundaries for next steps
- No changes to question/grading/attempt domain behavior unless strictly required by version decoupling.
- No large rewrite: move in small, test-protected, reversible steps.
- Preserve clean architecture boundaries: domain <- application <- infrastructure/transport.
