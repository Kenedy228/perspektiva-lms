# HTTP Application Routes

All non-auth application routes require:

```http
Authorization: Bearer <session-token>
Content-Type: application/json
```

Responses use the shared envelope:

```json
{
  "data": {},
  "_links": {
    "self": {"href": "/resource", "method": "GET"}
  }
}
```

Errors use:

```json
{
  "error": {
    "status": 400,
    "code": "invalid_input",
    "message": "request body is invalid",
    "details": {},
    "_links": {
      "self": {"href": "/resource", "method": "POST"}
    }
  }
}
```

## Routes

| Method | Path | Purpose |
| --- | --- | --- |
| GET | `/` | API index with HATEOAS links |
| GET | `/healthz` | Health check |
| POST | `/auth/login` | Create bearer session |
| POST | `/auth/logout` | Acknowledge client-side token disposal |
| GET | `/auth/session` | Return current session claims |
| GET | `/accounts` | List accounts |
| POST | `/accounts` | Create account |
| GET | `/accounts/{id}` | Get account |
| PATCH | `/accounts/{id}/login` | Change account login |
| PATCH | `/accounts/{id}/password` | Change account password |
| PATCH | `/accounts/{id}/role` | Change account role |
| POST | `/accounts/{id}/block` | Block account |
| POST | `/accounts/{id}/activate` | Activate account |
| DELETE | `/accounts/{id}` | Delete account |
| GET | `/organizations` | List organizations |
| POST | `/organizations` | Create organization |
| GET | `/organizations/{id}` | Get organization |
| PATCH | `/organizations/{id}` | Rename organization |
| PATCH | `/organizations/{id}/inn` | Change organization INN |
| DELETE | `/organizations/{id}` | Delete organization |
| GET | `/persons` | List persons |
| POST | `/persons` | Create person |
| GET | `/persons/{id}` | Get person |
| PATCH | `/persons/{id}` | Rename person |
| PUT | `/persons/{id}/profile` | Replace person profile |
| PATCH | `/persons/{id}/profile` | Change profile fields |
| DELETE | `/persons/{id}/profile` | Detach profile |
| PUT | `/persons/{id}/organization` | Assign organization |
| DELETE | `/persons/{id}/organization` | Remove organization |
| DELETE | `/persons/{id}` | Delete person |
| GET | `/banks` | List question banks |
| POST | `/banks` | Create question bank |
| GET | `/banks/{id}` | Get bank details |
| PATCH | `/banks/{id}` | Rename bank |
| POST | `/banks/{id}/questions` | Add questions to bank |
| DELETE | `/banks/{id}/questions` | Remove questions from bank |
| POST | `/questions` | Create question |
| GET | `/questions/{id}` | Get question |
| PATCH | `/questions/{id}` | Change question title |
| PUT | `/questions/{id}/content` | Replace question content |
| PUT | `/questions/{id}/attachment` | Replace question attachment |
| DELETE | `/questions/{id}/attachment` | Remove attachment |
| POST | `/questions/{id}/grade` | Grade one answer |
| POST | `/quizzes` | Create quiz |
| GET | `/quizzes/{id}` | Get quiz |
| PATCH | `/quizzes/{id}` | Rename quiz |
| PATCH | `/quizzes/{id}/limits` | Change quiz limits |
| PATCH | `/quizzes/{id}/shuffle` | Change shuffle policy |
| PUT | `/quizzes/{id}/sources` | Replace quiz sources |
| DELETE | `/quizzes/{id}` | Delete quiz |
| POST | `/attempts` | Start attempt |
| GET | `/attempts/{id}` | Get attempt |
| PUT | `/attempts/{id}/answers/{questionID}` | Add or replace answer |
| POST | `/attempts/{id}/finish` | Finish attempt |
| POST | `/attempts/{id}/cancel` | Cancel attempt |
| GET | `/courses` | List manageable or enrolled courses |
| POST | `/courses` | Create course |
| GET | `/courses/{id}` | Get course details |
| PATCH | `/courses/{id}` | Rename course |
| POST | `/courses/{courseID}/blocks` | Add block to course |
| POST | `/blocks/{blockID}/elements` | Add element to block |
| POST | `/courses/{courseID}/progress` | Mark course progress |
| GET | `/courses/{id}/ratings` | List course ratings |
| POST | `/enrollments` | Enroll student in version |
| GET | `/statistics/students` | List student statistics |

## Correct Examples

Create organization:

```http
POST /organizations
Authorization: Bearer <admin-token>
Content-Type: application/json

{"name":"Training Department","inn":"2536000000","inn_type":"legal_entity"}
```

```json
{
  "data": {"id": "8d25c84c-4eb7-4e1e-9fd0-8dbd27408d2d"},
  "_links": {
    "self": {"href": "/organizations/8d25c84c-4eb7-4e1e-9fd0-8dbd27408d2d", "method": "GET"},
    "collection": {"href": "/organizations", "method": "GET"}
  }
}
```

Create quiz with random source and shuffle:

```json
{
  "title": "Safety basics",
  "max_attempts": 2,
  "time_limit_seconds": 1800,
  "shuffle_questions": true,
  "sources": [
    {"bank_id": "3c91e092-823e-4768-9f69-63f5d83479b6", "criteria_type": "random", "question_count": 10}
  ]
}
```

Start attempt:

```json
{
  "account_id": "c62a4be7-9cb0-4821-884e-a16778feec03",
  "enrollment_id": "05ab0173-7e7d-43f9-980e-727244968884",
  "quiz_id": "734c06f8-55d9-4bb5-b502-ef104ac9b06d"
}
```

## Incorrect Examples

Missing token:

```json
{
  "error": {
    "status": 401,
    "code": "unauthorized",
    "message": "bearer token is required",
    "_links": {
      "self": {"href": "/courses", "method": "GET"}
    }
  }
}
```

Malformed JSON:

```http
POST /courses
Authorization: Bearer <admin-token>
Content-Type: application/json

{"title":
```

```json
{
  "error": {
    "status": 400,
    "code": "invalid_json",
    "message": "request body is invalid",
    "_links": {
      "self": {"href": "/courses", "method": "POST"}
    }
  }
}
```

Quiz source asks for more questions than the bank contains:

```json
{
  "error": {
    "status": 400,
    "code": "invalid_input",
    "message": "quiz source requests 25 questions, but bank contains 10",
    "_links": {
      "self": {"href": "/quizzes", "method": "POST"}
    }
  }
}
```
