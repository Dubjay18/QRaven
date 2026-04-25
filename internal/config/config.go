package config

import (
	"log"
	"path/filepath"
	"qraven/utils"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

// Setup initialize configuration
var (
	// Params ParamsConfiguration
	Config *Configuration
)

// Params = getConfig.Params
func Setup(logger *utils.Logger, name string) *Configuration {
	var baseConfiguration *BaseConfig
	v := viper.New()

	v.SetConfigType("env")
	v.AutomaticEnv()

	if name != "" {
		fileName := filepath.Base(name)
		fileExt := filepath.Ext(fileName)

		if fileExt != "" {
			v.SetConfigFile(name)
		} else {
			v.SetConfigName(fileName)
			v.AddConfigPath(".")
			dirName := filepath.Dir(name)
			if dirName != "." {
				v.AddConfigPath(dirName)
			}
		}
	}

	if err := v.ReadInConfig(); err != nil {
		log.Printf("Error reading config file, %s", err)
		log.Printf("Reading from environment variable")
	}

	var config BaseConfig
	err := BindKeys(v, config)
	if err != nil {
		log.Fatalf("Unable to bindkeys in struct, %v", err)
	}

	err = v.Unmarshal(&baseConfiguration)
	if err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	configuration := baseConfiguration.SetupConfigurationn()

	// Params = configuration.Params
	Config = configuration
	logger.Info("configurations loading successfully")
	return configuration
}

// GetConfig helps you to get configuration data
func GetConfig() *Configuration {
	return Config
}

func BindKeys(v *viper.Viper, input interface{}) error {

	envKeysMap := &map[string]interface{}{}
	if err := mapstructure.Decode(input, &envKeysMap); err != nil {
		return err
	}
	for k := range *envKeysMap {
		if bindErr := v.BindEnv(k); bindErr != nil {
			return bindErr
		}
	}

	return nil
}
