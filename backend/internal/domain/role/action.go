package role

type Action string

const (
	ActionRead    Action = "read"
	ActionWrite   Action = "write"
	ActionCreate  Action = "create"
	ActionUpdate  Action = "update"
	ActionDelete  Action = "delete"
	ActionPublish Action = "publish"
	ActionEnroll  Action = "enroll"
	ActionSubmit  Action = "submit"
	ActionGrade   Action = "grade"
	ActionManage  Action = "manage"
)
