package permission

import (
	"errors"
)

var (
	ErrPermissionExists  = errors.New("permission on resource already exists")
	ErrDuplicatedActions = errors.New("provided actions contains duplicate records")
	ErrNoResourceFound   = errors.New("permission set doesn't contain provided resource")
	ErrInvalidResource   = errors.New("resource is invalid")
	ErrInvalidAction     = errors.New("action is invalid")
)
