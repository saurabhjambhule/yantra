package aws

import (
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
)

func DoesImageExist(session *session.Session, ImageTag string, repoName string) (bool, string) {
	client := ecr.New(session)

	input := &ecr.DescribeImagesInput{
		ImageIds: []*ecr.ImageIdentifier{
			{
				ImageTag: aws.String(ImageTag),
			},
		},
		RepositoryName: aws.String(repoName),
	}

	imageOutput, err := client.DescribeImages(input)
	if err != nil {
		log.Fatal(err)
	}

	now := time.Now()
	pushedBefore := now.Sub(*imageOutput.ImageDetails[0].ImagePushedAt)

	return true, pushedBefore.String()
}
