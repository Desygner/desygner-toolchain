package client

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/desygner/desygner-toolchain/packages/aws/contracts"
	"github.com/goravel/framework/contracts/cache"
	config2 "github.com/goravel/framework/contracts/config"
)

type AWS struct {
	profile string
	cache   cache.Cache
	cfg     config2.Config
}

func (A *AWS) GetProfile() string {
	return A.profile
}

func New(cfg config2.Config, cacheClient cache.Cache) contracts.Aws {
	return &AWS{
		profile: cfg.GetString("aws.default"),
		cfg:     cfg,
		cache:   cacheClient,
	}
}

func (A *AWS) With(profile string) contracts.Aws {
	return &AWS{
		profile: profile,
	}
}

func (A *AWS) GetCredential() aws.Config {
	region := A.cfg.GetString(fmt.Sprintf("aws.profiles.%s.region", A.profile))
	roleArn := A.cfg.GetString(fmt.Sprintf("aws.profiles.%s.roleArn", A.profile))
	if region == "" {
		region = "us-east-1"
	}
	//customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
	//	if service == s3.ServiceID {
	//		return aws.Endpoint{
	//			URL:           "https://s3." + region + ".amazonaws.com",
	//			SigningRegion: region,
	//		}, nil
	//	}
	//	return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	//})
	opts := []func(*config.LoadOptions) error{
		config.WithRegion(region),
	}
	if roleArn != "" {
		opts = append(opts, config.WithCredentialsProvider(
			stscreds.NewAssumeRoleProvider(sts.NewFromConfig(A.GetCredential()), roleArn),
		))
	}

	cfg, _ := config.LoadDefaultConfig(context.Background(), opts...)
	return cfg
}

func (A *AWS) S3() *s3.Client {
	return s3.NewFromConfig(A.GetCredential(), func(o *s3.Options) {
		o.UsePathStyle = true
	})
}

func (A *AWS) SQS() *sqs.Client {
	return sqs.NewFromConfig(A.GetCredential())
}

func (A *AWS) Cloudfront() *cloudfront.Client {
	return cloudfront.NewFromConfig(A.GetCredential())
}

func (A *AWS) Ses() *ses.Client {
	return ses.NewFromConfig(A.GetCredential())
}

func (A *AWS) Cloudwatch() *cloudwatch.Client {
	return cloudwatch.NewFromConfig(A.GetCredential())
}

func (A *AWS) STS() *sts.Client {
	return sts.NewFromConfig(A.GetCredential())
}
