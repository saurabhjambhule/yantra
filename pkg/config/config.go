package config

import (
  "log"

  "github.com/spf13/viper"
)

const configType string = "yml"

func fetchConfigFromFile(configPath string, configFile string, configuration interface{}) interface{} {
  viper.SetConfigName(configFile)
	viper.AddConfigPath(configPath)
	viper.AutomaticEnv()
	viper.SetConfigType(configType)

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
