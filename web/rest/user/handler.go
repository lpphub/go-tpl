package user

import (
	"go-tpl/logic"
	"go-tpl/web/base"
	"go-tpl/web/types"

	"github.com/gin-gonic/gin"
)

func List(c *gin.Context) {
	var req types.UserQueryReq
	if err := c.ShouldBind(&req); err != nil {
		base.FailWithErr(c, err)
		return
	}

	data, err := logic.UserSvc.List(c, req)
	if err != nil {
		base.FailWithErr(c, err)
		return
	}
	base.OKWithData(c, data)
}

// GetByID 根据ID获取用户
func GetByID(c *gin.Context) {
	// todo 获取用户
	base.OKWithData(c, nil)
}
