package config

type Volume struct {
	Name string
}

type ContainerDependency struct {
	Condition     string
	ContainerName string
}

type KeyValuePair struct {
	Name   string
	Value  string
}

type HealthCheck struct {
	Command     []string
	Interval    int64
  Retries     int64
  StartPeriod int64
  Timeout     int64
}

type LogConfiguration struct {
	LogDriver string
	Options   map[string]*string
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
	Name              string
	Essential         bool
  Command           []string
	Cpu               int64
  DependsOn         []ContainerDependency
  EntryPoint        []string
	Environment       []KeyValuePair
  HealthCheck       HealthCheck
	Image             string
  LogConfiguration  LogConfiguration
	Memory            int64
	MemoryReservation int64
  MountPoints       []MountPoint
	PortMappings      []PortMappings
	StartTimeout      int64
  StopTimeout       int64
}

type TaskDefinition struct {
  ContainerDefinitions 		  []ContainerDefinition
  Cpu 											string
	Memory 										string
  ExecutionRoleArn 					string
	TaskRoleArn 							string
  Family 										string
  NetworkMode 							string
  RequiresCompatibilities   []string
  Tags											[]KeyValuePair
  Volumes 									[]Volume
}

type NetworkConfiguration struct {
    AwsvpcConfiguration AwsVpcConfiguration
}

type AwsVpcConfiguration struct {
    AssignPublicIp string
    SecurityGroups []string
    Subnets        []string
}

type ECSRunTask struct {
  Cluster 					   string
  Count 							 int64
  EnableECSManagedTags bool
  LaunchType 				   string
  NetworkConfiguration NetworkConfiguration
  PlatformVersion			 string
	Tags								 []KeyValuePair
  TaskDefinition 			 string
}

func GetTaskDefinition(configPath string, configFile string) TaskDefinition {
  conf := TaskDefinition{}
	conf = fetchConfigFromFile(configPath, configFile, conf).(TaskDefinition)

  return conf
}

func GetECSConfig(configPath string, configFile string) ECSRunTask {
  conf := ECSRunTask{}
	conf = fetchConfigFromFile(configPath, configFile, conf).(ECSRunTask)

  return conf
}
