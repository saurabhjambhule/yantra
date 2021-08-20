package config

type Volumes struct {
	Name string
}

type TaskDefinition struct {
	Family                  string
	Revision                int64
	ExecutionRoleArn        string
	ContainerDefinitions    []ContainerDefinition
  Cpu                     string
	Memory                  string
	NetworkMode             string
	Tags                    []KeyValuePair
  RequiresCompatibilities []string
	Volumes                 []Volumes
}

type ContainerDependency struct {
	Condition     string
	ContainerName string
}

type KeyValuePair struct {
	Name   string
	Value string
}

type HealthCheck struct {
	Command     string
	Interval    int64
  Retries     int64
  StartPeriod int64
  Timeout     int64
}

type LogConfiguration struct {
	LogDriver string
	Options   []KeyValuePair
}

type MountPoint struct {
	ContainerPath string
	ReadOnly      bool
  SourceVolume  string
}

type PortMappings struct {
	ContainerPort int64
	HostPort      int64
  Protocol      string
}

type ContainerDefinition struct {
	Name             string
	Essential        bool
  Command          []string
	Cpu              int64
  DependsOn        []ContainerDependency
  EntryPoint       string
	Environment      []KeyValuePair
  HealthCheck      HealthCheck
	Image            string
  LogConfiguration LogConfiguration
	Memory           int64
  MountPoints      []MountPoint
	PortMappings     []PortMappings
	StartTimeout     int64
  StopTimeout      int64
}

// type ECS struct {
// 	ClusterName    string
// 	ServiceName    string
// 	IamRoleArn     string
// 	TaskDefinition *ecs.RegisterTaskDefinitionInput
// 	TaskDefArn     *string
// 	TargetGroupArn *string
// }

func GetContainerDefinition(configPath string, configFile string) ContainerDefinition {
  conf := ContainerDefinition{}
	conf = fetchConfigFromFile(configPath, configFile, conf).(ContainerDefinition)

  return conf
}
