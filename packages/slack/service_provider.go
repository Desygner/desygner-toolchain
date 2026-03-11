package slack

import (
	"github.com/desygner/desygner-toolchain/app/extensions"
	"github.com/goravel/framework/contracts/foundation"
	"github.com/slack-go/slack"
)

const Binding = "slack"

var App foundation.Application

type ServiceProvider struct {
}

func (receiver *ServiceProvider) Register(app foundation.Application) {
	App = app

	app.Singleton(Binding, func(app foundation.Application) (any, error) {
		cfg := app.MakeConfig()
		botToken := cfg.GetString("slack.bot_token")
		logger := app.MakeLog()
		slackClient := slack.New(botToken, slack.OptionLog(extensions.NewSlackLogger(logger)))
		return NewSlack(slackClient), nil
	})
}

func (receiver *ServiceProvider) Boot(app foundation.Application) {
	app.Publishes("github.com/desygner/desygner-toolchain/packages/slack", map[string]string{
		"config/slack.go": app.ConfigPath("slack.go"),
	})

}
