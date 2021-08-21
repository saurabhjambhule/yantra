package ecs

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

func setAwsString(input string) *string {
	if utils.IsStringEmpty(input) {
		return nil
	}	else {
		return aws.String(input)
	}
}

func setAwsInt(input int64) *int64 {
	if utils.IsIntEmpty(input) {
		return nil
	}	else {
		return aws.Int64(input)
	}
}

func setTags(input []config.KeyValuePair) []*ecs.Tag {
	var tags []*ecs.Tag
	for _, t := range input {
		tag := &ecs.Tag {
				Key:   aws.String(t.Name),
				Value: aws.String(t.Value),
		}
		tags = append(tags, tag)
	}

	return tags
}

func setContainerDefinition(input config.ContainerDefinition) *ecs.ContainerDefinition {
	var command []*string
	for _, c := range input.Command {
		command = append(command, aws.String(c))
	}

	var entryPoint []*string
	for _, ep := range input.EntryPoint {
		entryPoint = append(entryPoint, aws.String(ep))
	}

	var portMappings []*ecs.PortMapping
	for _, p := range input.PortMappings {
		portMapping := &ecs.PortMapping {
				ContainerPort: aws.Int64(p.ContainerPort),
				HostPort:      aws.Int64(p.HostPort),
				Protocol:      aws.String(p.Protocol),
		}
		portMappings = append(portMappings, portMapping)
	}

	var dependencies []*ecs.ContainerDependency
	for _, d := range input.DependsOn {
		dependency := &ecs.ContainerDependency {
				Condition: 		 aws.String(d.Condition),
				ContainerName: aws.String(d.ContainerName),
		}
		dependencies = append(dependencies, dependency)
	}

	var mountPoints []*ecs.MountPoint
	for _, m := range input.MountPoints {
		mountPoint := &ecs.MountPoint {
				ContainerPath: aws.String(m.ContainerPath),
				ReadOnly: 		 aws.Bool(m.ReadOnly),
				SourceVolume:  aws.String(m.SourceVolume),
		}
		mountPoints = append(mountPoints, mountPoint)
	}

	h := input.HealthCheck
	var healthCommand []*string
	var healthCheck *ecs.HealthCheck = nil
	if !utils.IsStringSliceEmpty(h.Command) {
		for _, hc := range h.Command {
			healthCommand = append(healthCommand, aws.String(hc))
		}
		healthCheck = &ecs.HealthCheck {
				Command:     healthCommand,
				Interval:    aws.Int64(h.Interval),
			  Retries:     aws.Int64(h.Retries),
			  StartPeriod: aws.Int64(h.StartPeriod),
			  Timeout:     aws.Int64(h.Timeout),
		}
	}

	l := input.LogConfiguration
	logConfiguration := &ecs.LogConfiguration {
			LogDriver: aws.String(l.LogDriver),
			Options:   l.Options,
	}

	var environments []*ecs.KeyValuePair
	for _, e := range input.Environment {
		environment := &ecs.KeyValuePair {
				Name:  aws.String(e.Name),
				Value: aws.String(e.Value),
		}
		environments = append(environments, environment)
	}

	image := setAwsString(input.Image)
	name := setAwsString(input.Name)
	cpu := setAwsInt(input.Cpu)
	memory := setAwsInt(input.Memory)
	memoryReservation := setAwsInt(input.MemoryReservation)
	stopTimeout := setAwsInt(input.StopTimeout)
	startTimeout := setAwsInt(input.StartTimeout)

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

	family := setAwsString(input.Family)
	networkMode := setAwsString(input.NetworkMode)
	cpu := setAwsString(input.Cpu)
	memory := setAwsString(input.Memory)
	executionRoleArn := setAwsString(input.ExecutionRoleArn)
	taskRoleArn := setAwsString(input.TaskRoleArn)

	var compatibilities []*string
	for _, rc := range input.RequiresCompatibilities {
		compatibilities = append(compatibilities, aws.String(rc))
	}

	var volumes []*ecs.Volume
	for _, v := range input.Volumes {
		volume := &ecs.Volume {
				Name:  aws.String(v.Name),
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
	cluster := setAwsString(input.Cluster)
	count := setAwsInt(input.Count)
	launchType := setAwsString(input.LaunchType)
	platformVersion := setAwsString(input.PlatformVersion)
	taskDefinition := setAwsString(taskDefinitionArn)
	startedBy := setAwsString(startedByInput)
	tags := setTags(input.Tags)

	var securityGroups []*string
	for _, sg := range input.NetworkConfiguration.AwsvpcConfiguration.SecurityGroups {
		securityGroups = append(securityGroups, aws.String(sg))
	}
	var subnets []*string
	for _, sn := range input.NetworkConfiguration.AwsvpcConfiguration.Subnets {
		subnets = append(subnets, aws.String(sn))
	}
	assignPublicIp := setAwsString(input.NetworkConfiguration.AwsvpcConfiguration.AssignPublicIp)

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
	tasks = append(tasks, aws.String(taskInput))
	cluster := setAwsString(clusterInput)

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
