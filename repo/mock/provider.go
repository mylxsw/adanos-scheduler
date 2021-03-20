package mock

import (
	"github.com/mylxsw/glacier/infra"
)

type Provider struct {
}

func (p Provider) Register(app infra.Binder) {
	app.MustSingletonOverride(NewNodeRepo)
	app.MustSingletonOverride(NewCredentialRepo)
	app.MustSingletonOverride(NewNodeCredentialRepo)
	app.MustSingletonOverride(NewJobRepo)
	app.MustSingletonOverride(NewJobPlanRepo)
	app.MustSingletonOverride(NewJobHistoryRepo)
}

func (p Provider) Boot(app infra.Resolver) {
}
