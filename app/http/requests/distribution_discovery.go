package requests

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type DistributionDiscovery struct {
	URL string `form:"url" json:"url"`
}

func (r *DistributionDiscovery) Authorize(ctx http.Context) error {
	return nil
}

func (r *DistributionDiscovery) Filters(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *DistributionDiscovery) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"url": "required|full_url",
	}
}

func (r *DistributionDiscovery) Messages(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *DistributionDiscovery) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *DistributionDiscovery) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
