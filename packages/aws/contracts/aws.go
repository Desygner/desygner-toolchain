package contracts

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

type Aws interface {
	With(profile string) Aws
	GetCredential() aws.Config
	GetProfile() string

	// clients shortcuts

	S3() *s3.Client
	SQS() *sqs.Client
	Cloudfront() *cloudfront.Client
	Ses() *ses.Client
	Cloudwatch() *cloudwatch.Client
	STS() *sts.Client
}
