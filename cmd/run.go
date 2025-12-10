package cmd

import (
	"go-tpl/web"
)

func Serve() {
	app := web.New()
	app.Run()
}
