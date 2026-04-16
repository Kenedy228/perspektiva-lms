package permission

import (
	"gitflic.ru/lms/internal/domain/permission/action"
	"gitflic.ru/lms/internal/domain/permission/resource"
)

type Params struct {
	Resource resource.Resource
	Actions  []action.Action
}
