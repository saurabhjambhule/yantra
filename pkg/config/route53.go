package config

type Route53 struct {
    Route53HostedZone string
    DomainSuffix      string
    DomainPrefix      string
}

func GetRoute53Config(configPath string, configFile string) Route53 {
  conf := Route53{}
	conf = fetchConfigFromFile(configPath, configFile, conf).(Route53)

  return conf
}
