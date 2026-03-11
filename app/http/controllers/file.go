package controllers

import (
	"github.com/desygner/desygner-toolchain/app/handlers"
	"github.com/desygner/desygner-toolchain/app/http/requests"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

type ReplaceInput struct {
	Source       string `json:"url"`
	Target       string `json:"target"`
	SlackChannel string `json:"slack_channel"`
	SlackTS      string `json:"slack_ts"`
}
type File struct {
	// Dependent services
	fileHandler         *handlers.Replace
	invalidationHandler *handlers.Invalidate
}

func NewFile() *File {
	return &File{
		// Inject services
		fileHandler: &handlers.Replace{},
	}
}

func (r *File) Replace(ctx http.Context) http.Response {
	var req requests.FileReplace
	errors, err := ctx.Request().ValidateRequest(&req)
	if err != nil {
		return ctx.Response().Json(http.StatusBadRequest, map[string]any{
			"error": err.Error(),
		})
	}
	if errors != nil {
		return ctx.Response().Json(http.StatusBadRequest, map[string]any{
			"errors": errors.All(),
		})
	}
	facades.Log().Info(req.Source, req.Target, req)

	input := &handlers.ReplaceInput{
		Target: req.Target,
		Source: req.Source,
	}
	input.SlackChannel = req.SlackChannel
	input.SlackTS = req.SlackTS
	err = r.fileHandler.Handle(ctx.Context(), input)
	if err != nil {
		facades.Log().Error(err)
		return ctx.Response().Json(http.StatusInternalServerError, map[string]any{
			"error": err.Error(),
		})
	}
	return ctx.Response().Json(http.StatusOK, map[string]any{
		"message": "ok",
	})
}
