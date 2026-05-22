# PostgreSQL Implementation

The initial database schema is implemented in `migrations/00001_init.sql` and is managed with `goose`.

## Local Services

`docker-compose.yaml` starts:

- PostgreSQL 18 on host port `5433`
- MinIO latest on host ports `9000` and `9001`

PostgreSQL 18 stores data under `/var/lib/postgresql`, so the compose file uses a `postgres18_data` volume mounted at that path.

## Migration Command

```bash
/home/kenedy/go/bin/goose -dir migrations postgres 'postgres://lms:lms@localhost:5433/lms?sslmode=disable' up
```

The migration creates identity, account, question bank, question, quiz, course, enrollment, progress, attempt, storage metadata, and audit tables.

## Repository Adapters

The `backend/internal/infrastructure/postgres` package contains PostgreSQL adapters for the current application ports:

- account repository, query service, and audit recorder
- organization repository, query service, and audit recorder
- person repository, query service, and audit recorder
- question bank repository, query service, audit recorder, and quiz bank inspector
- quiz repository
- enrollment repository and progress initializer
- course, version, block, progress repositories, course query service, and course access policy
- attempt repository and attempt access policy
- question repository and attempt question provider

Question restoration, attempt restoration, answer restoration, and course element content persistence use explicit infrastructure DTOs. These DTOs keep PostgreSQL JSON payloads out of the domain model while preserving domain constructors and invariants during restore.

## Integration Tests

Repository/use-case integration coverage is under the `integration` build tag:

```bash
go test -tags=integration ./backend/internal/infrastructure/postgres
```

The test starts PostgreSQL 18 with testcontainers, applies `00001_init.sql`, and executes an organization use-case lifecycle through the PostgreSQL adapter.
