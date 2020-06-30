package healthz

import (
	"github.com/giantswarm/microendpoint/service/healthz"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"

	"github.com/giantswarm/k8s-kvm-health/service/healthz/kvm"
)

// Config represents the configuration used to create a healthz service.
type Config struct {
	// Dependencies.
	CheckAPI  bool
	IPAddress string
	Logger    micrologger.Logger
}

// New creates a new configured healthz service.
func New(config Config) (*Service, error) {
	var err error

	var kvmService healthz.Service
	{
		kvmServiceConfig := kvm.Config{
			CheckAPI: config.CheckAPI,
			IP:       config.IPAddress,
			Logger:   config.Logger,
		}

		kvmService, err = kvm.New(kvmServiceConfig)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	newService := &Service{
		KVM: kvmService,
	}

	return newService, nil
}

// Service is the healthz service collection.
type Service struct {
	KVM healthz.Service
}
