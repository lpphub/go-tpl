package role

import (
	"go-tpl/logic"
	"go-tpl/web/base"
	"go-tpl/web/types"

	"github.com/gin-gonic/gin"
)

func List(c *gin.Context) {
	var req types.RoleQueryReq
	if err := c.ShouldBind(&req); err != nil {
		base.FailWithErr(c, err)
		return
	}

	data, err := logic.RoleSvc.List(c, req.Pagination)
	if err != nil {
		base.FailWithErr(c, err)
		return
	}
	base.OKWithData(c, data)
}

// GetByID 根据ID获取
func GetByID(c *gin.Context) {
	// todo 获取
	base.OKWithData(c, nil)
}
