package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	courseports "gitflic.ru/lms/backend/internal/application/ports/course"
	coursecommands "gitflic.ru/lms/backend/internal/application/usecases/course/commands"
	coursequeries "gitflic.ru/lms/backend/internal/application/usecases/course/queries"
	"gitflic.ru/lms/backend/internal/domain/role"
	"gitflic.ru/lms/backend/internal/transport/http/middleware"
	"gitflic.ru/lms/backend/internal/transport/http/session"
	"github.com/google/uuid"
)

func newTestSession() *session.Manager {
	return session.NewManager([]byte("test-secret"), time.Hour)
}

func authToken(mgr *session.Manager, r role.Role) string {
	token, _, err := mgr.Issue(uuid.New(), uuid.New(), r)
	if err != nil {
		panic("issue test token: " + err.Error())
	}
	return token
}

func testAPI() *API {
	return &API{
		Courses: CourseUseCases{
			Create:               &stubCreateCourse{id: uuid.New()},
			Rename:               &stubRenameCourse{},
			AddBlock:             &stubAddBlock{id: uuid.New()},
			RemoveBlock:          &stubRemoveBlock{},
			MoveBlock:            &stubMoveBlock{},
			AddElement:           &stubAddElement{id: uuid.New()},
			RemoveElement:        &stubRemoveElement{},
			MoveElement:          &stubMoveElement{},
			ChangeCompletionMode: &stubChangeCompletionMode{},
			Progress:             &stubProgress{},
			UnmarkProgress:       &stubUnmarkProgress{},
			GetProgress:          &stubGetProgress{},
			List:                 &stubListQuery{},
			Ratings:              &stubRatingsQuery{},
			Statistics:           &stubStatisticsQuery{},
			Query:                &stubQueryService{},
		},
	}
}

func doRequest(api *API, mgr *session.Manager, method, path, body string, handler func(*API) http.HandlerFunc) *httptest.ResponseRecorder {
	var req *http.Request
	if body != "" {
		req = httptest.NewRequest(method, path, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, path, nil)
	}
	req.Header.Set("Authorization", "Bearer "+authToken(mgr, role.NewAdmin()))
	req.SetPathValue("id", "00000000-0000-0000-0000-000000000001")
	req.SetPathValue("courseID", "00000000-0000-0000-0000-000000000001")
	req.SetPathValue("blockID", "00000000-0000-0000-0000-000000000002")
	req.SetPathValue("elementID", "00000000-0000-0000-0000-000000000003")

	rec := httptest.NewRecorder()
	middleware.Auth(mgr)(handler(api))(rec, req)
	return rec
}

func doUnauthRequest(api *API, method, path string, handler func(*API) http.HandlerFunc) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, nil)
	rec := httptest.NewRecorder()
	handler(api)(rec, req)
	return rec
}

func assertStatus(t *testing.T, rec *httptest.ResponseRecorder, want int) {
	t.Helper()
	if rec.Code != want {
		t.Fatalf("expected status %d, got %d: %s", want, rec.Code, rec.Body.String())
	}
}

func assertJSONField(t *testing.T, rec *httptest.ResponseRecorder, field, want string) {
	t.Helper()
	var envelope struct {
		Data map[string]string `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &envelope); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if envelope.Data[field] != want {
		t.Fatalf("expected %q=%q, got %q", field, want, envelope.Data[field])
	}
}

// --- stubs ---

type stubCreateCourse struct{ id uuid.UUID }
type stubRenameCourse struct{}
type stubAddBlock struct{ id uuid.UUID }
type stubRemoveBlock struct{}
type stubMoveBlock struct{}
type stubAddElement struct{ id uuid.UUID }
type stubRemoveElement struct{}
type stubMoveElement struct{}
type stubChangeCompletionMode struct{}
type stubProgress struct{}
type stubUnmarkProgress struct{}
type stubGetProgress struct{}
type stubListQuery struct{}
type stubRatingsQuery struct{}
type stubStatisticsQuery struct{}
type stubQueryService struct{}

func (s *stubCreateCourse) Execute(_ context.Context, in coursecommands.CreateCourseInput) (*coursecommands.Output, error) {
	if in.Title == "" {
		return nil, errors.New("title required")
	}
	return &coursecommands.Output{ID: s.id.String()}, nil
}

func (s *stubRenameCourse) Execute(_ context.Context, _ coursecommands.RenameCourseInput) error {
	return nil
}

func (s *stubAddBlock) Execute(_ context.Context, _ coursecommands.AddBlockToCourseInput) (*coursecommands.Output, error) {
	return &coursecommands.Output{ID: s.id.String()}, nil
}

func (s *stubRemoveBlock) Execute(_ context.Context, _ coursecommands.RemoveBlockFromCourseInput) error {
	return nil
}

func (s *stubMoveBlock) Execute(_ context.Context, _ coursecommands.MoveCourseBlockInput) error {
	return nil
}

func (s *stubAddElement) Execute(_ context.Context, _ coursecommands.AddElementToBlockInput) (*coursecommands.Output, error) {
	return &coursecommands.Output{ID: s.id.String()}, nil
}

func (s *stubRemoveElement) Execute(_ context.Context, _ coursecommands.RemoveElementFromBlockInput) error {
	return nil
}

func (s *stubMoveElement) Execute(_ context.Context, _ coursecommands.MoveBlockElementInput) error {
	return nil
}

func (s *stubChangeCompletionMode) Execute(_ context.Context, _ coursecommands.ChangeElementCompletionModeInput) error {
	return nil
}

func (s *stubProgress) Execute(_ context.Context, _ coursecommands.MarkProgressInput) error {
	return nil
}

func (s *stubUnmarkProgress) Execute(_ context.Context, _ coursecommands.UnmarkElementCompletedInput) error {
	return nil
}

func (s *stubGetProgress) Execute(_ context.Context, _ coursecommands.GetProgressInput) (*coursecommands.ProgressOutput, error) {
	return &coursecommands.ProgressOutput{
		CompletedCount:      3,
		Percent:             75,
		CompletedElementIDs: []uuid.UUID{uuid.New(), uuid.New(), uuid.New()},
	}, nil
}

func (s *stubListQuery) Execute(_ context.Context, _ coursequeries.ListInput) (*coursequeries.ListOutput, error) {
	return &coursequeries.ListOutput{}, nil
}

func (s *stubRatingsQuery) Execute(_ context.Context, _ coursequeries.RatingsInput) (*coursequeries.RatingsOutput, error) {
	return &coursequeries.RatingsOutput{}, nil
}

func (s *stubStatisticsQuery) Execute(_ context.Context, _ coursequeries.StudentStatisticsInput) (*coursequeries.StudentStatisticsOutput, error) {
	return &coursequeries.StudentStatisticsOutput{}, nil
}

func (s *stubQueryService) ListManageable(_ context.Context, _ courseports.Filter) ([]courseports.ShortView, error) {
	return nil, nil
}
func (s *stubQueryService) ListVisibleForStudent(_ context.Context, _ courseports.Filter) ([]courseports.ShortView, error) {
	return nil, nil
}
func (s *stubQueryService) GetDetailsByID(_ context.Context, _ uuid.UUID) (courseports.DetailedView, error) {
	return courseports.DetailedView{}, nil
}
func (s *stubQueryService) ListRatings(_ context.Context, _ uuid.UUID, _, _ int) ([]courseports.StudentRatingView, error) {
	return nil, nil
}
func (s *stubQueryService) ListStudentStatistics(_ context.Context, _ courseports.StudentStatisticsFilter) ([]courseports.StudentRatingView, error) {
	return nil, nil
}

// --- tests ---

func TestListCoursesUnauthorized(t *testing.T) {
	rec := doUnauthRequest(testAPI(), http.MethodGet, "/courses", func(a *API) http.HandlerFunc { return a.ListCourses })
	assertStatus(t, rec, http.StatusUnauthorized)
}

func TestCreateCourseUnauthorized(t *testing.T) {
	rec := doUnauthRequest(testAPI(), http.MethodPost, "/courses", func(a *API) http.HandlerFunc { return a.CreateCourse })
	assertStatus(t, rec, http.StatusUnauthorized)
}

func TestCreateCourseSuccess(t *testing.T) {
	api := testAPI()
	mgr := newTestSession()
	rec := doRequest(api, mgr, http.MethodPost, "/courses", `{"title":"Test Course"}`, func(a *API) http.HandlerFunc { return a.CreateCourse })
	assertStatus(t, rec, http.StatusCreated)
	assertJSONField(t, rec, "id", api.Courses.Create.(*stubCreateCourse).id.String())
}

func TestAddCourseBlockUnauthorized(t *testing.T) {
	rec := doUnauthRequest(testAPI(), http.MethodPost, "/courses/id/blocks", func(a *API) http.HandlerFunc { return a.AddCourseBlock })
	assertStatus(t, rec, http.StatusUnauthorized)
}

func TestAddCourseBlockSuccess(t *testing.T) {
	api := testAPI()
	mgr := newTestSession()
	rec := doRequest(api, mgr, http.MethodPost, "/courses/1/blocks", `{"title":"Block"}`, func(a *API) http.HandlerFunc { return a.AddCourseBlock })
	assertStatus(t, rec, http.StatusCreated)
}

func TestRemoveCourseBlockUnauthorized(t *testing.T) {
	rec := doUnauthRequest(testAPI(), http.MethodDelete, "/courses/id/blocks/bid", func(a *API) http.HandlerFunc { return a.RemoveCourseBlock })
	assertStatus(t, rec, http.StatusUnauthorized)
}

func TestRemoveCourseBlockSuccess(t *testing.T) {
	api := testAPI()
	mgr := newTestSession()
	rec := doRequest(api, mgr, http.MethodDelete, "/courses/1/blocks/2", "", func(a *API) http.HandlerFunc { return a.RemoveCourseBlock })
	assertStatus(t, rec, http.StatusNoContent)
}

func TestMoveCourseBlockUnauthorized(t *testing.T) {
	rec := doUnauthRequest(testAPI(), http.MethodPatch, "/courses/id/blocks/move", func(a *API) http.HandlerFunc { return a.MoveCourseBlock })
	assertStatus(t, rec, http.StatusUnauthorized)
}

func TestMoveCourseBlockSuccess(t *testing.T) {
	api := testAPI()
	mgr := newTestSession()
	rec := doRequest(api, mgr, http.MethodPatch, "/courses/1/blocks/move", `{"from":0,"to":1}`, func(a *API) http.HandlerFunc { return a.MoveCourseBlock })
	assertStatus(t, rec, http.StatusOK)
}

func TestAddBlockElementUnauthorized(t *testing.T) {
	rec := doUnauthRequest(testAPI(), http.MethodPost, "/blocks/id/elements", func(a *API) http.HandlerFunc { return a.AddBlockElement })
	assertStatus(t, rec, http.StatusUnauthorized)
}

func TestRemoveBlockElementUnauthorized(t *testing.T) {
	rec := doUnauthRequest(testAPI(), http.MethodDelete, "/blocks/id/elements/eid", func(a *API) http.HandlerFunc { return a.RemoveBlockElement })
	assertStatus(t, rec, http.StatusUnauthorized)
}

func TestRemoveBlockElementSuccess(t *testing.T) {
	api := testAPI()
	mgr := newTestSession()
	rec := doRequest(api, mgr, http.MethodDelete, "/blocks/1/elements/2", "", func(a *API) http.HandlerFunc { return a.RemoveBlockElement })
	assertStatus(t, rec, http.StatusNoContent)
}

func TestMoveBlockElementUnauthorized(t *testing.T) {
	rec := doUnauthRequest(testAPI(), http.MethodPatch, "/blocks/id/elements/move", func(a *API) http.HandlerFunc { return a.MoveBlockElement })
	assertStatus(t, rec, http.StatusUnauthorized)
}

func TestMoveBlockElementSuccess(t *testing.T) {
	api := testAPI()
	mgr := newTestSession()
	rec := doRequest(api, mgr, http.MethodPatch, "/blocks/1/elements/move", `{"from":1,"to":0}`, func(a *API) http.HandlerFunc { return a.MoveBlockElement })
	assertStatus(t, rec, http.StatusOK)
}

func TestChangeElementCompletionModeUnauthorized(t *testing.T) {
	rec := doUnauthRequest(testAPI(), http.MethodPatch, "/elements/id/completion-mode", func(a *API) http.HandlerFunc { return a.ChangeElementCompletionMode })
	assertStatus(t, rec, http.StatusUnauthorized)
}

func TestChangeElementCompletionModeSuccess(t *testing.T) {
	api := testAPI()
	mgr := newTestSession()
	rec := doRequest(api, mgr, http.MethodPatch, "/elements/1/completion-mode", `{"completion_mode":"manual"}`, func(a *API) http.HandlerFunc { return a.ChangeElementCompletionMode })
	assertStatus(t, rec, http.StatusOK)
}

func TestMarkProgressUnauthorized(t *testing.T) {
	rec := doUnauthRequest(testAPI(), http.MethodPost, "/courses/id/progress", func(a *API) http.HandlerFunc { return a.MarkCourseProgress })
	assertStatus(t, rec, http.StatusUnauthorized)
}

func TestUnmarkProgressUnauthorized(t *testing.T) {
	rec := doUnauthRequest(testAPI(), http.MethodDelete, "/courses/id/progress", func(a *API) http.HandlerFunc { return a.UnmarkCourseProgress })
	assertStatus(t, rec, http.StatusUnauthorized)
}

func TestUnmarkProgressSuccess(t *testing.T) {
	api := testAPI()
	mgr := newTestSession()
	rec := doRequest(api, mgr, http.MethodDelete, "/courses/1/progress", `{"enrollment_id":"00000000-0000-0000-0000-000000000010","element_id":"00000000-0000-0000-0000-000000000011"}`, func(a *API) http.HandlerFunc { return a.UnmarkCourseProgress })
	assertStatus(t, rec, http.StatusNoContent)
}

func TestGetProgressUnauthorized(t *testing.T) {
	rec := doUnauthRequest(testAPI(), http.MethodGet, "/courses/id/progress", func(a *API) http.HandlerFunc { return a.GetCourseProgress })
	assertStatus(t, rec, http.StatusUnauthorized)
}

func TestGetProgressSuccess(t *testing.T) {
	api := testAPI()
	mgr := newTestSession()
	rec := doRequest(api, mgr, http.MethodGet, "/courses/1/progress?enrollment_id=00000000-0000-0000-0000-000000000010", "", func(a *API) http.HandlerFunc { return a.GetCourseProgress })
	assertStatus(t, rec, http.StatusOK)
}

func TestListRatingsUnauthorized(t *testing.T) {
	rec := doUnauthRequest(testAPI(), http.MethodGet, "/courses/id/ratings", func(a *API) http.HandlerFunc { return a.ListCourseRatings })
	assertStatus(t, rec, http.StatusUnauthorized)
}
