package aws

import (
	"fmt"
  "log"

  "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/aws/awserr"
  "github.com/aws/aws-sdk-go/service/route53"
	"github.com/saurabhjambhule/yantra/pkg/config"
)

func CreateRoute53RecordSet(session *session.Session, ecsConfigDir string, domainName string, ipAddress string)  {
  client := route53.New(session)
  _=client

  route53Config := config.GetRoute53Config(ecsConfigDir, ecsConfigFile)
  fmt.Println(route53Config)

  domainName = route53Config.DomainPrefix + domainName + route53Config.DomainSuffix

  fmt.Println(domainName)

  input := &route53.ChangeResourceRecordSetsInput{
    ChangeBatch: &route53.ChangeBatch{
      Changes: []*route53.Change{
        {
          Action: aws.String("UPSERT"),
          ResourceRecordSet: &route53.ResourceRecordSet{
            Name: aws.String(domainName),
            ResourceRecords: []*route53.ResourceRecord{
              {
                Value: aws.String(ipAddress),
              },
            },
            TTL:  aws.Int64(300),
            Type: aws.String("A"),
          },
        },
      },
    },
    HostedZoneId: checkAndSetAwsString(route53Config.Route53HostedZone),
  }

  result, err := client.ChangeResourceRecordSets(input)
  if err != nil {
    if aerr, ok := err.(awserr.Error); ok {
      switch aerr.Code() {
      case route53.ErrCodeNoSuchHostedZone:
          fmt.Println(route53.ErrCodeNoSuchHostedZone, aerr.Error())
      case route53.ErrCodeNoSuchHealthCheck:
          fmt.Println(route53.ErrCodeNoSuchHealthCheck, aerr.Error())
      case route53.ErrCodeInvalidChangeBatch:
          fmt.Println(route53.ErrCodeInvalidChangeBatch, aerr.Error())
      case route53.ErrCodeInvalidInput:
          fmt.Println(route53.ErrCodeInvalidInput, aerr.Error())
      case route53.ErrCodePriorRequestNotComplete:
          fmt.Println(route53.ErrCodePriorRequestNotComplete, aerr.Error())
      default:
          fmt.Println(aerr.Error())
      }
    } else {
      fmt.Println(err.Error())
    }
    log.Fatal(err)
  }

  fmt.Println(result)
}
