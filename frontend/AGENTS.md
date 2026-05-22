# AGENTS.md (frontend)

## Scope
- This file applies to the `frontend/` directory only.
- Frontend stack:
  - **React** (with hooks and functional components only).
  - **TypeScript** (strict mode for all new code).
  - Styling: CSS Modules or a minimal utility-first approach (if a styling choice already exists in the repo, follow it).
  - HTTP: a dedicated API client layer for talking to the Go backend (Fiber + PostgreSQL + MinIO).
- The backend (Go + Fiber + PostgreSQL + MinIO) is already implemented and is the single source of truth for business logic and data.
- The frontend is a web UI for the internal LMS similar to Moodle.

## High-level goals
- Build a **maintainable**, **type-safe**, and **predictable** React frontend for the LMS.
- Respect backend contracts: never "invent" fields or behaviors that do not exist in the backend API.
- Keep components simple and focused; push data fetching and business rules into dedicated layers.
- Optimize for clarity over cleverness.

## Project structure
Use a feature-first structure and keep cross-cutting concerns separate.

```text
frontend/
├── src/
│   ├── app/
│   │   ├── App.tsx
│   │   ├── router.tsx
│   │   └── providers/
│   ├── features/
│   │   ├── auth/
│   │   ├── courses/
│   │   ├── organizations/
│   │   ├── people/
│   │   ├── accounts/
│   │   ├── files/
│   │   └── ...
│   ├── components/
│   │   ├── ui/
│   │   └── layout/
│   ├── api/
│   │   ├── client.ts
│   │   ├── auth.ts
│   │   ├── courses.ts
│   │   ├── organizations.ts
│   │   └── ...
│   ├── lib/
│   │   ├── hooks/
│   │   ├── forms/
│   │   ├── table/
│   │   └── utils/
│   ├── types/
│   │   ├── auth.ts
│   │   ├── courses.ts
│   │   └── ...
│   ├── styles/
│   └── test/
├── package.json
├── tsconfig.json
└── vite.config.ts (or other bundler config)
```

Rules:
- `src/features/*` is the primary place for feature-specific UI and state.
- `src/components/ui` contains reusable UI primitives (buttons, inputs, modals, etc.).
- `src/api` contains API clients and typed request/response DTOs.
- `src/lib` contains shared hooks, form helpers, and utilities.
- `src/types` contains shared TypeScript types that match backend contracts.

## Code style
- Use **TypeScript** with `strict` mode enabled for all new files.
- Only **functional React components** with hooks.
- Use React 18+ idioms (no legacy lifecycle methods).
- Prefer **named exports** over default exports.
- Components use **PascalCase** filenames: `CourseList.tsx`, `LoginForm.tsx`.
- Non-component modules use **camelCase** filenames: `useAuth.ts`, `formatDate.ts`.
- Do not use `any`. If unavoidable, isolate it and annotate with a TODO and reasoning.
- Prefer explicit types over implicit `any` from dynamic JSON.
- Use existing ESLint and Prettier configuration; do not override it without necessity.

## Component patterns

### Dumb vs. smart components
- **UI components** (`src/components/ui`) must be dumb/presentational:
  - Receive data and callbacks via props.
  - No direct data fetching.
  - No knowledge of backend endpoints.
- **Feature components** (`src/features/*`) may:
  - Use hooks to fetch data from the API client.
  - Compose UI components.
  - Contain feature-level state and view logic.

### Hooks
- Use dedicated hooks for data fetching: e.g. `useCourses`, `useCourseDetails`, `useAuth`.
- Custom hooks live in `src/lib/hooks` or inside the corresponding feature folder.
- Hooks should be composable and small.
- Avoid putting complex business logic directly inside components; move it into hooks or pure functions.

### Forms
- Use a single form library across the app (e.g., React Hook Form, if already present).
- Encapsulate form schemas and validation rules close to the feature.
- Handle validation both on the client (for UX) and rely on backend validation as the source of truth.

## API access
- **IMPORTANT**: All HTTP requests must go through `src/api`.
- Do **NOT** call `fetch` or axios directly from components.
- `src/api/client.ts` should:
  - Configure base URL for the Go backend.
  - Handle auth headers/tokens.
  - Provide a small wrapper around `fetch` (or chosen HTTP client) with proper error handling.
- Each domain area gets its own API module:
  - `auth.ts` for login/logout/current user.
  - `courses.ts` for course listing, details, enrollment.
  - `organizations.ts` for organization management.
  - `people.ts`, `accounts.ts`, `files.ts`, etc.
- API modules should:
  - Export typed functions like `getCourses()`, `getCourse(id)`, `createCourse(payload)`.
  - Use TypeScript types from `src/types` to represent DTOs.
  - Never expose raw HTTP details to components.

### Error handling
- API client should throw or return a structured error object (status code, message, optional details).
- UI components and feature components should:
  - Show clear error states (message, retry button if appropriate).
  - Not leak raw backend error payloads directly to the user.

## State management
- Prefer **React Query / TanStack Query** (if present) or simple hooks with local state for server state.
- Avoid global state libraries unless already in use (Redux, Zustand, etc.).
- For now, assume **no new global state manager** should be introduced without strong reason.

Guidelines:
- Server state (data coming from backend) should be cached and synced via React Query or equivalent.
- Local UI state lives inside components or small hooks.
- Do not manually implement ad-hoc caching; use the chosen data-fetching library.

## Routing
- Use React Router (if present) or the existing routing solution.
- Define routes in `src/app/router.tsx`.
- Route components should:
  - Fetch their own data via hooks.
  - Compose feature/UI components.
- Keep route configuration declarative and grouped by domain (auth, courses, admin, etc.).

## Testing
- Use **Vitest** or **Jest** (depending on what is already configured) for tests.
- Use **React Testing Library** for component testing.
- Prefer testing behavior and user-visible output over implementation details.

Guidelines:
- For new code, add tests for:
  - core feature components (screens, complex flows);
  - critical hooks (auth, course loading, enrollment actions);
  - utilities with non-trivial logic.
- Aim for **80%+** coverage on new feature modules.
- Mock API calls at the boundary (API client), not deep inside components.

## Integration with backend
- The Go backend defines the API contract.
- Do not invent fields, endpoints, or error shapes.
- If the backend API changes, update `src/types` and `src/api` first, then adjust features.
- Always keep TypeScript DTOs in sync with backend responses.

## Security & auth
- Never store tokens in localStorage or sessionStorage if the project uses HTTP-only cookies (check existing implementation).
- Do not expose sensitive data (tokens, internal IDs meant for backend only) in the UI.
- Protect routes on the frontend according to backend auth model (e.g., redirect to login if user is not authenticated).
- Do not implement your own crypto in the frontend.

## Performance considerations
- Avoid unnecessary re-renders:
  - Use memoization (`useMemo`, `useCallback`) where profiling shows benefit.
  - Split large components into smaller ones when appropriate.
- Avoid loading huge chunks on initial load:
  - Use code-splitting for big feature routes.
- Prefer server-driven pagination for large lists (courses, users, submissions, etc.).
- Do not fetch data that is not needed for the current screen.

## Styling
- Follow the existing styling approach in the repo (CSS Modules, Tailwind, etc.).
- Keep styles scoped to components.
- Avoid inline styles for anything non-trivial.
- Use design tokens / CSS variables if already present.

## Frontend change workflow for Codex
When implementing a feature or change:

1. **Understand the domain**
   - Identify which backend endpoints and types are involved.
   - Check existing feature implementations (auth, courses, etc.) and follow their patterns.

2. **Work from the API inward**
   - Add or update types in `src/types/*` as needed (matching backend DTOs).
   - Add or update API functions in `src/api/*`.

3. **Implement hooks and feature logic**
   - Create/update hooks in `src/lib/hooks` or the relevant feature folder.
   - Encapsulate data fetching and state handling.

4. **Compose UI**
   - Use `src/components/ui` primitives where possible.
   - Keep feature components focused and readable.

5. **Add tests**
   - Add or update tests for key components and hooks.
   - Use React Testing Library and mocks for API calls.

6. **Check behavior**
   - Ensure loading, empty, error, and success states are handled.
   - Validate that the UI respects backend constraints and domain rules.

## Things to avoid
- Do not call the backend directly from UI components using raw `fetch`.
- Do not bypass `src/api` when talking to the backend.
- Do not add new global state libraries without necessity.
- Do not use `any` or `unknown` without a strong reason and a comment.
- Do not duplicate types; reuse shared types from `src/types`.
- Do not embed business logic in dumb UI components.

## Expected output from Codex in frontend/
- Changes should be minimal, coherent, and consistent with existing patterns.
- Prefer extending existing patterns over inventing new ones.
- Explain trade-offs in comments only when a choice is non-obvious.
- When requirements are unclear, choose the safest maintainable default and state the assumption in code comments.
