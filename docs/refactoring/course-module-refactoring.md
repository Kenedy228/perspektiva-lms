# Course Module Refactoring (Final)

## Scope
- Removed legacy `course/version` flow from course module boundaries.
- Course application flow is now `Course -> Block -> Element`.
- Updated persistence mapping to support:
  - ordered `course -> blocks` links;
  - ordered `block -> elements` links;
  - `completion_mode` for course elements;
  - progress persistence without `course_progress.version_id` dependency.

## Domain model
- `Course`:
  - owns ordered `blockIDs`;
  - supports add/remove/move block operations.
- `Block`:
  - owns ordered `elementIDs`;
  - supports add/remove/move element operations.
- `Element`:
  - supports content types:
    - `test`
    - `lecture_material`
    - `download_file`
  - supports completion modes:
    - `none`
    - `manual`

## Application layer changes
- Removed version-centered course use cases.
- Added:
  - `AddBlockToCourseUseCase`
  - `MoveCourseBlockUseCase`
  - `AddElementToBlockUseCase`
  - `MoveBlockElementUseCase`

## HTTP transport changes
- Removed version-centered endpoints from course flow.
- New endpoints:
  - `POST /courses/{courseID}/blocks`
  - `POST /blocks/{blockID}/elements`
  - `POST /courses/{courseID}/progress`

## Persistence changes
- Added migration `00003_course_model_schema_update.sql`:
  - adds `course_blocks_links (course_id, block_id, position)`;
  - adds `course_elements.completion_mode`;
  - drops `course_progress.version_id`;
  - includes backfill and rollback logic.
- Postgres mapping updated accordingly.

## Notes
- Enrollment and attempt modules still use `enrollments.version_id` and their own version-oriented constraints.
- Full removal of version-dependent artifacts outside course module should be done in dedicated follow-up steps.
