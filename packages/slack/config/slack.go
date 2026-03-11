package config

import (
	"github.com/goravel/framework/facades"
)

func init() {
	config := facades.Config()
	config.Add("slack", map[string]any{
		"bot_token": config.Env("SLACK_BOT_TOKEN", ""),
	})
}
