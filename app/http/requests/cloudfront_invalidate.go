package requests

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type CloudfrontInvalidate struct {
	DistributionID string `form:"distribution_id" json:"distribution_id"`
	URI            string `form:"uri" json:"uri"`
	SlackChannel   string `form:"slack_channel" json:"slack_channel"`
	SlackTS        string `form:"slack_ts" json:"slack_ts"`
}

func (r *CloudfrontInvalidate) Authorize(ctx http.Context) error {
	return nil
}

func (r *CloudfrontInvalidate) Filters(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *CloudfrontInvalidate) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"distribution_id": "required",
		"uri":             "required|starts_with:/",
	}
}

func (r *CloudfrontInvalidate) Messages(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *CloudfrontInvalidate) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *CloudfrontInvalidate) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
