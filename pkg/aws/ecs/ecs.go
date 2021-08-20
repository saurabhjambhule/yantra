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

func createTaskDefination(session *session.Session, input config.TaskDefinition) string {
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

	var tags []*ecs.Tag
	for _, t := range input.Tags {
		tag := &ecs.Tag {
				Key:   aws.String(t.Name),
				Value: aws.String(t.Value),
		}
		tags = append(tags, tag)
	}

	var volumes []*ecs.Volume
	for _, v := range input.Volumes {
		volume := &ecs.Volume {
				Name:  aws.String(v.Name),
		}
		volumes = append(volumes, volume)
	}

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

	client := ecs.New(session)
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

func RunECSTask(session *session.Session)  {
	input := config.GetTaskDefinition("/Users/saurabhjambhule/workspace/go/src/github.com/saurabhjambhule/yantra/assets/config", "ecs_sample.yaml")

	taskDefinitionArn := createTaskDefination(session, input)
	fmt.Println(taskDefinitionArn)
}
