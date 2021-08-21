package config

import (
  "log"

  "github.com/spf13/viper"
)

const (
	configType = "yml"
	envPrefix  = "yantra"
)


func fetchConfigFromFile(configPath string, configFile string, configuration interface{}) interface{} {
  viper.SetConfigName(configFile)
  viper.SetConfigType(configType)
	viper.AddConfigPath(configPath)
  viper.SetEnvPrefix(envPrefix)
	viper.AutomaticEnv()

	err := viper.ReadInConfig();
  if err != nil {
		log.Fatal(err)
	}

  // var configuration ContainerDefinition
	err = viper.Unmarshal(&configuration)
  if err != nil {
		log.Fatal(err)
	}

  return configuration
}
