package deploy

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/saurabhjambhule/yantra/pkg/aws"
	"github.com/saurabhjambhule/yantra/pkg/aws/ecr"
	"github.com/saurabhjambhule/yantra/pkg/aws/ecs"
	// "github.com/saurabhjambhule/yantra/pkg/config"
	"github.com/saurabhjambhule/yantra/internal/utils"
)

const (
	gitBranch = "YANTRA_GIT_BRANCH"
)

func checkImageFromECR(session *session.Session, imageTag string, repoName string) string {
	_, createdAt := ecr.DoesImageExist(session, imageTag, repoName)

	return createdAt
}

func runECSTask()  {

	awsProfile := "default"
	awsRegion := "us-east-1"
	session := aws.StartSession(awsProfile, awsRegion)

	imageTag := getImageTag()
	os.Setenv(gitBranch, imageTag)
	createdAt := checkImageFromECR(session, imageTag, "dash-test")
	fmt.Println("The image: " + imageTag + " created " + createdAt + " before!")
	utils.UserConfirmation()

	// taskDefinitionConfig := config.GetTaskDefinition("/Users/saurabhjambhule/workspace/go/src/github.com/saurabhjambhule/yantra/examples/config/ecs", "task_defination")
	// fmt.Println(taskDefinitionConfig)

	taskIP := ecs.RunECSTask(session, "/Users/saurabhjambhule/workspace/go/src/github.com/saurabhjambhule/yantra/examples/config/ecs", "saurbh")
	fmt.Println(taskIP)
}
