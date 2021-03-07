package api

import (
	"github.com/mylxsw/adanos-scheduler/api/controller"
	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/container"
	"github.com/mylxsw/glacier/web"
)

func routers(cc container.Container) func(router *web.Router, mw web.RequestMiddleware) {
	return func(router *web.Router, mw web.RequestMiddleware) {
		mws := make([]web.HandlerDecorator, 0)
		mws = append(mws, mw.AccessLog(log.Module("api")), mw.CORS("*"))

		router.WithMiddleware(mws...).Controllers(
			"/api",
			controller.NewJobController(cc),
			controller.NewNodeController(cc),
			controller.NewCredentialController(cc),
		)
	}
}
