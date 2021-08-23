package plume

import (
	"github.com/iakud/plume/service"
)

type options struct {
	services []service.Service
}

type Option func(*options)

func WithServices(services ...service.Service) Option {
	return func(opt *options) {
		opt.services = services
	}
}