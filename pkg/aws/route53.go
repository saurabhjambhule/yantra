package aws

// import (
// 	"fmt"
//
// 	"github.com/aws/aws-sdk-go/aws/session"
// 	"github.com/saurabhjambhule/yantra/pkg/aws"
// 	"github.com/saurabhjambhule/yantra/pkg/config"
// )
//
// func AddRecordSet(session *session.Session, domainName string, ipAddress string)  {
//   svc := route53.New(session)
//
//   input := &route53.ChangeResourceRecordSetsInput{
//       ChangeBatch: &route53.ChangeBatch{
//           Changes: []*route53.Change{
//               {
//                   Action: aws.String("UPSERT"),
//                   ResourceRecordSet: &route53.ResourceRecordSet{
//                       Name: aws.String(),
//                       ResourceRecords: []*route53.ResourceRecord{
//                           {
//                               Value: aws.String("192.0.2.44"),
//                           },
//                       },
//                       TTL:  aws.Int64(60),
//                       Type: aws.String("A"),
//                   },
//               },
//           },
//           Comment: aws.String("Web server for example.com"),
//       },
//       HostedZoneId: aws.String("Z3M3LMPEXAMPLE"),
//   }
//
//   result, err := svc.ChangeResourceRecordSets(input)
//   if err != nil {
//       if aerr, ok := err.(awserr.Error); ok {
//           switch aerr.Code() {
//           case route53.ErrCodeNoSuchHostedZone:
//               fmt.Println(route53.ErrCodeNoSuchHostedZone, aerr.Error())
//           case route53.ErrCodeNoSuchHealthCheck:
//               fmt.Println(route53.ErrCodeNoSuchHealthCheck, aerr.Error())
//           case route53.ErrCodeInvalidChangeBatch:
//               fmt.Println(route53.ErrCodeInvalidChangeBatch, aerr.Error())
//           case route53.ErrCodeInvalidInput:
//               fmt.Println(route53.ErrCodeInvalidInput, aerr.Error())
//           case route53.ErrCodePriorRequestNotComplete:
//               fmt.Println(route53.ErrCodePriorRequestNotComplete, aerr.Error())
//           default:
//               fmt.Println(aerr.Error())
//           }
//       } else {
//           // Print the error, cast err to awserr.Error to get the Code and
//           // Message from an error.
//           fmt.Println(err.Error())
//       }
//       return
//   }
//
//   fmt.Println(result)
// }
