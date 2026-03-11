package commands

import (
	"context"
	"github.com/desygner/desygner-toolchain/app/handlers"
	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/facades"
)

type CloudfrontInvalidate struct {
}

// Signature The name and signature of the console command.
func (receiver *CloudfrontInvalidate) Signature() string {
	return "cloudfront:invalidate"
}

// Description The console command description.
func (receiver *CloudfrontInvalidate) Description() string {
	return "Send Invalidation Request to Cloudfront"
}

// Extend The console command extend.
func (receiver *CloudfrontInvalidate) Extend() command.Extend {
	return command.Extend{
		ArgsUsage: " <url>",
		Category:  "cloudfront",
		Flags: []command.Flag{
			&command.StringFlag{
				Name:    "slack_channel",
				Aliases: []string{"ch"},
				Usage:   "Slack Channel",
			},
			&command.StringFlag{
				Name:    "slack_ts",
				Aliases: []string{"ts"},
				Usage:   "Slack Thread TS",
			},
		},
	}
}

// Handle Execute the console command.
func (receiver *CloudfrontInvalidate) Handle(ctx console.Context) error {
	url := ctx.Argument(0)
	var input *handlers.InvalidationInput
	slackChannel := ctx.Option("slack_channel")
	slackTS := ctx.Option("slack_ts")
	logger := facades.Log()
	logger.With(map[string]any{"ch": slackChannel, "ts": slackTS}).Info("check if it has slack channel and ts")
	input = &handlers.InvalidationInput{
		URI: url,
	}
	if slackChannel != "" && slackTS != "" {
		input.SlackChannel = slackChannel
		input.SlackTS = slackTS
	}
	_, err := (&handlers.Invalidate{}).Handle(context.Background(), input)
	return err
}
