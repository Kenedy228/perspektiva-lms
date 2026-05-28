# Диаграмма последовательности домена Question

Диаграмма отражает три ключевых сценария: создание попытки, добавление ответа и проверка/оценка ответа.

```mermaid
sequenceDiagram
    actor Student
    participant Transport as HTTP Handler
    participant App as GradeUseCase
    participant Repo as QuestionRepository
    participant Reg as CheckerRegistry
    participant Validator as AnswerValidator
    participant Checker as grading.Checker
    participant Domain as Domain Objects

    %% ═══════════════════════════════════════════
    %% СЦЕНАРИЙ 1: Создание попытки
    %% ═══════════════════════════════════════════
    rect rgb(235, 245, 255)
        Note over Student, Domain: Сценарий 1 — Создание попытки (Attempt.New)

        Student->>Transport: POST /attempts
        Transport->>Domain: attempt.New(params, startedAt)

        activate Domain
        Domain->>Domain: validateEnrollmentID()
        Domain->>Domain: validateQuizID()
        Domain->>Domain: validateQuestions([]Question)

        loop ∀ question ∈ params.Questions
            Domain->>Domain: item.New(question)
            Note over Domain: question.Clone() → snapshot
            Domain->>Domain: items = append(items, item)
        end

        Note over Domain: status = in_progress
        Note over Domain: answers = make(map[uuid]Entry)
        Domain-->>Transport: *Attempt
        deactivate Domain

        Transport-->>Student: 201 Created
    end

    %% ═══════════════════════════════════════════
    %% СЦЕНАРИЙ 2: Добавление ответа
    %% ═══════════════════════════════════════════
    rect rgb(255, 245, 235)
        Note over Student, Domain: Сценарий 2 — Добавление ответа (Attempt.AddAnswer)

        Student->>Transport: PUT /attempts/{id}/answers
        Transport->>Domain: attempt.AddAnswer(questionID, answer, at)

        activate Domain

        Domain->>Domain: CanModify()?
        alt status ≠ in_progress
            Domain-->>Transport: ❌ ErrStateConflict
            Transport-->>Student: 409 Conflict
        else status = in_progress
            Domain->>Domain: ensureBeforeDeadline(at)
            alt at > deadlineAt
                Domain-->>Transport: ❌ ErrStateConflict
                Transport-->>Student: 409 Conflict
            else within deadline
                Domain->>Domain: findItem(items, questionID)
                alt question not found
                    Domain-->>Transport: ❌ ErrNotFound
                    Transport-->>Student: 404 Not Found
                else question found
                    Domain->>Domain: answer.New(questionID, answer, at)
                    Note over Domain: answer.Clone() → immutable snapshot
                    Domain->>Domain: answers[questionID] = entry
                    Domain-->>Transport: ✔ ok
                    deactivate Domain
                    Transport-->>Student: 200 OK
                end
            end
        end
    end

    %% ═══════════════════════════════════════════
    %% СЦЕНАРИЙ 3: Проверка ответа
    %% ═══════════════════════════════════════════
    rect rgb(235, 255, 235)
        Note over Student, Domain: Сценарий 3 — Проверка ответа (GradeUseCase / ValidateAnswerUseCase)

        Student->>Transport: POST /answers/validate
        Transport->>App: ValidateAnswerUseCase.Execute(ValidateAnswerInput)

        activate App
        App->>App: validateInput()
        Note over App: questionID не пустой, answer не nil

        App->>Repo: FindByID(ctx, questionID)
        Repo-->>App: question.Question

        App->>App: validators[q.Type()]
        App->>Validator: Validate(question, answer)

        activate Validator

        rect rgb(250, 250, 240)
            Note over Validator: Runtime type assertions

            alt TypeSelectable
                Validator->>Validator: q.(*selectable.Question) ?
                Validator->>Validator: a.(selectable/answer.Answer) ?
            else TypeMatching
                Validator->>Validator: q.(*matching.Question) ?
                Validator->>Validator: a.(matching/answer.Answer) ?
            else TypeSequence
                Validator->>Validator: q.(*sequence.Question) ?
                Validator->>Validator: a.(sequence/answer.Answer) ?
            else TypeShort
                Validator->>Validator: q.(*short.Question) ?
                Validator->>Validator: a.(short/answer.Answer) ?
            end

            alt types match
                Validator-->>App: nil (ok)
            else type mismatch or nil
                Validator-->>App: ErrInvalidQuestionType / ErrNilAnswer
                deactivate Validator
                deactivate App
                App-->>Transport: error
                Transport-->>Student: 400 Bad Request
            end
        end

        deactivate Validator
        Note over App: ValidateAnswerUseCase returns nil (valid)
        deactivate App

        App-->>Transport: nil
        Transport-->>Student: 204 No Content
    end

    rect rgb(240, 255, 240)
        Note over Student, Domain: Сценарий 3b — Оценка ответа (GradeUseCase с Registry)

        Student->>Transport: POST /answers/grade
        Transport->>App: GradeUseCase.Execute(GradeInput)

        activate App
        App->>App: validateInput()
        App->>Repo: FindByID(ctx, questionID)
        Repo-->>App: question.Question

        App->>App: validators[q.Type()]
        App->>Validator: Validate(question, answer)
        Validator-->>App: nil

        App->>Reg: Get(q.Type())
        Reg-->>App: grading.Checker

        App->>Checker: Check(question, answer)

        activate Checker

        rect rgb(250, 250, 240)
            Note over Checker: Runtime type assertions (defensive)

            alt TypeSelectable
                Checker->>Checker: qCast := q.(*selectable.Question)
                Checker->>Checker: aCast := a.(selectable/answer.Answer)
                Checker->>Checker: exactly all correct ∧ no incorrect?
                Checker-->>App: Score(0.0 | 1.0)
            else TypeMatching
                Checker->>Checker: qCast := q.(*matching.Question)
                Checker->>Checker: aCast := a.(matching/answer.Answer)
                Checker->>Checker: all correct pairs ∈ student answer?
                Checker-->>App: Score(0.0 | 1.0)
            else TypeSequence
                Checker->>Checker: qCast := q.(*sequence.Question)
                Checker->>Checker: aCast := a.(sequence/answer.Answer)
                Checker->>Checker: SHA1(value) match for all positions?
                Checker-->>App: Score(0.0 | 1.0)
            else TypeShort
                Checker->>Checker: qCast := q.(*short.Question)
                Checker->>Checker: aCast := a.(short/answer.Answer)
                Checker->>Checker: normalize(answer) ∈ normalize(variants)?
                Checker-->>App: Score(0.0 | 1.0)
            end
        end

        deactivate Checker
        deactivate App

        App-->>Transport: GradeOutput{Score}
        Transport-->>Student: 200 { score: 0 | 1 }
    end

    %% ═══════════════════════════════════════════
    %% СЦЕНАРИЙ 4: Завершение попытки
    %% ═══════════════════════════════════════════
    rect rgb(245, 235, 255)
        Note over Student, Domain: Сценарий 4 — Завершение попытки (Attempt.Finish)

        Student->>Transport: POST /attempts/{id}/finish
        Transport->>Domain: attempt.Finish(at)

        activate Domain
        Domain->>Domain: status = in_progress?
        alt status ≠ in_progress
            Domain-->>Transport: ❌ ErrStateConflict
            Transport-->>Student: 409 Conflict
        else status = in_progress
            Domain->>Domain: at.IsZero()?
            Domain->>Domain: at.Before(startedAt)?
            Domain->>Domain: ensureBeforeDeadline(at)

            Domain->>Domain: status = finished
            Domain->>Domain: finishedAt = at
            Domain-->>Transport: ✔ ok
            deactivate Domain
            Transport-->>Student: 200 OK
        end
    end
```

---

## Сценарии

### 1. Создание попытки (`Attempt.New`)
- Принимает параметры: `enrollmentID`, `quizID`, список вопросов, ограничение по времени.
- Каждый `Question` клонируется через `Clone()` и сохраняется как `Item` — это **immutable snapshot**.
- Устанавливается начальный статус `in_progress` и дедлайн (если задан лимит времени).
- Ответы инициализируются пустой картой.

### 2. Добавление ответа (`Attempt.AddAnswer`)
- Проверяется, что попытка в статусе `in_progress` (`CanModify()`).
- Проверяется, что текущее время не превысило дедлайн.
- Проверяется, что вопрос с данным `questionID` существует среди `items`.
- Ответ клонируется через `Answer.Clone()` и сохраняется в `answers` как `Entry`.

### 3. Проверка / оценка ответа (`GradeUseCase` / `ValidateAnswerUseCase`)

**ValidateAnswerUseCase (только валидация):**
- Загружает вопрос по ID через `Repository`.
- Находит `AnswerValidator` по `q.Type()` через `validators` map.
- Вызывает `Validator.Validate(q, a)`, который проверяет:
  - `q != nil`, `a != nil`;
  - тип ответа соответствует типу вопроса (через runtime type assertions).
- **Не выполняет** вычисление score. Возвращает только ошибку или nil.

**GradeUseCase (валидация + оценка):**
- Загружает вопрос, валидирует ответ через `AnswerValidator`.
- Получает `Checker` через `Registry.Get(q.Type())` — O(1) map lookup.
- Вызывает `Checker.Check(q, a)`, который выполняет **defensive type assertions** и вычисляет score.
- Логика проверки специфична для каждого типа:
  - **selectable**: студент должен выбрать ровно все правильные варианты (без лишних);
  - **matching**: все правильные пары должны присутствовать в ответе студента;
  - **sequence**: полное совпадение длины и порядка (ID опций вычисляются через `uuid.NewSHA1`);
  - **short**: нормализованный ответ студента должен совпасть с любым вариантом.

### 4. Завершение попытки (`Attempt.Finish`)
- Проверяется статус `in_progress`, валидность времени и дедлайн.
- Статус меняется на `finished`, фиксируется `finishedAt`.
- Альтернативные переходы: `SetExpired()`, `Interrupt()`, `Cancel()`.
