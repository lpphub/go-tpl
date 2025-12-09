package main

import (
	"fmt"
	"go-tpl/infra"
	"go-tpl/logic"
	"go-tpl/web"
)

func main() {
	// 1.初始化基础设施
	infra.Init()
	// 2.初始化逻辑层
	logic.Init()

	// 3.配置并启动web
	app := web.New()
	app.Run(fmt.Sprintf(":%d", infra.Cfg.Server.Port))
}
