# AGENTS.md

## Scope
- This file applies to the `backend/` directory only.
- Backend stack:
    - **Golang**
    - **Fiber** for HTTP transport
    - **PostgreSQL** via `pgx` using the standard `database/sql` package
    - **MinIO** for object storage
- This backend is part of an internal LMS platform similar to Moodle and must support long-term maintainability, security, and evolving business rules.

## Backend goals
- Build a production-grade backend for an LMS with clear domain boundaries.
- Optimize for correctness, maintainability, observability, and testability.
- Keep HTTP, business logic, persistence, and storage concerns separated.
- Prefer explicit architecture over framework-driven magic.

## Architectural style
- Follow **Clean Architecture** / **DDD-lite** principles.
- Organize the backend into these layers:
    - `domain` — entities, value objects, repository contracts, domain services, invariants.
    - `application` — use cases, commands, queries, DTOs, transaction orchestration.
    - `infrastructure` — PostgreSQL repositories, MinIO storage adapter, external integrations, config, logging.
    - `transport/http` — Fiber app setup, handlers, routing, middleware, request/response DTOs.
- Dependencies must point inward.
- Domain code must not depend on Fiber, SQL details, MinIO SDK, or config implementations.
- Application services may depend on interfaces, never on framework-specific types.

## Recommended folder structure
```text
backend/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── domain/
│   │   ├── user/
│   │   ├── course/
│   │   ├── enrollment/
│   │   ├── assignment/
│   │   ├── submission/
│   │   ├── quiz/
│   │   ├── grade/
│   │   ├── file/
│   │   └── shared/
│   ├── application/
│   │   ├── user/
│   │   ├── course/
│   │   ├── enrollment/
│   │   ├── assignment/
│   │   ├── submission/
│   │   ├── quiz/
│   │   ├── grade/
│   │   └── common/
│   ├── infrastructure/
│   │   ├── config/
│   │   ├── db/
│   │   ├── persistence/
│   │   ├── storage/
│   │   ├── auth/
│   │   ├── logger/
│   │   ├── clock/
│   │   └── id/
│   ├── transport/
│   │   └── http/
│   │       ├── handler/
│   │       ├── middleware/
│   │       ├── request/
│   │       ├── response/
│   │       ├── presenter/
│   │       └── router/
│   └── platform/
│       └── app/
├── migrations/
├── test/
├── docs/
├── go.mod
└── go.sum
```

## Domain modeling rules
- Use domain terminology from the LMS, not generic CRUD naming.
- Likely aggregates/modules include:
    - users
    - roles
    - courses
    - course modules
    - lessons
    - enrollments
    - assignments
    - submissions
    - quizzes
    - attempts
    - grades
    - certificates
    - files
    - notifications
    - audit logs
- Protect invariants inside domain logic.
- Do not place business rules in Fiber handlers or SQL repository code.
- Avoid anemic domain models where rules are scattered in services and handlers.

## Naming rules
- Name packages by business meaning, not technical mechanism.
- Good: `course`, `submission`, `grade`, `auth`, `storage`.
- Bad: `utils`, `helpers`, `commonstuff`, `serviceimpl`.
- Use explicit names for use cases, for example:
    - `CreateCourse`
    - `EnrollUserInCourse`
    - `SubmitAssignment`
    - `PublishGrade`
- Repository interfaces should express domain intent, not SQL vocabulary.

## Fiber rules
- Fiber is only the HTTP delivery mechanism.
- Handlers must stay thin:
    - parse request;
    - validate boundary input;
    - call application use case;
    - map result to HTTP response.
- Do not pass `*fiber.Ctx` outside transport layer.
- Do not put business logic into middleware.
- Middleware should be limited to cross-cutting concerns: auth, request ID, logging, panic recovery, CORS, rate limiting, tracing, metrics.
- Standardize error responses and success envelopes where appropriate.
- Keep routing declarative and grouped by bounded context or feature.

## HTTP API rules
- APIs must be consistent in:
    - naming;
    - resource structure;
    - pagination;
    - filtering;
    - sorting;
    - error payloads;
    - auth and permission behavior.
- Validate all external input at the boundary.
- Distinguish between:
    - malformed input;
    - unauthorized;
    - forbidden;
    - not found;
    - conflict;
    - validation failure;
    - internal error.
- Never leak raw DB or MinIO errors directly to clients.
- Keep request/response DTOs separate from domain entities.

## Go coding conventions
- Write idiomatic Go.
- Prefer composition over inheritance-like abstractions.
- Keep interfaces near the consumer.
- Do not create interfaces prematurely.
- Return errors explicitly and wrap them with context using `fmt.Errorf("...: %w", err)`.
- Use `context.Context` for DB, storage, network, and long-running operations.
- Avoid global mutable state.
- Prefer constructor injection.
- Keep functions focused and coherent.
- Prefer small packages with strong cohesion.

## Error handling rules
- All errors must carry useful context.
- Sentinel errors may be used for stable business cases such as:
    - not found;
    - duplicate enrollment;
    - forbidden action;
    - invalid attempt state.
- Map domain/application errors to HTTP status codes centrally.
- Log internal errors with context, but do not expose internals in API responses.
- Do not ignore returned errors, including `rows.Err()`, `tx.Commit()`, `tx.Rollback()`, and storage cleanup errors.

## PostgreSQL rules
- Use PostgreSQL through `pgx` with `database/sql` compatibility via `github.com/jackc/pgx/v5/stdlib`.
- Open DB connections through one centralized initialization path.
- Configure connection pool limits explicitly.
- Check connectivity during startup with `PingContext`.
- Repositories must depend on `*sql.DB` or transaction abstractions, not on ad-hoc globals.
- Keep SQL explicit; avoid hidden ORM behavior.
- Prefer handwritten SQL with clear column lists.
- Always use placeholders and parameterized queries.
- Use `QueryContext`, `QueryRowContext`, and `ExecContext`.
- Close rows properly and check iteration errors.
- Treat transactions as application-level boundaries.

## SQL style rules
- Store SQL close to the repository/use case it belongs to.
- Use explicit column names, never `SELECT *` in production code.
- Make scan order obvious and stable.
- Be careful with nullable fields; use proper Go nullable handling.
- Optimize for query clarity first, then performance.
- Consider indexes whenever adding new filtering, sorting, join-heavy, or reporting paths.
- Watch for N+1 problems in course trees, submission lists, gradebooks, and reporting endpoints.

## Transaction rules
- Transactions belong in the application layer when a use case spans multiple repositories.
- Repository methods should support execution under either `*sql.DB` or `*sql.Tx` through a narrow abstraction when needed.
- Keep transactions short.
- Never perform external network calls to MinIO or third-party systems while holding a DB transaction unless explicitly required and carefully justified.
- If a workflow spans DB + object storage, design compensating behavior for partial failure.

## Migration rules
- All schema changes must go through versioned migrations.
- Never edit already-applied migrations in shared environments.
- Each migration should be reversible when practical.
- Migrations must be reviewed for:
    - backward compatibility;
    - lock impact;
    - index creation cost;
    - data backfill safety.
- Document destructive or high-risk migrations.

## MinIO rules
- Use MinIO through a dedicated storage adapter in `internal/infrastructure/storage`.
- Do not spread MinIO SDK calls across handlers or use cases.
- All object storage operations must accept `context.Context`.
- Centralize bucket names, path conventions, and content-type handling.
- Define clear object key conventions, for example by entity type and ownership scope.
- Validate file size, media type, and allowed upload category before storage.
- Prefer streaming where appropriate.
- For larger uploads, rely on multipart-capable upload behavior supported by the SDK.
- Handle partial-failure scenarios explicitly when DB metadata and object upload can diverge.
- Presigned URLs must have short expiration and proper access control.

## File handling rules
- Treat uploaded files as security-sensitive input.
- Never trust client-provided MIME type alone.
- Sanitize filenames if filenames are stored or displayed.
- Store metadata in PostgreSQL and binary objects in MinIO.
- Define lifecycle rules for temporary files, draft uploads, and orphan cleanup.
- Access to course materials, submissions, and private attachments must be authorization-checked at application level.

## Security rules
- Treat this LMS as a system containing personal, academic, and administrative data.
- Enforce authentication and authorization everywhere.
- Prefer explicit permission checks over implicit assumptions.
- Sensitive operations must be audited.
- Never log secrets, tokens, passwords, file presigned URLs, or personal data unnecessarily.
- Hash passwords with a modern password hashing algorithm if local auth is used.
- Validate and constrain upload endpoints carefully.
- Protect admin-only and grading operations with role/permission checks.

## Config rules
- All configuration must come from explicit config structures and environment variables.
- No hardcoded secrets, endpoints, bucket names, or DSNs.
- Fail fast on invalid configuration.
- Configuration should include at least:
    - HTTP server settings;
    - database DSN and pool settings;
    - MinIO endpoint, credentials, region, bucket names, secure mode;
    - auth settings;
    - logging level;
    - CORS / trusted origins;
    - feature flags if used.

## Logging and observability
- Use structured logs.
- Include request ID / correlation ID in logs.
- Log important business events and failures at appropriate levels.
- Keep logs useful for debugging production issues.
- Do not log sensitive payloads.
- Prefer metrics/tracing hooks for critical workflows if observability stack exists.
- Important flows to observe include:
    - login;
    - enrollment;
    - assignment submission;
    - grading;
    - file upload/download;
    - certificate generation.

## Testing rules
- Every non-trivial change must include tests.
- Prefer **testify** as the default testing framework for assertions, suites where appropriate, and mocks, so tests stay compact and readable.
- Prefer `require` for hard preconditions and `assert` for multiple checks in the same test.
- Prefer `testify/mock` for mocking external boundaries when a lightweight fake is not clearer.
- Keep test code expressive and short; avoid overly verbose manual assertion boilerplate.
- Follow TDD where practical:
    1. write a failing test for the behavior or invariant;
    2. implement the smallest change that makes it pass;
    3. refactor while keeping tests green.
- Prioritize tests for:
    - domain invariants;
    - application use cases;
    - repository behavior;
    - HTTP handlers for critical paths;
    - storage adapter behavior where feasible.
- Domain tests must verify invariants as fully as practical, including valid flows, invalid transitions, edge cases, duplicate actions, forbidden states, boundary values, and regression scenarios.
- For aggregates and entities, prefer exhaustive invariant coverage over superficial happy-path checks.
- Every business rule added or changed should have at least one test proving the intended behavior and one or more tests proving invalid states are rejected.
- Use table-driven tests when they improve coverage and readability, especially for validation rules, permissions, state transitions, and edge cases.
- Bug fixes should include regression coverage.
- Avoid brittle tests tied to irrelevant implementation details.
- Prefer integration tests for repository logic against real PostgreSQL when feasible.
- Mock only true external boundaries such as object storage, external APIs, message brokers, clock/id generators when needed, or other infrastructure dependencies.
- Do not over-mock domain logic; prefer real domain objects in domain and application tests.
- Handler tests should verify HTTP contract: status code, response body, validation behavior, auth behavior, and error mapping.
- Repository integration tests should verify scanning, constraints, transactions, filtering, and pagination behavior.

## Performance rules
- Be careful with LMS hotspots:
    - course tree loading;
    - gradebook queries;
    - quiz attempts;
    - submission listing;
    - file delivery metadata;
    - reporting endpoints.
- Avoid loading large result sets unnecessarily.
- Paginate collection endpoints by default.
- Use indexes intentionally.
- Keep object storage operations out of hot synchronous paths when possible.
- Measure before optimizing, but do not ignore obvious inefficiencies.

## Documentation rules
- Update backend docs when changing:
    - architecture;
    - environment variables;
    - migrations;
    - API contracts;
    - file storage behavior;
    - auth model.
- Add ADRs for major architectural decisions.
- Keep setup instructions reproducible.
- Document assumptions if requirements are unclear.

## Preferred implementation style
- Make the smallest coherent change that solves the problem.
- Match existing conventions unless they are clearly harmful.
- When adding a new feature:
    1. identify the domain concept;
    2. define or update use case contract;
    3. implement domain/application logic;
    4. implement repository/storage adapters;
    5. expose via Fiber handler;
    6. add tests;
    7. update docs.
- When fixing a bug:
    1. identify the failure path;
    2. confirm root cause;
    3. apply minimal safe fix;
    4. add regression test.

## Things to avoid
- Do not put SQL in Fiber handlers.
- Do not put MinIO calls directly in handlers.
- Do not pass framework DTOs into domain logic.
- Do not create god-services with mixed responsibilities.
- Do not mix auth, validation, persistence, and business rules in one function.
- Do not introduce unnecessary abstractions for one implementation unless they define a real boundary.
- Do not use generic repository patterns that erase domain meaning.
- Do not use `SELECT *`.
- Do not hold DB transactions open during slow file/network operations without strong reason.

## Expected output from Codex in backend/
- Produce reviewable, minimal, production-oriented changes.
- Preserve architecture boundaries.
- Explain trade-offs when multiple valid solutions exist.
- Call out missing requirements before making irreversible decisions.
- If requirements are ambiguous, choose the safest maintainable default and state the assumption.

## Clarifications worth asking when needed
- auth method: local auth, SSO, LDAP, OAuth2;
- role matrix: student, teacher, curator, admin, reviewer, manager;
- course publication workflow;
- grading policy and attempt rules;
- submission retention and file size limits;
- audit/compliance requirements;
- reporting requirements;
- notification and eventing strategy;
- multi-tenant vs single-tenant scope.
