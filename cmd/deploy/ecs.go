package deploy

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/saurabhjambhule/yantra/pkg/aws"
	"github.com/saurabhjambhule/yantra/pkg/aws/ecr"
	"github.com/saurabhjambhule/yantra/pkg/aws/ecs"
)

func checkImageFromECR(session *session.Session, imageTag string, repoName string) string {
	_, createdAt := ecr.DoesImageExist(session, imageTag, repoName)

	return createdAt
}

func runECSTask()  {

	awsProfile := "default"
	awsRegion := "us-east-1"
	session := aws.StartSession(awsProfile, awsRegion)

	// imageTag := getImageTag()
	// createdAt := checkImageFromECR(session, imageTag, "dash-test")
	// fmt.Println("The image: " + imageTag + " created " + createdAt + " before!")
	// utils.UserConfirmation()

	ecs.RunECSTask(session, "/Users/saurabhjambhule/workspace/go/src/github.com/saurabhjambhule/yantra/examples/config/ecs", "track2883")
}
