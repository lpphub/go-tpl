package types

import "go-tpl/logic/shared"

type UserQueryReq struct {
	shared.Pagination
	Username string `json:"username"`
}

type RoleQueryReq struct {
	shared.Pagination
}
