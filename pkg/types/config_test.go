package types

import (
	"flag"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
)

func TestMergeConfig(t *testing.T) {
	tmpFile, err := ioutil.TempFile("./", "*.yml")
	if !assert.Nil(t, err, "failed to create tmp config yml") {
		return
	}
	filename := tmpFile.Name()

	flagSet := flag.NewFlagSet("test", flag.PanicOnError)
	_ = flagSet.String("config-file", filename, "")

	testApp := cli.NewApp()
	testContext := cli.NewContext(testApp, flagSet, nil)

	if err := testContext.Set("config-file", filename); !assert.Nil(t, err, "failed to set temporary config file name as flag") {
		return
	}
	Config.Port = 8081
	testConfig := config{
		Port: 8080,
		DB: &db{
			HostIP: "1.2.3.4",
			Port:   1234,
		},
	}
	data, err := yaml.Marshal(testConfig)
	if !assert.Nil(t, err, "failed to unmarshal config to yaml format") {
		return
	}
	if _, err = tmpFile.Write(data); !assert.Nil(t, err, "failed to write data to temporary file") {
		return
	}

	if err := tmpFile.Close(); !assert.Nil(t, err, "failed to close temporary file") {
		return
	}
	defer os.Remove(filename)
	if err := MergeConfig(testContext); !assert.Nil(t, err, "failed to merge configuration file") {
		return
	}
	assert.Equal(t, 8081, Config.Port, "port configuration is not matched")
	assert.Equal(t, "1.2.3.4", Config.DB.HostIP, "database host IP configuration is not matched")
	assert.Equal(t, 1234, Config.DB.Port, "Database port configuration is not matched")
}
