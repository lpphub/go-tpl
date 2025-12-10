package cmd

import (
	"go-tpl/web"
)

func Serve() {
	app := web.New()
	// 初始化资源
	app.Init()
	// 运行服务
	app.Run()
}
