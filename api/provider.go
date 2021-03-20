package api

import (
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gorilla/mux"
	"github.com/mylxsw/adanos-scheduler/api/controller"
	"github.com/mylxsw/adanos-scheduler/repo"
	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/glacier/infra"
	"github.com/mylxsw/glacier/listener"
	"github.com/mylxsw/glacier/web"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Provider struct{}

func (p Provider) Aggregates() []infra.Provider {
	return []infra.Provider{
		web.Provider(
			listener.FlagContext("listen"),
			web.SetIgnoreLastSlashOption(true),
			web.SetExceptionHandlerOption(p.exceptionHandler),
			web.SetRouteHandlerOption(p.routes),
			web.SetMuxRouteHandlerOption(p.muxRoutes),
		),
	}
}

func (p Provider) muxRoutes(router *mux.Router) {
	// prometheus metrics
	router.PathPrefix("/metrics").Handler(promhttp.Handler())
	// health check
	router.PathPrefix("/health").Handler(HealthCheck{})
}

func (p Provider) routes(cc infra.Resolver, router web.Router, mw web.RequestMiddleware) {
	mws := make([]web.HandlerDecorator, 0)
	mws = append(mws, mw.AccessLog(log.Module("api")), mw.CORS("*"))

	router.WithMiddleware(mws...).Controllers(
		"/api",
		controller.NewJobController(cc),
		controller.NewNodeController(cc),
		controller.NewCredentialController(cc),
	)
}

func (p Provider) exceptionHandler(ctx web.Context, err interface{}) web.Response {
	if errTyped, ok := err.(error); ok {
		if errors.Is(errTyped, repo.ErrNotFound) {
			return ctx.JSONError(fmt.Sprintf("%v", err), http.StatusNotFound)
		}
	}

	log.With(string(debug.Stack())).Errorf("request handle error: %v", err)
	return ctx.JSONError(fmt.Sprintf("%v", err), http.StatusInternalServerError)
}

func (p Provider) Register(app infra.Binder) {
}

func (p Provider) Boot(app infra.Resolver) {
}

