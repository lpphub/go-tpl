package web

import (
	"context"
	"errors"
	"go-tpl/infra/logger/logx"
	"go-tpl/infra/monitor"
	"go-tpl/web/rest"
	"go-tpl/web/rest/permission"
	"go-tpl/web/rest/role"
	"go-tpl/web/rest/user"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

type App struct {
	Engine *gin.Engine
}

func New() *App {
	app := &App{
		Engine: gin.Default(),
	}

	app.setupRouter()
	return app
}

func (a *App) setupRouter() {
	r := a.Engine

	// pprof and metrics
	//monitor.StartPprof()
	monitor.RegisterMetrics(r)

	api := r.Group("/api")
	// 公共中间件
	api.Use(logx.GinLogMiddleware())

	// 公共路由
	api.GET("/test", rest.Test)
	api.POST("/register", rest.Register)
	api.POST("/login", rest.Login)
	api.POST("/refresh", rest.RefreshToken)

	// 注册接口处理
	user.Register(api)
	role.Register(api)
	permission.Register(api)
}

func (a *App) Run(addr string) {
	srv := &http.Server{
		Addr:    addr,
		Handler: a.Engine,
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
