package config

import (
	"github.com/goravel/framework/facades"
)

func init() {
	config := facades.Config()
	config.Add("aws", map[string]any{
		"default": config.Env("AWS_CURRENT_PROFILE", "default"),
		"ttl":     config.Env("AWS_TTL", 3600),
		"profiles": map[string]any{
			"default": map[string]any{},
			"desygner_dev": map[string]any{
				"roleArn": config.Env("AWS_DESYGNER_DEV_ROLE_ARN", ""),
			},
			"fluer_dev": map[string]any{
				"roleArn": config.Env("AWS_FLUER_DEV_ROLE_ARN", ""),
			},
			"fluer_prod": map[string]any{
				"roleArn": config.Env("AWS_FLUER_PROD_ROLE_ARN", ""),
			},
			"knowz_test": map[string]any{
				"roleArn": config.Env("AWS_KNOWZ_TEST_ROLE_ARN", ""),
				"region":  config.Env("AWS_KNOWZ_TEST_REGION", "ap-southeast-2"),
			},
			"knowz_prod": map[string]any{
				"roleArn": config.Env("AWS_KNOWZ_PROD_ROLE_ARN", ""),
				"region":  config.Env("AWS_KNOWZ_PROD_REGION", "us-east-1"),
			},
		},
	})
}
