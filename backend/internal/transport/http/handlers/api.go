package handlers

import (
	"context"
	"net/http"

	attemptports "gitflic.ru/lms/backend/internal/application/ports/attempt"
	courseports "gitflic.ru/lms/backend/internal/application/ports/course"
	quizports "gitflic.ru/lms/backend/internal/application/ports/quiz"
	accountcommands "gitflic.ru/lms/backend/internal/application/usecases/account/commands"
	accountqueries "gitflic.ru/lms/backend/internal/application/usecases/account/queries"
	attemptcommands "gitflic.ru/lms/backend/internal/application/usecases/attempt/commands"
	bankcommands "gitflic.ru/lms/backend/internal/application/usecases/bank/commands"
	bankqueries "gitflic.ru/lms/backend/internal/application/usecases/bank/queries"
	coursecommands "gitflic.ru/lms/backend/internal/application/usecases/course/commands"
	coursequeries "gitflic.ru/lms/backend/internal/application/usecases/course/queries"
	enrollmentcommands "gitflic.ru/lms/backend/internal/application/usecases/enrollment/commands"
	enrollmentqueries "gitflic.ru/lms/backend/internal/application/usecases/enrollment/queries"
	orgcommands "gitflic.ru/lms/backend/internal/application/usecases/organization/commands"
	orgqueries "gitflic.ru/lms/backend/internal/application/usecases/organization/queries"
	personcommands "gitflic.ru/lms/backend/internal/application/usecases/person/commands"
	personqueries "gitflic.ru/lms/backend/internal/application/usecases/person/queries"
	questioncommands "gitflic.ru/lms/backend/internal/application/usecases/question/commands"
	questiongrading "gitflic.ru/lms/backend/internal/application/usecases/question/grading"
	quizcommands "gitflic.ru/lms/backend/internal/application/usecases/quiz/commands"
	"gitflic.ru/lms/backend/internal/domain/account"
	attemptdomain "gitflic.ru/lms/backend/internal/domain/attempt"
	"gitflic.ru/lms/backend/internal/domain/grading/registry"
	questdomain "gitflic.ru/lms/backend/internal/domain/question"
	"gitflic.ru/lms/backend/internal/domain/role"
	"gitflic.ru/lms/backend/internal/transport/http/response"
	"github.com/google/uuid"
)

// API groups business HTTP handlers backed by application use cases.
type API struct {
	Accounts      AccountUseCases
	Organizations OrganizationUseCases
	Persons       PersonUseCases
	Banks         BankUseCases
	Questions     QuestionUseCases
	Quizzes       QuizUseCases
	Attempts      AttemptUseCases
	Courses       CourseUseCases
	Enrollments   EnrollmentUseCases
}

type AccountUseCases struct {
	Create         *accountcommands.CreateUseCase
	ChangeLogin    *accountcommands.ChangeLoginUseCase
	ChangePassword *accountcommands.ChangePasswordUseCase
	ChangeRole     *accountcommands.ChangeRoleUseCase
	Block          *accountcommands.BlockUseCase
	Activate       *accountcommands.ActivateUseCase
	Delete         *accountcommands.DeleteUseCase
	List           *accountqueries.ListQuery
	Get            *accountqueries.GetByIDQuery
}

type OrganizationUseCases struct {
	Create    *orgcommands.CreateUseCase
	Rename    *orgcommands.RenameUseCase
	ChangeINN *orgcommands.ChangeINNUseCase
	Delete    *orgcommands.DeleteByIDUseCase
	ListName  *orgqueries.ListByNameQuery
	ListINN   *orgqueries.ListByINNQuery
	Get       *orgqueries.GetDetailsByIDQuery
}

type PersonUseCases struct {
	Create            *personcommands.CreateUseCase
	CreateWithProfile *personcommands.CreateWithProfileUseCase
	Rename            *personcommands.RenameUseCase
	AttachProfile     *personcommands.AttachProfileUseCase
	ReplaceProfile    *personcommands.ReplaceProfileUseCase
	DetachProfile     *personcommands.DetachProfileUseCase
	ChangeSNILS       *personcommands.ChangeSNILSUseCase
	ChangeDOB         *personcommands.ChangeDateOfBirthUseCase
	ChangeJobTitle    *personcommands.ChangeJobTitleUseCase
	ChangeEducation   *personcommands.ChangeEducationUseCase
	AssignOrg         *personcommands.AssignOrganizationUseCase
	RemoveOrg         *personcommands.RemoveOrganizationUseCase
	Delete            *personcommands.DeleteByIDUseCase
	Get               *personqueries.GetDetailsByIdQuery
	ListLastName      *personqueries.ListByLastNameQuery
	ListOrg           *personqueries.ListByOrganizationIDQuery
	ListSNILS         *personqueries.ListBySnilsQuery
}

type BankUseCases struct {
	Create *bankcommands.CreateUseCase
	Rename *bankcommands.RenameUseCase
	Add    *bankcommands.AddQuestionsUseCase
	Remove *bankcommands.RemoveQuestionsUseCase
	Clear  *bankcommands.ClearQuestionsUseCase
	Delete *bankcommands.DeleteUseCase
	List   *bankqueries.ListQuery
	Get    *bankqueries.GetDetailsByIDQuery
}

type QuestionUseCases struct {
	Create         *questioncommands.CreateUseCase
	ChangeTitle    *questioncommands.ChangeTitleUseCase
	Selectable     *questioncommands.ChangeSelectableOptionsUseCase
	Sequence       *questioncommands.ChangeSequenceOptionsUseCase
	Matching       *questioncommands.ChangeMatchingPairsUseCase
	Short          *questioncommands.ChangeShortVariantsUseCase
	Grade          *questiongrading.GradeUseCase
	ValidateAnswer *questiongrading.ValidateAnswerUseCase
	Repository     interface {
		FindByID(rctx context.Context, id uuid.UUID) (questdomain.Question, error)
		DeleteByID(rctx context.Context, id uuid.UUID) error
	}
}

type QuizUseCases struct {
	Create        *quizcommands.CreateUseCase
	Rename        *quizcommands.RenameUseCase
	ChangeLimits  *quizcommands.ChangeLimitsUseCase
	ChangeShuffle *quizcommands.ChangeShufflePolicyUseCase
	Replace       *quizcommands.ReplaceSourcesUseCase
	Repository    quizports.Repository
}

type AttemptUseCases struct {
	Start         *attemptcommands.StartUseCase
	Answer        *attemptcommands.AddAnswerUseCase
	Finish        *attemptcommands.FinishUseCase
	Cancel        *attemptcommands.CancelUseCase
	Repository    attemptReader
	Query         attemptports.QueryService
	GradeRegistry *registry.Registry
}

type attemptReader interface {
	FindByID(context.Context, uuid.UUID) (*attemptdomain.Attempt, error)
}

type CourseUseCases struct {
	Create               *coursecommands.CreateCourseUseCase
	Rename               *coursecommands.RenameCourseUseCase
	AddBlock             *coursecommands.AddBlockToCourseUseCase
	RemoveBlock          *coursecommands.RemoveBlockFromCourseUseCase
	MoveBlock            *coursecommands.MoveCourseBlockUseCase
	AddElement           *coursecommands.AddElementToBlockUseCase
	RemoveElement        *coursecommands.RemoveElementFromBlockUseCase
	MoveElement          *coursecommands.MoveBlockElementUseCase
	ChangeCompletionMode *coursecommands.ChangeElementCompletionModeUseCase
	GetProgress          *coursecommands.GetProgressUseCase
	UploadContent        *coursecommands.UploadElementContentUseCase
	DownloadContent      *coursecommands.DownloadElementContentUseCase
	List                 *coursequeries.ListQuery
	Ratings              *coursequeries.RatingsQuery
	Statistics           *coursequeries.StudentStatisticsQuery
	Query                courseports.QueryService
}

type EnrollmentUseCases struct {
	Create *enrollmentcommands.CreateUseCase
	Get    *enrollmentqueries.GetByIDQuery
	List   *enrollmentqueries.ListQuery
}

func roleFromString(value string) (role.Role, error) {
	t, err := role.ParseType(value)
	if err != nil {
		return role.Role{}, err
	}
	return role.New(t)
}

func (api *API) ListAccounts(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorRole(r)
	if !ok {
		response.WriteError(w, r, response.NewError(http.StatusUnauthorized, "unauthorized", "session is required"))
		return
	}
	limit, offset := limitOffset(r)
	out, err := api.Accounts.List.Execute(r.Context(), accountqueries.ListInput{
		ActorRole: actor.role,
		Role:      role.Type(r.URL.Query().Get("role")),
		Status:    account.Status(r.URL.Query().Get("status")),
		Login:     r.URL.Query().Get("login"),
		Limit:     limit,
		Offset:    offset,
	})
	if err != nil {
		writeHandlerError(w, r, err)
		return
	}
	writeOK(w, r, out.Views, response.Links{"self": {Href: r.URL.RequestURI(), Method: http.MethodGet}})
}

func (api *API) CreateAccount(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorRole(r)
	if !ok {
		response.WriteError(w, r, response.NewError(http.StatusUnauthorized, "unauthorized", "session is required"))
		return
	}
	var req CreateAccountRequest
	if err := response.DecodeJSON(r, &req); err != nil {
		response.WriteError(w, r, response.NewError(http.StatusBadRequest, "invalid_json", "request body is invalid"))
		return
	}
	accountRole, err := roleFromString(req.Role)
	if err != nil {
		writeHandlerError(w, r, err)
		return
	}
	out, err := api.Accounts.Create.Execute(r.Context(), accountcommands.CreateInput{
		ActorRole: actor.role,
		Login:     req.Login,
		Password:  req.Password,
		Role:      accountRole,
		PersonID:  req.PersonID,
	})
	if err != nil {
		writeHandlerError(w, r, err)
		return
	}
	writeCreated(w, "/accounts", out.ID)
}

func (api *API) GetAccount(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorRole(r)
	if !ok {
		response.WriteError(w, r, response.NewError(http.StatusUnauthorized, "unauthorized", "session is required"))
		return
	}
	out, err := api.Accounts.Get.Execute(r.Context(), accountqueries.GetByIDInput{ActorRole: actor.role, AccountID: r.PathValue("id")})
	if err != nil {
		writeHandlerError(w, r, err)
		return
	}
	writeOK(w, r, out.View, response.Links{"self": {Href: r.URL.Path, Method: http.MethodGet}, "person": {Href: "/persons/" + out.View.PersonID, Method: http.MethodGet}})
}

func (api *API) ChangeAccountLogin(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorRole(r)
	if !ok {
		response.WriteError(w, r, response.NewError(http.StatusUnauthorized, "unauthorized", "session is required"))
		return
	}
	var req ChangeLoginRequest
	if err := response.DecodeJSON(r, &req); err != nil {
		response.WriteError(w, r, response.NewError(http.StatusBadRequest, "invalid_json", "request body is invalid"))
		return
	}
	out, err := api.Accounts.ChangeLogin.Execute(r.Context(), accountcommands.ChangeLoginInput{ActorRole: actor.role, AccountID: r.PathValue("id"), Login: req.Login})
	if err != nil {
		writeHandlerError(w, r, err)
		return
	}
	writeOK(w, r, map[string]string{"id": out.ID}, nil)
}

func (api *API) ChangeAccountPassword(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorRole(r)
	if !ok {
		response.WriteError(w, r, response.NewError(http.StatusUnauthorized, "unauthorized", "session is required"))
		return
	}
	var req ChangePasswordRequest
	if err := response.DecodeJSON(r, &req); err != nil {
		response.WriteError(w, r, response.NewError(http.StatusBadRequest, "invalid_json", "request body is invalid"))
		return
	}
	out, err := api.Accounts.ChangePassword.Execute(r.Context(), accountcommands.ChangePasswordInput{ActorRole: actor.role, AccountID: r.PathValue("id"), Password: req.Password})
	if err != nil {
		writeHandlerError(w, r, err)
		return
	}
	writeOK(w, r, map[string]string{"id": out.ID}, nil)
}

func (api *API) ChangeAccountRole(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorRole(r)
	if !ok {
		response.WriteError(w, r, response.NewError(http.StatusUnauthorized, "unauthorized", "session is required"))
		return
	}
	var req ChangeRoleRequest
	if err := response.DecodeJSON(r, &req); err != nil {
		response.WriteError(w, r, response.NewError(http.StatusBadRequest, "invalid_json", "request body is invalid"))
		return
	}
	accountRole, err := roleFromString(req.Role)
	if err != nil {
		writeHandlerError(w, r, err)
		return
	}
	out, err := api.Accounts.ChangeRole.Execute(r.Context(), accountcommands.ChangeRoleInput{ActorRole: actor.role, AccountID: r.PathValue("id"), AccountRole: accountRole})
	if err != nil {
		writeHandlerError(w, r, err)
		return
	}
	writeOK(w, r, map[string]string{"id": out.ID}, nil)
}

func (api *API) BlockAccount(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorRole(r)
	if !ok {
		response.WriteError(w, r, response.NewError(http.StatusUnauthorized, "unauthorized", "session is required"))
		return
	}
	out, err := api.Accounts.Block.Execute(r.Context(), accountcommands.BlockInput{ActorRole: actor.role, AccountID: r.PathValue("id")})
	if err != nil {
		writeHandlerError(w, r, err)
		return
	}
	writeOK(w, r, map[string]string{"id": out.ID}, nil)
}

func (api *API) ActivateAccount(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorRole(r)
	if !ok {
		response.WriteError(w, r, response.NewError(http.StatusUnauthorized, "unauthorized", "session is required"))
		return
	}
	out, err := api.Accounts.Activate.Execute(r.Context(), accountcommands.ActivateInput{ActorRole: actor.role, AccountID: r.PathValue("id")})
	if err != nil {
		writeHandlerError(w, r, err)
		return
	}
	writeOK(w, r, map[string]string{"id": out.ID}, nil)
}

func (api *API) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorRole(r)
	if !ok {
		response.WriteError(w, r, response.NewError(http.StatusUnauthorized, "unauthorized", "session is required"))
		return
	}
	if err := api.Accounts.Delete.Execute(r.Context(), accountcommands.DeleteInput{ActorRole: actor.role, AccountID: r.PathValue("id")}); err != nil {
		writeHandlerError(w, r, err)
		return
	}
	writeNoContent(w)
}
