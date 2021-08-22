package config

import (
  "reflect"
	"testing"
)

func TestGetRoute53Config(t *testing.T) {
  testConfig := Route53 {
    Route53HostedZone: "Z2VPLB83J447ZF",
    DomainSuffix: "track",
    DomainPrefix: ".testing.siftery.com",
  }

	route53Config := GetRoute53Config("../../examples/config/ecs", "config")

  if reflect.DeepEqual(testConfig, route53Config) {
    t.Errorf("got %q, wanted %q", route53Config, testConfig)
  }
}
