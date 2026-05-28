package handlers

// LoginRequest is the request body for POST /auth/login.
type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// SessionResponse is returned after a successful authentication.
type SessionResponse struct {
	Token     string     `json:"token"`
	TokenType string     `json:"token_type"`
	ExpiresAt int64      `json:"expires_at"`
	Account   AccountRef `json:"account"`
	Person    PersonRef  `json:"person"`
	Role      string     `json:"role"`
}

// AccountRef references an account resource.
type AccountRef struct {
	ID string `json:"id"`
}

// PersonRef references a person resource.
type PersonRef struct {
	ID string `json:"id"`
}

// CreateAccountRequest is the intended DTO for account creation.
type CreateAccountRequest struct {
	PersonID string `json:"person_id"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type ChangeLoginRequest struct {
	Login string `json:"login"`
}

type ChangePasswordRequest struct {
	Password string `json:"password"`
}

type ChangeRoleRequest struct {
	Role string `json:"role"`
}

// OrganizationRequest is used by organization create/update endpoints.
type OrganizationRequest struct {
	Name    string `json:"name"`
	INN     string `json:"inn,omitempty"`
	INNType string `json:"inn_type,omitempty"`
}

// PersonRequest is used by person create/update endpoints.
type PersonRequest struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	MiddleName string `json:"middle_name,omitempty"`
}

// ProfileRequest is used by person profile endpoints.
type ProfileRequest struct {
	SNILS          string `json:"snils"`
	DateOfBirth    string `json:"date_of_birth"`
	JobTitle       string `json:"job_title,omitempty"`
	Education      string `json:"education,omitempty"`
	OrganizationID string `json:"organization_id,omitempty"`
}

// BankRequest is used by question bank create/update endpoints.
type BankRequest struct {
	Title string `json:"title"`
}

// QuestionIDsRequest is used to add or remove question references.
type QuestionIDsRequest struct {
	QuestionIDs []string `json:"question_ids"`
}

// QuestionRequest is the intended DTO for question create/update endpoints.
type QuestionRequest struct {
	Type              string                  `json:"type"`
	Title             string                  `json:"title"`
	Attachment        *AttachmentInput        `json:"attachment,omitempty"`
	SelectableOptions []SelectableOptionInput `json:"selectable_options,omitempty"`
	SequenceOptions   []SequenceOptionInput   `json:"sequence_options,omitempty"`
	MatchingPairs     []MatchingPairInput     `json:"matching_pairs,omitempty"`
	ShortVariants     []ShortVariantInput     `json:"short_variants,omitempty"`
}

// AttachmentInput references an already uploaded media object.
type AttachmentInput struct {
	MediaType string `json:"media_type"`
	FileName  string `json:"file_name"`
	SizeBytes int64  `json:"size_bytes"`
}

type SelectableOptionInput struct {
	Text      string `json:"text"`
	IsCorrect bool   `json:"is_correct"`
}

type SequenceOptionInput struct {
	Text string `json:"text"`
}

type MatchingPairInput struct {
	Prompt string `json:"prompt"`
	Match  string `json:"match"`
}

type ShortVariantInput struct {
	Text string `json:"text"`
}

// QuizRequest is used by quiz create/update endpoints.
type QuizRequest struct {
	Title            string       `json:"title"`
	MaxAttempts      int          `json:"max_attempts"`
	TimeLimitSeconds int          `json:"time_limit_seconds"`
	ShuffleQuestions bool         `json:"shuffle_questions"`
	Sources          []QuizSource `json:"sources"`
}

// QuizSource describes a source bank and selection criteria.
type QuizSource struct {
	BankID        string   `json:"bank_id"`
	CriteriaType  string   `json:"criteria_type"`
	QuestionCount int      `json:"question_count,omitempty"`
	QuestionIDs   []string `json:"question_ids,omitempty"`
}

// CourseRequest is used by course create/update endpoints.
type CourseRequest struct {
	Title string `json:"title"`
}

type CourseBlockRequest struct {
	Title string `json:"title"`
}

type CourseElementRequest struct {
	Title          string `json:"title"`
	Type           string `json:"type"`
	FileName       string `json:"file_name,omitempty"`
	SizeBytes      int64  `json:"size_bytes,omitempty"`
	QuizID         string `json:"quiz_id,omitempty"`
	CompletionMode string `json:"completion_mode,omitempty"`
}

// EnrollmentRequest creates an enrollment.
type EnrollmentRequest struct {
	AccountID     string `json:"account_id"`
	CourseID      string `json:"course_id"`
	ActivatedAt   string `json:"activated_at"`
	DeactivatedAt string `json:"deactivated_at"`
}

// AttemptStartRequest creates a quiz attempt.
type AttemptStartRequest struct {
	AccountID    string `json:"account_id"`
	EnrollmentID string `json:"enrollment_id"`
	QuizID       string `json:"quiz_id"`
}

type AttemptAnswerRequest struct {
	Type          string            `json:"type"`
	OptionIDs     []string          `json:"option_ids,omitempty"`
	MatchingPairs map[string]string `json:"matching_pairs,omitempty"`
	ShortInput    string            `json:"short_input,omitempty"`
}

type MoveBlockRequest struct {
	From int `json:"from"`
	To   int `json:"to"`
}

type MoveElementRequest struct {
	From int `json:"from"`
	To   int `json:"to"`
}

type ChangeCompletionModeRequest struct {
	CompletionMode string `json:"completion_mode"`
}

type ProgressResponse struct {
	EnrollmentID        string   `json:"enrollment_id"`
	CompletedCount      int      `json:"completed_count"`
	Percent             int      `json:"percent"`
	TotalTrackedItems   int      `json:"total_tracked_items"`
	CompletedElementIDs []string `json:"completed_element_ids"`
}
