package main

import (
	"go-tpl/infra"
	"go-tpl/infra/monitor"
	"go-tpl/logic"
	"go-tpl/web"
	"log"
)

func main() {
	// 1.初始化基础设施
	infra.Init()
	// 2.初始化逻辑层
	logic.Init()

	// 3.启动服务
	app := web.SetupRouter()

	// 4.启动监控服务
	monitor.SetupMetrics(app)
	//monitor.SetupPprof()

	err := app.Run(":8080")
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
