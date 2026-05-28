package role

type Resource string

const (
	ResourceUser        Resource = "user"
	ResourceCourse      Resource = "course"
	ResourceBlock       Resource = "block"
	ResourceElement     Resource = "element"
	ResourceEnrollment  Resource = "enrollment"
	ResourceQuiz        Resource = "quiz"
	ResourceBank        Resource = "bank"
	ResourceAttempt     Resource = "attempt"
	ResourceSubmission  Resource = "submission"
	ResourceGrade       Resource = "grade"
	ResourceCertificate Resource = "certificate"
	ResourceFile        Resource = "file"
	ResourceAuditLog    Resource = "audit_log"
)
