package aws

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
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
