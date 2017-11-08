package service

import (
	"github.com/giantswarm/k8s-kvm-health/flag"
	"github.com/giantswarm/microerror"
	"strings"
	"testing"
)

func Test_Flannel_ParseIP(t *testing.T) {
	tests := []struct {
		config             func(flannelFile []byte) (string, error)
		flannelFileContent []byte
		expectedIP         string
		expectedErr        error
	}{
		// test 0
		{
			config: func(flannelFile []byte) (string, error) {
				conf := DefaultConfig()
				conf.Flag = flag.New()
				err := conf.parseIPs(flannelFile)
				return conf.Flag.Service.IPAddress, err
			},
			expectedIP: "172.23.3.66",
			flannelFileContent: []byte(`FLANNEL_NETWORK=172.23.3.0/24
FLANNEL_SUBNET=172.23.3.65/30
FLANNEL_MTU=1450
FLANNEL_IPMASQ=false`),
			expectedErr: nil,
		},
		// test 1
		{
			config: func(flannelFile []byte) (string, error) {
				conf := DefaultConfig()
				conf.Flag = flag.New()
				err := conf.parseIPs(flannelFile)
				return conf.Flag.Service.IPAddress, err
			},
			expectedIP: "198.168.0.2",
			flannelFileContent: []byte(`FLANNEL_NETWORK=198.168.0.0/24
FLANNEL_SUBNET=198.168.0.1/30
FLANNEL_MTU=1450
FLANNEL_IPMASQ=false`),
			expectedErr: nil,
		},
		// test 2 - missing FLANNEL_SUBNET
		{
			config: func(flannelFile []byte) (string, error) {
				conf := DefaultConfig()
				conf.Flag = flag.New()
				err := conf.parseIPs(flannelFile)
				return conf.Flag.Service.IPAddress, err
			},
			expectedIP: "",
			flannelFileContent: []byte(`FLANNEL_NETWORK=192.168.0.0/24
FLANNEL_MTU=1450
FLANNEL_IPMASQ=false`),
			expectedErr: invalidKVMConfigurationError,
		},
		// test 3 - invalid subnet in kvm file
		{
			config: func(flannelFile []byte) (string, error) {
				conf := DefaultConfig()
				conf.Flag = flag.New()
				err := conf.parseIPs(flannelFile)
				return conf.Flag.Service.IPAddress, err
			},
			expectedIP: "",
			flannelFileContent: []byte(`FLANNEL_NETWORK=198.168.0.0/24
FLANNEL_SUBNET=_x.68.c.0/30
FLANNEL_MTU=1450
FLANNEL_IPMASQ=false`),
			expectedErr: invalidKVMConfigurationError,
		},
		// test 4 - empty kvm file
		{
			config: func(flannelFile []byte) (string, error) {
				conf := DefaultConfig()
				conf.Flag = flag.New()
				err := conf.parseIPs(flannelFile)
				return conf.Flag.Service.IPAddress, err
			},
			expectedIP:         "",
			flannelFileContent: []byte(``),
			expectedErr:        invalidKVMConfigurationError,
		},
		// test 5 - non kvm file
		{
			config: func(flannelFile []byte) (string, error) {
				conf := DefaultConfig()
				conf.Flag = flag.New()
				err := conf.parseIPs(flannelFile)
				return conf.Flag.Service.IPAddress, err
			},
			expectedIP: "",
			flannelFileContent: []byte(`machine:
  services:
    - docker

dependencies:
  override:
    - |
      wget -q $(curl -sS -H "Authorization: token $RELEASE_TOKEN" https://api.github.com/repos/giantswarm/architect/releases/latest | grep browser_download_url | head -n 1 | cut -d '"' -f 4)
    - chmod +x ./architect
    - ./architect version
`),
			expectedErr: invalidKVMConfigurationError,
		},
	}

	for index, test := range tests {
		ip, err := test.config(test.flannelFileContent)

		if microerror.Cause(err) != microerror.Cause(test.expectedErr) {
			t.Fatalf("%d: unexcepted error, expected %#v but got %#v", index, test.expectedErr, err)
		}
		if test.expectedErr == nil {
			if strings.Compare(ip, test.expectedIP) != 0 {
				t.Fatalf("%d: Incorrent ip, expected %s but got %s.", index, test.expectedIP, ip)
			}
		}
	}
}
