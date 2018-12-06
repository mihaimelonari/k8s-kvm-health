package service

import (
	"github.com/giantswarm/k8s-kvm-health/flag"
	"github.com/giantswarm/k8s-kvm-health/service/healthz"
	"github.com/giantswarm/microendpoint/service/version"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"strings"
	"sync"
)

// Config represents the configuration used to create a new service.
type Config struct {
	// Dependencies.
	Logger micrologger.Logger

	// Settings.
	Flag *flag.Flag

	Description string
	GitCommit   string
	Name        string
	Source      string
}

// DefaultConfig provides a default configuration to create a new service by
// best effort.
func DefaultConfig() Config {
	return Config{
		// Dependencies.
		Logger: nil,

		// Settings.
		Flag: nil,

		Description: "",
		GitCommit:   "",
		Name:        "",
		Source:      "",
	}
}

// New creates a new configured service object.
func New(config Config) (*Service, error) {
	// Settings.
	if config.Flag == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Flag must not be empty")
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "config.Logger must not be empty")
	}

	var err error

	// load kvm network configuration
	err = config.LoadFlannelConfig()
	if err != nil {
		return nil, microerror.Mask(err)
	}

	var healthzService *healthz.Service
	{
		healthzConfig := healthz.Config{
			CheckAPI:  false,
			IPAddress: config.Flag.Service.IPAddress,
			Logger:    config.Logger,
		}

		if config.Flag.Service.CheckAPI == strings.ToLower("true") {
			healthzConfig.CheckAPI = true
		}

		healthzService, err = healthz.New(healthzConfig)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	var versionService *version.Service
	{
		versionConfig := version.DefaultConfig()

		versionConfig.Description = config.Description
		versionConfig.GitCommit = config.GitCommit
		versionConfig.Name = config.Name
		versionConfig.Source = config.Source

		versionService, err = version.New(versionConfig)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	newService := &Service{
		// Dependencies.
		Healthz: healthzService,
		Version: versionService,

		// Internals
		bootOnce: sync.Once{},
	}

	return newService, nil
}

type Service struct {
	// Dependencies.
	Healthz *healthz.Service
	Version *version.Service

	// Internals.
	bootOnce sync.Once
}
