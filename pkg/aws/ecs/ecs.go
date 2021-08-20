package ecs

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/saurabhjambhule/yantra/pkg/config"
)


func createTaskDefination(session *session.Session)  {
  client := ecs.New(session)

  input := &ecs.RegisterTaskDefinitionInput{
      ContainerDefinitions: []*ecs.ContainerDefinition{
          {
              Command: []*string{
                  aws.String("sleep"),
                  aws.String("360"),
              },
              Cpu:       aws.Int64(10),
              Essential: aws.Bool(true),
              Image:     aws.String("busybox"),
              Memory:    aws.Int64(10),
              Name:      aws.String("sleep"),
          },
      },
      Family:      aws.String("sleep360"),
      TaskRoleArn: aws.String(""),
  }

  result, err := client.RegisterTaskDefinition(input)
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
          // Print the error, cast err to awserr.Error to get the Code and
          // Message from an error.
          fmt.Println(err.Error())
      }
      return
  }

  fmt.Println(result)
}

func RunECSTask()  {

	containerDefinition := config.GetContainerDefinition("/Users/saurabhjambhule/workspace/go/src/github.com/saurabhjambhule/yantra/assets/config", "ecs_sample0.yaml")
	fmt.Println("Database is\t", containerDefinition.Image)

}
