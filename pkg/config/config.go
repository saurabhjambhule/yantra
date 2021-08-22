package config

import (
  "log"
  "os"
  "regexp"
  "strings"

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


func UpdatePlaceholder(input string) string {
  re, _ := regexp.Compile(`.*(\$\{.*\}).*`)

	for {
		matched := re.MatchString(input)
		if !matched {
			break;
		}

		matches := re.FindStringSubmatch(input)
    placeholder := strings.TrimSuffix(strings.TrimPrefix(matches[1],"${"), "}")
    value := getPlaceholderValue(placeholder)
		input = strings.Replace(input, matches[1], value, -1)
	}

  return input
}

// TODO: Add support for multiple placeholder locations
func getPlaceholderValue(input string) string {
  value := os.Getenv(input)

  return value
}
