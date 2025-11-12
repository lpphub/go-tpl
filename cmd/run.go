package main

import (
	"context"
	"errors"
	"go-tpl/infra"
	"go-tpl/infra/monitor"
	"go-tpl/logic"
	"go-tpl/web"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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

	run(app)
}

func run(handler http.Handler) {
	srv := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}
	go func() {
		log.Printf("Server starting on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	} else {
		log.Println("Server shutdown completed")
	}
}
