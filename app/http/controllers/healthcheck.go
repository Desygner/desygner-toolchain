package controllers

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

type Healthcheck struct {
	// Dependent services
}

func NewHealthcheck() *Healthcheck {
	return &Healthcheck{
		// Inject services
	}
}

func (r *Healthcheck) Index(ctx http.Context) http.Response {
	cfg := facades.Config()
	return ctx.Response().Success().Json(http.Json{
		"status":  "ok",
		"name":    cfg.GetString("app.name"),
		"version": cfg.GetString("app.version"),
	})
}
