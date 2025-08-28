package vo

type UserVO struct {
	Uid      int64  `json:"uid" binding:"required"`
	Nickname string `json:"nickname"`
}
