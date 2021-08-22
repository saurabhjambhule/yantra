package aws

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/saurabhjambhule/yantra/pkg/config"
	"github.com/saurabhjambhule/yantra/internal/utils"
)

const (
	ecsConfigFile  = "config"
	taskDefinitionFile  = "task_defination"
)

func checkAndSetAwsString(input string) *string {
	if utils.IsStringEmpty(input) {
		return nil
	}	else {
		return setAwsString(input)
	}
}

func checkAndSetAwsInt(input int64) *int64 {
	if utils.IsIntEmpty(input) {
		return nil
	}	else {
		return setAwsInt(input)
	}
}

func setTags(input []config.KeyValuePair) []*ecs.Tag {
	var tags []*ecs.Tag
	for _, t := range input {
		tag := &ecs.Tag {
				Key:   setAwsString(t.Name),
				Value: setAwsString(t.Value),
		}
		tags = append(tags, tag)
	}

	return tags
}

func setContainerDefinition(input config.ContainerDefinition) *ecs.ContainerDefinition {
	var command []*string
	for _, c := range input.Command {
		command = append(command, setAwsString(c))
	}

	var entryPoint []*string
	for _, ep := range input.EntryPoint {
		entryPoint = append(entryPoint, setAwsString(ep))
	}

	var portMappings []*ecs.PortMapping
	for _, p := range input.PortMappings {
		portMapping := &ecs.PortMapping {
				ContainerPort: setAwsInt(p.ContainerPort),
				HostPort:      setAwsInt(p.HostPort),
				Protocol:      setAwsString(p.Protocol),
		}
		portMappings = append(portMappings, portMapping)
	}

	var dependencies []*ecs.ContainerDependency
	for _, d := range input.DependsOn {
		dependency := &ecs.ContainerDependency {
				Condition: 		 setAwsString(d.Condition),
				ContainerName: setAwsString(d.ContainerName),
		}
		dependencies = append(dependencies, dependency)
	}

	var mountPoints []*ecs.MountPoint
	for _, m := range input.MountPoints {
		mountPoint := &ecs.MountPoint {
				ContainerPath: setAwsString(m.ContainerPath),
				ReadOnly: 		 aws.Bool(m.ReadOnly),
				SourceVolume:  setAwsString(m.SourceVolume),
		}
		mountPoints = append(mountPoints, mountPoint)
	}

	h := input.HealthCheck
	var healthCommand []*string
	var healthCheck *ecs.HealthCheck = nil
	if !utils.IsStringSliceEmpty(h.Command) {
		for _, hc := range h.Command {
			healthCommand = append(healthCommand, setAwsString(hc))
		}
		healthCheck = &ecs.HealthCheck {
				Command:     healthCommand,
				Interval:    setAwsInt(h.Interval),
			  Retries:     setAwsInt(h.Retries),
			  StartPeriod: setAwsInt(h.StartPeriod),
			  Timeout:     setAwsInt(h.Timeout),
		}
	}

	l := input.LogConfiguration
	logConfiguration := &ecs.LogConfiguration {
			LogDriver: setAwsString(l.LogDriver),
			Options:   l.Options,
	}

	var environments []*ecs.KeyValuePair
	for _, e := range input.Environment {
		environment := &ecs.KeyValuePair {
				Name:  setAwsString(e.Name),
				Value: setAwsString(e.Value),
		}
		environments = append(environments, environment)
	}

	image := checkAndSetAwsString(input.Image)
	name := checkAndSetAwsString(input.Name)
	cpu := checkAndSetAwsInt(input.Cpu)
	memory := checkAndSetAwsInt(input.Memory)
	memoryReservation := checkAndSetAwsInt(input.MemoryReservation)
	stopTimeout := checkAndSetAwsInt(input.StopTimeout)
	startTimeout := checkAndSetAwsInt(input.StartTimeout)

	containerDefinition := &ecs.ContainerDefinition{
		Image:     				 image,
		Name:      				 name,
		Command: 					 command,
		Environment: 			 environments,
		EntryPoint: 			 entryPoint,
		Cpu:       				 cpu,
		Memory:    				 memory,
		MemoryReservation: memoryReservation,
		Essential: 				 aws.Bool(input.Essential),
		PortMappings: 		 portMappings,
		DependsOn: 				 dependencies,
		MountPoints: 			 mountPoints,
		HealthCheck:   		 healthCheck,
		LogConfiguration:  logConfiguration,
		StartTimeout:      startTimeout,
		StopTimeout:       stopTimeout,
	}

	return containerDefinition
}

func createTaskDefination(client *ecs.ECS, input config.TaskDefinition) string {
	var containerDefinitions []*ecs.ContainerDefinition
	for _, cd := range input.ContainerDefinitions {
		containerDefinitions = append(containerDefinitions, setContainerDefinition(cd))
	}

	family := checkAndSetAwsString(input.Family)
	networkMode := checkAndSetAwsString(input.NetworkMode)
	cpu := checkAndSetAwsString(input.Cpu)
	memory := checkAndSetAwsString(input.Memory)
	executionRoleArn := checkAndSetAwsString(input.ExecutionRoleArn)
	taskRoleArn := checkAndSetAwsString(input.TaskRoleArn)

	var compatibilities []*string
	for _, rc := range input.RequiresCompatibilities {
		compatibilities = append(compatibilities, setAwsString(rc))
	}

	var volumes []*ecs.Volume
	for _, v := range input.Volumes {
		volume := &ecs.Volume {
				Name:  setAwsString(v.Name),
		}
		volumes = append(volumes, volume)
	}

	tags := setTags(input.Tags)

  taskDefinition := &ecs.RegisterTaskDefinitionInput{
    ContainerDefinitions:    containerDefinitions,
	  Cpu: 										 cpu,
		Memory: 								 memory,
	  ExecutionRoleArn: 			 executionRoleArn,
		TaskRoleArn: 					   taskRoleArn,
	  Family: 								 family,
	  NetworkMode: 						 networkMode,
	  RequiresCompatibilities: compatibilities,
	  Tags:										 tags,
	  Volumes: 								 volumes,
  }

  result, err := client.RegisterTaskDefinition(taskDefinition)
  if err != nil {
    if aerr, ok := err.(awserr.Error); ok {
      switch aerr.Code() {
      case ecs.ErrCodeServerException:
        fmt.Println(ecs.ErrCodeServerException, aerr.Error())
      case ecs.ErrCodeClientException:
        fmt.Println(ecs.ErrCodeClientException, aerr.Error())
      case ecs.ErrCodeInvalidParameterException:
        fmt.Println(ecs.ErrCodeInvalidParameterException, aerr.Error())
      default:
        fmt.Println(aerr.Error())
      }
    } else {
      fmt.Println(err.Error())
    }
    log.Fatal(err)
  }

  return *result.TaskDefinition.TaskDefinitionArn
}

func setECSTask(input config.ECSRunTask, taskDefinitionArn string, startedByInput string) *ecs.RunTaskInput {
	cluster := checkAndSetAwsString(input.Cluster)
	count := checkAndSetAwsInt(input.Count)
	launchType := checkAndSetAwsString(input.LaunchType)
	platformVersion := checkAndSetAwsString(input.PlatformVersion)
	taskDefinition := checkAndSetAwsString(taskDefinitionArn)
	startedBy := checkAndSetAwsString(startedByInput)
	tags := setTags(input.Tags)

	var securityGroups []*string
	for _, sg := range input.NetworkConfiguration.AwsvpcConfiguration.SecurityGroups {
		securityGroups = append(securityGroups, setAwsString(sg))
	}
	var subnets []*string
	for _, sn := range input.NetworkConfiguration.AwsvpcConfiguration.Subnets {
		subnets = append(subnets, setAwsString(sn))
	}
	assignPublicIp := checkAndSetAwsString(input.NetworkConfiguration.AwsvpcConfiguration.AssignPublicIp)

	awsVpcConfig := &ecs.AwsVpcConfiguration{
		AssignPublicIp: assignPublicIp,
		SecurityGroups: securityGroups,
		Subnets:        subnets,
	}
	networkConfig := &ecs.NetworkConfiguration{
		AwsvpcConfiguration: awsVpcConfig,
	}

	ecsRunTaskInput := &ecs.RunTaskInput{
    Cluster:        			cluster,
    TaskDefinition: 			taskDefinition,
		Count: 				        count,
		EnableECSManagedTags: aws.Bool(input.EnableECSManagedTags),
		LaunchType: 				  launchType,
		PlatformVersion: 			platformVersion,
		StartedBy: 						startedBy,
		NetworkConfiguration: networkConfig,
		Tags: 								tags,
	}

	return ecsRunTaskInput
}

func getECSTaskStatus(client *ecs.ECS, clusterInput string, taskInput string) {
	var tasks []*string
	tasks = append(tasks, setAwsString(taskInput))
	cluster := checkAndSetAwsString(clusterInput)

	input := &ecs.DescribeTasksInput{
    Cluster: cluster,
    Tasks:   tasks,
	}

	err := client.WaitUntilTasksRunning(input)
	if err != nil {
		log.Fatal(err)
	}
}

func isTaskRunningStartedBy(client *ecs.ECS, cluster *string, startedBy *string) (bool, []*string) {
	input := &ecs.ListTasksInput{
	    Cluster: 	 cluster,
			StartedBy: startedBy,
	}

	result, err := client.ListTasks(input)
	if err != nil {
    if aerr, ok := err.(awserr.Error); ok {
      switch aerr.Code() {
      case ecs.ErrCodeServerException:
        fmt.Println(ecs.ErrCodeServerException, aerr.Error())
      case ecs.ErrCodeClientException:
        fmt.Println(ecs.ErrCodeClientException, aerr.Error())
      case ecs.ErrCodeInvalidParameterException:
        fmt.Println(ecs.ErrCodeInvalidParameterException, aerr.Error())
      case ecs.ErrCodeClusterNotFoundException:
        fmt.Println(ecs.ErrCodeClusterNotFoundException, aerr.Error())
      case ecs.ErrCodeServiceNotFoundException:
        fmt.Println(ecs.ErrCodeServiceNotFoundException, aerr.Error())
      default:
        fmt.Println(aerr.Error())
      }
    } else {
      fmt.Println(err.Error())
    }
    log.Fatal(err)
	}

	if utils.IsStringSliceEmpty(aws.StringValueSlice(result.TaskArns)) {
		return false, result.TaskArns
	} else {
		return true, result.TaskArns
	}
}

func stopRunningTasks(client *ecs.ECS, cluster *string, tasks []*string) {
	for _, t := range tasks {
		input := &ecs.StopTaskInput{
	    Cluster: cluster,
	    Task: t,
		}

		_, err := client.StopTask(input)
		if err != nil {
	    if aerr, ok := err.(awserr.Error); ok {
        switch aerr.Code() {
        case ecs.ErrCodeServerException:
          fmt.Println(ecs.ErrCodeServerException, aerr.Error())
        case ecs.ErrCodeClientException:
          fmt.Println(ecs.ErrCodeClientException, aerr.Error())
        case ecs.ErrCodeInvalidParameterException:
          fmt.Println(ecs.ErrCodeInvalidParameterException, aerr.Error())
        case ecs.ErrCodeClusterNotFoundException:
          fmt.Println(ecs.ErrCodeClusterNotFoundException, aerr.Error())
        case ecs.ErrCodeServiceNotFoundException:
          fmt.Println(ecs.ErrCodeServiceNotFoundException, aerr.Error())
        default:
          fmt.Println(aerr.Error())
        }
	    } else {
        fmt.Println(err.Error())
	    }
	    log.Fatal(err)
		}
	}
}

func getTaskIP(client *ecs.ECS, cluster *string, task *string) string {
	input := &ecs.DescribeTasksInput{
		Cluster: cluster,
    Tasks: []*string{
      task,
    },
	}

	result, err := client.DescribeTasks(input)
	if err != nil {
    if aerr, ok := err.(awserr.Error); ok {
      switch aerr.Code() {
      case ecs.ErrCodeServerException:
        fmt.Println(ecs.ErrCodeServerException, aerr.Error())
      case ecs.ErrCodeClientException:
        fmt.Println(ecs.ErrCodeClientException, aerr.Error())
      case ecs.ErrCodeInvalidParameterException:
        fmt.Println(ecs.ErrCodeInvalidParameterException, aerr.Error())
      case ecs.ErrCodeClusterNotFoundException:
        fmt.Println(ecs.ErrCodeClusterNotFoundException, aerr.Error())
      default:
        fmt.Println(aerr.Error())
      }
    } else {
      fmt.Println(err.Error())
    }
    log.Fatal(err)
	}

	return *result.Tasks[0].Containers[0].NetworkInterfaces[0].PrivateIpv4Address
}

func RunECSTask(session *session.Session, ecsConfigDir string, startedBy string) string {
	client := ecs.New(session)

	taskDefinitionConfig := config.GetTaskDefinition(ecsConfigDir, taskDefinitionFile)
	taskDefinitionArn := createTaskDefination(client, taskDefinitionConfig)

	ecsConfig := config.GetECSConfig(ecsConfigDir, ecsConfigFile)
	ecsRunTaskInput := setECSTask(ecsConfig, taskDefinitionArn, startedBy)

	isRunning, existingTasks := isTaskRunningStartedBy(client, ecsRunTaskInput.Cluster, ecsRunTaskInput.StartedBy)

	result, err := client.RunTask(ecsRunTaskInput)
	if err != nil {
    if aerr, ok := err.(awserr.Error); ok {
      switch aerr.Code() {
      case ecs.ErrCodeServerException:
        fmt.Println(ecs.ErrCodeServerException, aerr.Error())
      case ecs.ErrCodeClientException:
        fmt.Println(ecs.ErrCodeClientException, aerr.Error())
    	case ecs.ErrCodeInvalidParameterException:
        fmt.Println(ecs.ErrCodeInvalidParameterException, aerr.Error())
      case ecs.ErrCodeClusterNotFoundException:
        fmt.Println(ecs.ErrCodeClusterNotFoundException, aerr.Error())
      case ecs.ErrCodeUnsupportedFeatureException:
        fmt.Println(ecs.ErrCodeUnsupportedFeatureException, aerr.Error())
      case ecs.ErrCodePlatformUnknownException:
        fmt.Println(ecs.ErrCodePlatformUnknownException, aerr.Error())
      case ecs.ErrCodePlatformTaskDefinitionIncompatibilityException:
        fmt.Println(ecs.ErrCodePlatformTaskDefinitionIncompatibilityException, aerr.Error())
      case ecs.ErrCodeAccessDeniedException:
        fmt.Println(ecs.ErrCodeAccessDeniedException, aerr.Error())
      case ecs.ErrCodeBlockedException:
        fmt.Println(ecs.ErrCodeBlockedException, aerr.Error())
      default:
        fmt.Println(aerr.Error())
      }
    } else {
      fmt.Println(err.Error())
    }
    log.Fatal(err)
	}

	getECSTaskStatus(client, *result.Tasks[0].ClusterArn, *result.Tasks[0].TaskArn)
	taskIP := getTaskIP(client, result.Tasks[0].ClusterArn, result.Tasks[0].TaskArn)

	if isRunning {
		stopRunningTasks(client, ecsRunTaskInput.Cluster, existingTasks)
	}

	return taskIP
}
