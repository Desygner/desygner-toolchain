package contracts

import "github.com/slack-go/slack"

type Slack interface {
	Client() *slack.Client
}
