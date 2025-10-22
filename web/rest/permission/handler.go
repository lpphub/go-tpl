package permission

import (
	"go-tpl/web/base"

	"github.com/gin-gonic/gin"
)

func List(c *gin.Context) {
	// todo 列表
	base.OKWithData(c, nil)
}

// GetByID 根据ID获取用户
func GetByID(c *gin.Context) {
	// todo 获取
	base.OKWithData(c, nil)
}
