## Project overview
- This project is an internal Learning Management System (LMS) for the organization, conceptually similar to Moodle.
- The system is split into two main parts: `backend` and `frontend`.
- Backend is implemented in **Golang**.
- Frontend is implemented in **React**.
- The goal is to build a maintainable, secure, extensible, and production-ready platform aligned with organizational requirements.

## Repository structure
- `backend/` — Go services, domain logic, application layer, infrastructure, integrations, background jobs, migrations, tests.
- `frontend/` — React application, UI components, pages, state management, API integration, tests.
- `docs/` — architecture notes, ADRs, API contracts, business rules, diagrams.
- `deploy/` or `infra/` — Docker, CI/CD, environment configs, deployment manifests.

## General engineering principles
- Preserve clean architecture and clear separation of concerns.
- Prefer explicit, simple, readable code over clever abstractions.
- Keep modules small, cohesive, and testable.
- Avoid hidden coupling between backend modules and frontend features.
- Every non-trivial decision should optimize for maintainability and long-term evolution of the LMS.
- Before adding a dependency, justify why the standard library or current stack is insufficient.
- Do not introduce large-scale refactors unless they are required for the task.
- Do not break existing public contracts without documenting the change.

## Domain expectations
- Treat the LMS as a business-critical system.
- Core domain concepts may include: users, roles, courses, modules, lessons, enrollments, groups, assignments, submissions, quizzes, attempts, grades, certificates, files, notifications, schedules, and audit logs.
- Prefer explicit domain terminology in code and documentation.
- Business rules must live in domain/application logic, not be scattered through handlers, controllers, or UI code.
- Protect invariants at the domain level.

## Architecture rules
- Follow a layered or clean architecture approach:
    - `domain` — entities, value objects, domain services, business invariants, repository interfaces.
    - `application` — use cases, orchestration, DTOs, commands, queries, transaction boundaries.
    - `infrastructure` — database, storage, message brokers, HTTP clients, email, cache, external integrations.
    - `interfaces` / `transport` — HTTP handlers, gRPC handlers, middleware, presenters.
- Dependencies must point inward: infrastructure depends on application/domain, not vice versa.
- Avoid leaking persistence models into domain and transport layers.
- Use interfaces where they express a boundary, not as ceremony.

## Backend rules (Golang)
- Target idiomatic Go.
- Prefer composition over inheritance-style patterns.
- Keep interfaces close to the consumer.
- Return explicit errors; never swallow them.
- Wrap errors with context using `fmt.Errorf("...: %w", err)`.
- Handle context propagation correctly for I/O, DB, HTTP, and long-running operations.
- Do not use global mutable state unless absolutely necessary and documented.
- Keep functions focused and short where practical.
- Prefer constructor-based dependency injection.
- Use table-driven tests where they improve clarity.

## Go project conventions
- Recommended package organization inside `backend/`:
    - `internal/domain/...`
    - `internal/application/...`
    - `internal/infrastructure/...`
    - `internal/transport/http/...`
    - `internal/platform/...`
- Use `cmd/` for service entry points.
- Keep reusable public packages minimal; prefer `internal/` by default.
- Package names should be short, lowercase, and meaningful.
- Avoid generic utility packages like `common`, `helpers`, or `misc` unless the content is truly cross-cutting and coherent.

## API and integration rules
- APIs must be consistent in naming, error shape, pagination, filtering, and authorization behavior.
- Validate input at the boundary layer.
- Keep transport DTOs separate from domain entities.
- Document contract changes in `docs/`.
- For breaking API changes, note migration impact clearly.

## Data and persistence
- Database schema changes must be done through explicit migrations.
- Never modify production-critical data logic without considering backward compatibility and rollback.
- Repositories should encapsulate persistence details.
- Keep transactions scoped to the application use case boundary.
- Be careful with N+1 queries, large scans, locking, and missing indexes.
- Prefer explicit queries over magical ORM behavior.

## Security and compliance
- Treat authentication, authorization, and personal data carefully.
- Enforce least-privilege access.
- Never expose sensitive information in logs, errors, or API responses.
- Validate and sanitize file upload workflows.
- Audit important actions such as enrollment changes, grading changes, role changes, and administrative actions.
- Assume the platform may store personal and academic data; design accordingly.

## Frontend rules (React)
- Keep UI components small, predictable, and reusable.
- Separate presentational concerns from business/data-fetching concerns where practical.
- Co-locate feature code when it improves maintainability.
- Avoid tightly coupling UI state to backend response shapes.
- Handle loading, empty, error, and permission-denied states explicitly.
- Prefer accessible components and semantic HTML.

## Testing expectations
- Add or update tests for every non-trivial change.
- Backend tests should include unit tests for business rules and integration tests for repositories/HTTP flows where needed.
- Frontend tests should cover critical UI behavior and user flows.
- Bug fixes should include a regression test when practical.
- Do not remove tests just to make the build pass.

## Observability and operations
- Add structured logging for important flows and failures.
- Logs should be useful for debugging but must not leak secrets or personal data.
- Prefer metrics and tracing hooks for critical business operations if observability is already present in the project.
- Make failure modes visible and actionable.

## Documentation
- Update documentation when changing architecture, domain rules, setup flow, environment variables, or API contracts.
- Write ADRs for major architectural decisions.
- Keep README/setup instructions accurate.
- Document assumptions when requirements are ambiguous.

## Change workflow for Codex
- First read nearby code and infer the local architectural pattern before editing.
- Match the existing style unless it is clearly harmful.
- When implementing a feature:
    1. Identify the domain concept and use case.
    2. Define or update contracts.
    3. Implement backend logic.
    4. Add/update tests.
    5. Update frontend integration if needed.
    6. Update documentation.
- When fixing a bug:
    1. Reproduce or infer the failure path.
    2. Find the root cause.
    3. Apply the smallest safe fix.
    4. Add regression coverage.

## Code style expectations
- Use clear names based on domain language.
- Avoid comments that restate the code; prefer comments for intent, invariants, or non-obvious decisions.
- Keep files reasonably sized.
- Extract functions when it improves readability, not just to reduce line count.
- Prefer explicitness over hidden framework magic.

## Things to avoid
- Do not mix transport, business logic, and persistence in one place.
- Do not create massive service objects with unrelated responsibilities.
- Do not introduce circular dependencies.
- Do not hardcode organization-specific constants across the codebase; centralize configuration.
- Do not bypass authorization checks for “internal-only” flows.
- Do not implement silent fallbacks that hide failures.

## Preferred output from Codex
- Produce changes that are minimal, coherent, and ready for review.
- Explain trade-offs when a task has multiple valid implementation paths.
- Highlight unclear requirements before making irreversible architectural choices.
- When requirements are missing, choose the safest maintainable default and state the assumption.

## If requirements are unclear
Ask for clarification about:
- multi-tenant vs single-tenant architecture;
- authentication method (local auth, SSO, LDAP, OAuth2, etc.);
- role model (student, teacher, curator, admin, reviewer, etc.);
- grading and certification rules;
- file storage and video delivery;
- notification channels;
- reporting and audit requirements;
- integration with existing organizational systems.