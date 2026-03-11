package aws

import (
	"github.com/desygner/desygner-toolchain/packages/aws/client"
	"github.com/goravel/framework/contracts/foundation"
)

const Binding = "aws"

var App foundation.Application

type ServiceProvider struct {
}

func (receiver *ServiceProvider) Register(app foundation.Application) {
	App = app

	app.Singleton(Binding, func(app foundation.Application) (any, error) {
		return client.New(app.MakeConfig(), app.MakeCache()), nil
	})
}

func (receiver *ServiceProvider) Boot(app foundation.Application) {
	app.Publishes("github.com/desygner/desygner-toolchain/packages/aws", map[string]string{
		"config/aws.go": app.ConfigPath("aws.go"),
	}, "aws-config")
}
