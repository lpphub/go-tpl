package rest

import (
	"errors"
	"go-tpl/infra/jwt"
	"go-tpl/infra/logging"
	"go-tpl/logic"
	"go-tpl/logic/shared"
	"go-tpl/web/base"
	"go-tpl/web/types"

	"github.com/gin-gonic/gin"
)

func Test(c *gin.Context) {
	logging.Info(c, "test1")

	logging.Info(c.Request.Context(), "test2")

	logging.Info(c, "test3")

	logging.Errorw(c, errors.New("test"))

	base.OKWithData(c, "ok")
}

func Register(c *gin.Context) {
	var req types.RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		base.FailWithError(c, shared.ErrInvalidParam)
		return
	}

	// 创建用户
	user, err := logic.Svc.User.Create(c, types.CreateUserReq{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		base.FailWithError(c, err)
		return
	}

	// 生成 token
	token, err := jwt.GenerateToken(user.ID)
	if err != nil {
		logging.Errorw(c, err)
		base.Fail(c, 500, "生成token失败")
		return
	}

	base.OKWithData(c, gin.H{
		"token": token,
		"user":  user,
	})
}

func Login(c *gin.Context) {
	var req types.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		base.FailWithError(c, shared.ErrInvalidParam)
		return
	}

	// 验证用户名和密码
	user, err := logic.Svc.User.ValidateLogin(c, req.Username, req.Password)
	if err != nil {
		base.FailWithError(c, err)
		return
	}

	// 生成 token
	token, err := jwt.GenerateToken(user.ID)
	if err != nil {
		logging.Errorw(c, err)
		base.Fail(c, 500, "生成token失败")
		return
	}

	base.OKWithData(c, gin.H{
		"token": token,
		"user":  user,
	})
}
