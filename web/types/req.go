package types

import "go-tpl/logic/base"

type UserQueryReq struct {
	base.Pagination
	Username string `json:"username"`
}

type RoleQueryReq struct {
	base.Pagination
}
