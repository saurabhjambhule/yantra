package aws

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/saurabhjambhule/yantra/pkg/config"
	"github.com/saurabhjambhule/yantra/internal/utils"
)

const (
	ecsConfigFile  = "config"
	taskDefinitionFile  = "task_defination"
)

func StartSession(awsProfile string, awsRegion string) *session.Session {
	session, err := session.NewSessionWithOptions(session.Options{
		Profile: awsProfile,
		Config: aws.Config{
			Region: aws.String(awsRegion),
		},
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		log.Fatal(err)
	}


	return session
}

func setAwsString(input string) *string {
	input = config.UpdatePlaceholder(input)

	return aws.String(input)
}

// TODO: Placeholder for integer values
func setAwsInt(input int64) *int64 {
	return aws.Int64(input)
}

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
