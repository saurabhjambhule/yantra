package deploy

import (
	// "fmt"
	_"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/saurabhjambhule/yantra/pkg/aws"
	// "github.com/saurabhjambhule/yantra/pkg/config"
	_"github.com/saurabhjambhule/yantra/internal/utils"
)

const (
	gitBranch = "YANTRA_GIT_BRANCH"
)

func checkImageFromECR(session *session.Session, imageTag string, repoName string) string {
	_, createdAt := aws.DoesImageExist(session, imageTag, repoName)

	return createdAt
}

func runECSTask()  {

	awsProfile := "default"
	awsRegion := "us-east-1"
	session := aws.StartSession(awsProfile, awsRegion)

	// imageTag := getImageTag()
	// os.Setenv(gitBranch, imageTag)
	// createdAt := checkImageFromECR(session, imageTag, "dash-test")
	// fmt.Println("The image: " + imageTag + " created " + createdAt + " before!")
	// utils.UserConfirmation()

	// route53Config := config.GetRoute53Config("/Users/saurabhjambhule/workspace/go/src/github.com/saurabhjambhule/yantra/examples/config/ecs", "config")
	// fmt.Println(route53Config)

	aws.CreateRoute53RecordSet(session, "/Users/saurabhjambhule/workspace/go/src/github.com/saurabhjambhule/yantra/examples/config/ecs", "testing", "10.123.12.1")

	// taskDefinitionConfig := config.GetTaskDefinition("/Users/saurabhjambhule/workspace/go/src/github.com/saurabhjambhule/yantra/examples/config/ecs", "task_defination")
	// fmt.Println(taskDefinitionConfig)

	// taskIP := aws.RunECSTask(session, "/Users/saurabhjambhule/workspace/go/src/github.com/saurabhjambhule/yantra/examples/config/ecs", "saurbh")
	// fmt.Println(taskIP)
}
