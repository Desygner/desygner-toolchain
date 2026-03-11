package slack

import "github.com/slack-go/slack"

type Slack struct {
	cli *slack.Client
}

func NewSlack(client *slack.Client) *Slack {
	return &Slack{client}
}

func (receiver *Slack) Client() *slack.Client {
	return receiver.cli
}
