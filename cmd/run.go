package main

import (
	"go-tpl/web"
)

func main() {
	app := web.New()
	// 初始化资源
	app.Init()
	// 运行服务
	app.Run()
}
