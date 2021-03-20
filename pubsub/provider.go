package pubsub

import (
	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/glacier/event"
	"github.com/mylxsw/glacier/infra"
)

type Provider struct{}

func (p Provider) Aggregates() []infra.Provider {
	return []infra.Provider{
		event.Provider(p.listeners),
	}
}

func (p Provider) Register(app infra.Binder) {
}

func (p Provider) Boot(app infra.Resolver) {
}

func (p Provider) listeners(cc infra.Resolver, listener event.Listener) {
	listener.Listen(func(ev SystemUpDownEvent) {
		log.With(ev).Debugf("system-up-down event received")
	})
}
