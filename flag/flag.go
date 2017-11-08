package flag

import (
	"github.com/giantswarm/k8s-kvm-health/flag/service"
	"github.com/giantswarm/microkit/flag"
)

type Flag struct {
	Service service.Service
}

func New() *Flag {
	f := &Flag{}
	flag.Init(f)
	return f
}
