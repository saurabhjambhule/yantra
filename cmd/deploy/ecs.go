package deploy

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/saurabhjambhule/yantra/pkg/aws/ecr"
	"github.com/saurabhjambhule/yantra/pkg/aws/ecs"
)

func checkImageFromECR(session *session.Session, imageTag string, repoName string) string {
	_, createdAt := ecr.DoesImageExist(session, imageTag, repoName)

	return createdAt
}

func runECSTask()  {
  ecs.RunECSTask()
}
