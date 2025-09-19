package ext

import (
	"net/http"
	_ "net/http/pprof"

	"github.com/felixge/fgprof"
	"github.com/gin-gonic/gin"
)

func SetupPprof(engine *gin.Engine) {
	http.DefaultServeMux.Handle("/debug/fgprof", fgprof.Handler())

	go func() {
		if err := http.ListenAndServe(":6060", nil); err != nil {
			panic("go profiler server start error: " + err.Error())
		}
	}()
}
