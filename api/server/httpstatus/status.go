package httpstatus

import "errors"

var (
	ParamsError        = errors.New("参数错误")
	OperateFailed      = errors.New("操作失败")
	NoPermission       = errors.New("无权限")
	InnerError         = errors.New("inner error")
	NoUserIdError      = errors.New("请登录")
	RoleExistError     = errors.New("角色已存在")
	RoleNotExistError  = errors.New("角色不存在")
	MenusExistError    = errors.New("权限已存在")
	MenusNtoExistError = errors.New("权限不存在")
)
