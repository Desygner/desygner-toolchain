package requests

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type FileReplace struct {
	SlackChannel string `form:"slack_channel" json:"slack_channel"`
	SlackTS      string `form:"slack_ts" json:"slack_ts"`
	Source       string `form:"source" json:"source"`
	Target       string `form:"target" json:"target"`
}

func (r *FileReplace) Authorize(ctx http.Context) error {
	return nil
}

func (r *FileReplace) Filters(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *FileReplace) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"source": "required|full_url",
		"target": "required|full_url",
	}
}

func (r *FileReplace) Messages(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *FileReplace) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *FileReplace) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
