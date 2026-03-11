package mcp

import (
	"github.com/desygner/desygner-toolchain/packages/mcp/service"
	"github.com/goravel/framework/contracts/foundation"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const Binding = "mcp"

var App foundation.Application

type ServiceProvider struct {
}

func (receiver *ServiceProvider) Register(app foundation.Application) {
	App = app

	app.Singleton(Binding, func(app foundation.Application) (any, error) {
		config := app.MakeConfig()
		srv := mcp.NewServer(&mcp.Implementation{
			Name:    config.GetString("app.name"),
			Version: config.GetString("app.version"),
		}, nil)
		return service.NewMcp(srv, nil), nil
	})
}

func (receiver *ServiceProvider) Boot(app foundation.Application) {
	app.Publishes("github.com/desygner/desygner-toolchain/packages/mcp", map[string]string{
		"config/mcp.go": app.ConfigPath("mcp.go"),
	})

}
