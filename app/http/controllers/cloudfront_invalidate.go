package controllers

import (
	"github.com/desygner/desygner-toolchain/app/handlers"
	"github.com/desygner/desygner-toolchain/app/http/requests"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/log"
	"github.com/goravel/framework/facades"
)

type CloudfrontInvalidate struct {
	// Dependent services
	handler *handlers.Invalidate
	logger  log.Log
}

func NewCloudfrontInvalidate() *CloudfrontInvalidate {
	return &CloudfrontInvalidate{
		// Inject services
		logger:  facades.Log(),
		handler: &handlers.Invalidate{},
	}
}

func (r *CloudfrontInvalidate) Invalidate(ctx http.Context) http.Response {
	var input requests.CloudfrontInvalidate
	errors, err := ctx.Request().ValidateRequest(&input)
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
	invalidationInput := &handlers.InvalidationInput{
		DistributionID: input.DistributionID,
		URI:            input.URI,
	}
	if input.SlackTS != "" && input.SlackChannel != "" {
		invalidationInput.SlackChannel = input.SlackChannel
		invalidationInput.SlackTS = input.SlackTS
	}
	res, err := r.handler.Handle(ctx.Context(), invalidationInput)
	if err != nil {
		r.logger.Error(err)
		return ctx.Response().Json(http.StatusInternalServerError, map[string]any{
			"error": err.Error(),
		})
	}
	return ctx.Response().Json(http.StatusOK, map[string]any{
		"invalidation_id": res.InvalidationID,
		"distribution_id": res.DistributionID,
	})
}
func (r *CloudfrontInvalidate) GetBucket(ctx http.Context) http.Response {
	var input requests.DistributionDiscovery
	errors, err := ctx.Request().ValidateRequest(&input)
	if err != nil {
		return ctx.Response().Json(http.StatusBadRequest, map[string]any{
			"errors": errors,
		})
	}
	if errors != nil {
		return ctx.Response().Json(http.StatusBadRequest, map[string]any{
			"errors": errors.All(),
		})
	}
	res, err := r.handler.ResolveBucket(ctx.Context(), input.URL)
	if err != nil {
		r.logger.Error(err)
		return ctx.Response().Json(http.StatusInternalServerError, map[string]any{
			"error": err.Error(),
		})
	}
	return ctx.Response().Json(http.StatusOK, map[string]any{
		"bucket":          res.Bucket,
		"distribution_id": res.DistributionID,
	})
}
