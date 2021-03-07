package service

import (
	"github.com/mylxsw/container"
	"github.com/mylxsw/glacier/infra"
)

type Provider struct {
}

func (p Provider) Register(app container.Container) {
	app.MustSingletonOverride(NewJobService)
	app.MustSingletonOverride(NewNodeService)
	app.MustSingletonOverride(NewCredentialService)
}

func (p Provider) Boot(app infra.Glacier) {
}
