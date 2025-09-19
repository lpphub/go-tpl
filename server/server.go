package server

import (
	"fmt"
	"go-tpl/ext"
	"go-tpl/server/api"
	"go-tpl/server/common/errs"
	"go-tpl/server/infra/global"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lpphub/golib/logger"
	"github.com/lpphub/golib/logger/logx"
	"github.com/lpphub/golib/web"
)

func Serve() {
	global.InitResource()
	defer global.Clear()

	app := gin.New()

	web.Bootstraps(app, web.BootstrapConf{
		Cors: true,
		AccessLog: web.AccessLogConfig{
			Enable:    true,
			SkipPaths: []string{"/metrics"},
		},
		CustomRecovery: func(ctx *gin.Context, err any) {
			logx.Error(ctx, fmt.Sprintf("server error: %+v", err))
			web.JsonWithError(ctx, errs.ErrServerError)
		},
	})

	ext.SetupPprof(app)
	ext.SetupMetrics(app)

	api.SetupRoute(app)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: app,
	}
	logger.Log().Info().Msgf("Listening and serving HTTP on %s", srv.Addr)
	web.ListenAndServe(srv)
}
