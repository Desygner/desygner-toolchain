package config

import (
	"github.com/goravel/framework/facades"
)

func init() {
	config := facades.Config()
	config.Add("mcp", map[string]any{
		"port": config.Env("MCP_PORT", 8080),
	})
}
