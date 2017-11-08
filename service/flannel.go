package service

import (
	"fmt"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	"io/ioutil"
	"net"
	"os"
	"regexp"
	"strings"
	"time"
)

const (
	MaxRetry = 100
)

func (c *Config) LoadFlannelConfig() error {
	err := c.waitForFlannelFile(c.Logger)
	if err != nil {
		return microerror.Mask(err)
	}

	// fetch configuration from OS env
	confFile, err := c.readFlannelFile()
	if err != nil {
		return microerror.Mask(err)
	}

	// parse config and generate IP for interfaces
	err = c.parseIPs(confFile)
	if err != nil {
		return microerror.Mask(err)
	}
	// debug output
	c.Logger.Log("debug", fmt.Sprintf("Loaded Config: %+v", c.Flag.Service))
	return nil
}

// readFlannelFile fetches ENVs values and read kvm file
func (c *Config) readFlannelFile() ([]byte, error) {
	fileContent, err := ioutil.ReadFile(c.Flag.Service.FlannelFile)
	if err != nil {
		return nil, microerror.Maskf(invalidFlannelFileError, "%s", c.Flag.Service.FlannelFile)
	}

	return fileContent, nil
}

// parseIPs parses kvm configuration file and generate ips for interface
func (c *Config) parseIPs(confFile []byte) error {
	// get FLANNEL_SUBNET from kvm file via regexp
	r, _ := regexp.Compile("FLANNEL_SUBNET=[0-9]+.[0-9]+.[0-9]+.[0-9]+/[0-9]+")
	flannelLine := r.Find(confFile)
	// check if regexp returned non-empty line
	if len(flannelLine) < 5 {
		return microerror.Mask(invalidKVMConfigurationError)
	}

	// parse kvm subnet
	flannelSubnetStr := strings.Split(string(flannelLine), "=")[1]
	flannelIP, _, err := net.ParseCIDR(flannelSubnetStr)
	if err != nil {
		return microerror.Maskf(failedParsingFlannelSubnetError, "%v", err)
	}
	// force ipv4 for later trick
	flannelIP = flannelIP.To4()

	// get kvm ip, which is just one number bigger than bridge hence the [3]++ trick
	flannelIP[3]++
	c.Flag.Service.IPAddress = flannelIP.String()

	return nil
}

// waitForFlannelFile waits until flannel file is created
func (c *Config) waitForFlannelFile(newLogger micrologger.Logger) error {
	// wait for file creation
	for count := 0; ; count++ {
		// don't wait forever, if file is not created within retry limit, exit with failure
		if count > MaxRetry {
			return microerror.Maskf(invalidFlannelFileError, "After 100sec flannel file is not created. Exiting")
		}
		// check if file exists
		if _, err := os.Stat(c.Flag.Service.FlannelFile); !os.IsNotExist(err) {
			break
		}
		newLogger.Log("debug", fmt.Sprintf("Waiting for file '%s' to be created.", c.Flag.Service.FlannelFile))
		time.Sleep(1 * time.Second)
	}
	// all good
	return nil
}
