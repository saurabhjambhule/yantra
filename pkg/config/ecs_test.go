package config

import (
	"testing"
)

func TestGetTaskDefinition(t *testing.T) {
  volume := Volume {
  	Name: "app",
  }
  _=volume

	tdConfig := GetTaskDefinition("../../examples/config/ecs", "task_defination")
  _=tdConfig
}

func TestGetECSConfig(t *testing.T) {
	ecsConfig := GetECSConfig("../../examples/config/ecs", "config")
  _=ecsConfig
}
