package config

import (
	"bytes"
	"flag"
	"os"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	yaml "gopkg.in/yaml.v2"
)

func viperDecoderConfigOptions(config *mapstructure.DecoderConfig) {
	// We have some lists predefined and want them to be overriden when they are set in the config. Else we keep them.
	config.ZeroFields = true
}

func LoadConfiguration() (*Configuration, error) {
	envPrefix := "MEDIANEXUS"

	configFile := os.Getenv(envPrefix + "_CONFIG")
	if len(configFile) < 1 {
		flag.Parse()
		if flag.NArg() > 0 {
			configFile = flag.Arg(0)
		}
	}

	return parseConfiguration(configFile, envPrefix)
}

func parseConfiguration(configFilePath string, envPrefix string) (*Configuration, error) {
	if err := initializeViper(envPrefix); err != nil {
		return nil, err
	}

	if len(configFilePath) > 0 {
		if err := readConfigFile(configFilePath); err != nil {
			return nil, err
		}
	}

	configuration := NewConfiguration()
	if err := viper.Unmarshal(&configuration, viperDecoderConfigOptions); err != nil {
		return nil, errors.Wrap(err, "unable to deserialize config file")
	}

	if err := configuration.Validate(); err != nil {
		return nil, errors.Wrap(err, "configuration invalid")
	}

	return &configuration, nil
}

func initializeViper(envPrefix string) error {
	if err := registerConfigurationSchema(); err != nil {
		return err
	}

	viper.SetEnvPrefix(envPrefix)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	return nil
}

func registerConfigurationSchema() error {
	configuration := NewConfiguration()

	output, err := yaml.Marshal(configuration)
	if err != nil {
		return errors.Wrap(err, "failed to generate configuration schema")
	}

	viper.SetConfigType("yaml")
	err = viper.ReadConfig(bytes.NewBuffer(output))
	if err != nil {
		return errors.Wrap(err, "failed to read config schema")
	}

	return nil
}

func readConfigFile(configFilePath string) error {
	viper.SetConfigFile(configFilePath)

	// why merge and not read? We already read one, the configuration schema. So we need to merge this new config
	// with the already existing one.
	if err := viper.MergeInConfig(); err != nil {
		return errors.Wrap(err, "error reading config file")
	}

	return nil
}
