package pubsub

import (
	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/container"
	"github.com/mylxsw/glacier/event"
	"github.com/mylxsw/glacier/infra"
)

type Provider struct{}

func (p Provider) Register(app container.Container) {
}

func (p Provider) Boot(app infra.Glacier) {
	app.MustResolve(func(em event.Manager) {
		em.Listen(func(ev SystemUpDownEvent) {
			log.With(ev).Debugf("system-up-down event received")
		})
	})
}
