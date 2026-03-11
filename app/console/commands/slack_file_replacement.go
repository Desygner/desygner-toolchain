package commands

import (
	"context"
	"github.com/desygner/desygner-toolchain/app/handlers"
	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
)

type SlackFileReplacement struct {
}

// Signature The name and signature of the console command.
func (receiver *SlackFileReplacement) Signature() string {
	return "slack:file_replace"
}

// Description The console command description.
func (receiver *SlackFileReplacement) Description() string {
	return "Replace file on cloudfront url from slack private file"
}

// Extend The console command extend.
func (receiver *SlackFileReplacement) Extend() command.Extend {
	return command.Extend{
		ArgsUsage: " <source> <target>",
		Category:  "slack",
		Flags: []command.Flag{
			&command.StringFlag{
				Name:    "slack_channel",
				Aliases: []string{"ch"},
				Usage:   "Slack Channel",
			},
			&command.StringFlag{
				Name:    "slack_ts",
				Aliases: []string{"ts"},
			},
		},
	}
}

// Handle Execute the console command.
func (receiver *SlackFileReplacement) Handle(ctx console.Context) error {
	input := &handlers.ReplaceInput{
		Source: ctx.Argument(0),
		Target: ctx.Argument(1),
	}
	input.SlackChannel = ctx.Option("slack_channel")
	input.SlackTS = ctx.Option("slack_ts")
	return (&handlers.Replace{}).Handle(context.Background(), input)
}
