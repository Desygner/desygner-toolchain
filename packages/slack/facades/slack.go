package facades

import (
	"log"

	"github.com/desygner/desygner-toolchain/packages/slack"
	"github.com/desygner/desygner-toolchain/packages/slack/contracts"
)

func Slack() contracts.Slack {
	instance, err := slack.App.Make(slack.Binding)
	if err != nil {
		log.Println(err)
		return nil
	}

	return instance.(contracts.Slack)
}
