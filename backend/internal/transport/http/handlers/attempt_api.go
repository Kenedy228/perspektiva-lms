package handlers

import (
	"net/http"
	"time"

	attemptcommands "gitflic.ru/lms/backend/internal/application/usecases/attempt/commands"
	attemptdomain "gitflic.ru/lms/backend/internal/domain/attempt"
	"gitflic.ru/lms/backend/internal/domain/grading/registry"
	questdomain "gitflic.ru/lms/backend/internal/domain/question"
	matchingq "gitflic.ru/lms/backend/internal/domain/question/matching"
	selectableq "gitflic.ru/lms/backend/internal/domain/question/selectable"
	sequenceq "gitflic.ru/lms/backend/internal/domain/question/sequence"
	"gitflic.ru/lms/backend/internal/domain/role"
	"gitflic.ru/lms/backend/internal/transport/http/response"
	"github.com/google/uuid"
)

func (api *API) ListAttempts(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorRole(r)
	if !ok {
		response.WriteError(w, r, response.NewError(http.StatusUnauthorized, "unauthorized", "session is required"))
		return
	}
	if actor.role.Kind() != role.TypeAdmin && actor.role.Kind() != role.TypeCreator {
		response.WriteError(w, r, response.NewError(http.StatusForbidden, "forbidden", "Операция запрещена."))
		return
	}
	enrollmentIDStr := r.URL.Query().Get("enrollment_id")
	if enrollmentIDStr == "" {
		response.WriteError(w, r, response.NewError(http.StatusBadRequest, "invalid_input", "enrollment_id обязателен."))
		return
	}
	enrollmentID, err := uuid.Parse(enrollmentIDStr)
	if err != nil {
		response.WriteError(w, r, response.NewError(http.StatusBadRequest, "invalid_input", "enrollment_id должен быть UUID."))
		return
	}
	limit, offset := limitOffset(r)
	views, err := api.Attempts.Query.ListByEnrollmentID(r.Context(), enrollmentID, limit, offset)
	if err != nil {
		writeHandlerError(w, r, err)
		return
	}

	items := make([]map[string]any, 0, len(views))
	for _, v := range views {
		item := map[string]any{
			"id":              v.ID,
			"enrollment_id":   v.EnrollmentID,
			"quiz_id":         v.QuizID,
			"status":          v.Status,
			"started_at":      v.StartedAt,
			"questions_count": v.QuestionsCount,
			"answers_count":   v.AnswersCount,
		}
		if !v.DeadlineAt.IsZero() {
			item["deadline_at"] = v.DeadlineAt
		}
		if !v.FinishedAt.IsZero() {
			item["finished_at"] = v.FinishedAt
		}
		items = append(items, item)
	}
	writeOK(w, r, items, nil)
}

func (api *API) StartAttempt(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorRole(r)
	if !ok {
		response.WriteError(w, r, response.NewError(http.StatusUnauthorized, "unauthorized", "session is required"))
		return
	}
	var req AttemptStartRequest
	if err := response.DecodeJSON(r, &req); err != nil {
		response.WriteError(w, r, response.NewError(http.StatusBadRequest, "invalid_json", "request body is invalid"))
		return
	}
	out, err := api.Attempts.Start.Execute(r.Context(), attemptcommands.StartInput{
		ActorRole:    actor.role,
		AccountID:    req.AccountID,
		EnrollmentID: req.EnrollmentID,
		QuizID:       req.QuizID,
		Now:          time.Now().UTC(),
	})
	if err != nil {
		writeHandlerError(w, r, err)
		return
	}
	writeCreated(w, "/attempts", out.ID)
}

func (api *API) GetAttempt(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		writeHandlerError(w, r, err)
		return
	}
	a, err := api.Attempts.Repository.FindByID(r.Context(), id)
	if err != nil {
		writeHandlerError(w, r, err)
		return
	}

	answeredIDs := make([]string, 0, a.CountAnswers())
	for qID := range a.Answers() {
		answeredIDs = append(answeredIDs, qID.String())
	}

	data := map[string]any{
		"id":                    a.ID().String(),
		"enrollment_id":         a.EnrollmentID().String(),
		"quiz_id":               a.QuizID().String(),
		"status":                a.Status().String(),
		"started_at":            a.StartedAt(),
		"deadline_at":           a.DeadlineAt(),
		"finished_at":           a.FinishedAt(),
		"questions_count":       a.CountItems(),
		"answers_count":         a.CountAnswers(),
		"answered_question_ids": answeredIDs,
		"questions":             serializeAttemptQuestions(a),
	}

	if a.Status() != attemptdomain.StatusInProgress && api.Attempts.GradeRegistry != nil {
		questionScores, totalScore := computeAttemptScores(a, api.Attempts.GradeRegistry)
		data["question_scores"] = questionScores
		data["total_score"] = totalScore
	}

	writeOK(w, r, data, response.Links{
		"self":   {Href: r.URL.Path, Method: http.MethodGet},
		"finish": {Href: r.URL.Path + "/finish", Method: http.MethodPost},
		"cancel": {Href: r.URL.Path + "/cancel", Method: http.MethodPost},
	})
}

func serializeAttemptQuestions(a *attemptdomain.Attempt) []map[string]any {
	items := a.Items()
	result := make([]map[string]any, 0, len(items))
	for _, itm := range items {
		q := itm.Snapshot()
		qMap := map[string]any{
			"id":          q.ID().String(),
			"type":        q.Type().String(),
			"title":       q.Title().Value(),
			"instruction": q.Instruction(),
		}
		switch typed := q.(type) {
		case *selectableq.Question:
			opts := make([]map[string]any, 0, len(typed.Options()))
			for _, o := range typed.Options() {
				opts = append(opts, map[string]any{
					"id":   o.ID().String(),
					"text": o.Value(),
				})
			}
			qMap["options"] = opts
		case *sequenceq.Question:
			opts := make([]map[string]any, 0, len(typed.Options()))
			for _, o := range typed.Options() {
				opts = append(opts, map[string]any{
					"id":   uuid.NewSHA1(uuid.NameSpaceOID, []byte(o.Value())).String(),
					"text": o.Value(),
				})
			}
			qMap["options"] = opts
		case *matchingq.Question:
			pairs := make([]map[string]any, 0, len(typed.Pairs()))
			for _, p := range typed.Pairs() {
				pairs = append(pairs, map[string]any{
					"prompt_id": p.PromptID().String(),
					"prompt":    p.Prompt().Value(),
					"match_id":  p.MatchID().String(),
					"match":     p.Match().Value(),
				})
			}
			qMap["pairs"] = pairs
		default:
			_ = questdomain.TypeShort
		}
		result = append(result, qMap)
	}
	return result
}

func (api *API) AddAttemptAnswer(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorRole(r)
	if !ok {
		response.WriteError(w, r, response.NewError(http.StatusUnauthorized, "unauthorized", "session is required"))
		return
	}
	var req AttemptAnswerRequest
	if err := response.DecodeJSON(r, &req); err != nil {
		response.WriteError(w, r, response.NewError(http.StatusBadRequest, "invalid_json", "request body is invalid"))
		return
	}
	ans, err := buildAnswer(req)
	if err != nil {
		writeHandlerError(w, r, err)
		return
	}
	err = api.Attempts.Answer.Execute(r.Context(), attemptcommands.AddAnswerInput{
		ActorRole:  actor.role,
		AttemptID:  r.PathValue("id"),
		QuestionID: r.PathValue("questionID"),
		Answer:     ans,
		AnsweredAt: time.Now().UTC(),
	})
	if err != nil {
		writeHandlerError(w, r, err)
		return
	}
	writeOK(w, r, map[string]string{"id": r.PathValue("id")}, nil)
}

func (api *API) FinishAttempt(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorRole(r)
	if !ok {
		response.WriteError(w, r, response.NewError(http.StatusUnauthorized, "unauthorized", "session is required"))
		return
	}
	if err := api.Attempts.Finish.Execute(r.Context(), attemptcommands.FinishInput{ActorRole: actor.role, AttemptID: r.PathValue("id"), FinishedAt: time.Now().UTC()}); err != nil {
		writeHandlerError(w, r, err)
		return
	}
	writeOK(w, r, map[string]string{"id": r.PathValue("id")}, nil)
}

func (api *API) CancelAttempt(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorRole(r)
	if !ok {
		response.WriteError(w, r, response.NewError(http.StatusUnauthorized, "unauthorized", "session is required"))
		return
	}
	if err := api.Attempts.Cancel.Execute(r.Context(), attemptcommands.CancelInput{ActorRole: actor.role, AttemptID: r.PathValue("id"), CancelledAt: time.Now().UTC()}); err != nil {
		writeHandlerError(w, r, err)
		return
	}
	writeOK(w, r, map[string]string{"id": r.PathValue("id")}, nil)
}

// computeAttemptScores grades each answered question using the question snapshot
// already stored in the attempt, avoiding an extra DB round-trip.
// Returns per-question scores (0.0–1.0) and the overall average score.
func computeAttemptScores(a *attemptdomain.Attempt, reg *registry.Registry) (questionScores map[string]float64, totalScore float64) {
	questionByID := make(map[uuid.UUID]questdomain.Question, len(a.Items()))
	for _, itm := range a.Items() {
		questionByID[itm.ID()] = itm.Snapshot()
	}

	questionScores = make(map[string]float64, len(a.Answers()))
	var sum float64
	for questionID, entry := range a.Answers() {
		q, ok := questionByID[questionID]
		if !ok {
			continue
		}
		checker, err := reg.Get(q.Type())
		if err != nil {
			continue
		}
		s, err := checker.Check(q, entry.Answer())
		if err != nil {
			continue
		}
		questionScores[questionID.String()] = s.Value()
		sum += s.Value()
	}

	if n := len(a.Items()); n > 0 {
		totalScore = sum / float64(n)
	}
	return questionScores, totalScore
}
