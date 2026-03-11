package handlers

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
	facades2 "github.com/desygner/desygner-toolchain/packages/aws/facades"
	"github.com/desygner/desygner-toolchain/packages/slack/extension"
	"github.com/goravel/framework/facades"
	"strconv"
	"time"
)

type Invalidate struct {
	CDN
}

type InvalidationInput struct {
	DistributionID, URI string
	extension.Input
}

type InvalidationResponse struct {
	DistributionID string
	InvalidationID string
}

func (receiver *Invalidate) Handle(ctx context.Context, input *InvalidationInput) (res *InvalidationResponse, err error) {
	defer func() {
		if err != nil {
			input.SendMessage(ctx, fmt.Sprintf(":x: Failed to invalidate %s. Reason: %s", input.URI, err.Error()))
		}
	}()
	cfn := facades2.Aws().Cloudfront()
	cfg := facades.Config()
	caller := cfg.GetString("app.name") + "/" + cfg.GetString("app.env") + "/" + strconv.Itoa(int(time.Now().Unix()))
	resp, err := cfn.CreateInvalidation(ctx, &cloudfront.CreateInvalidationInput{
		DistributionId: aws.String(input.DistributionID),
		InvalidationBatch: &types.InvalidationBatch{
			CallerReference: aws.String(caller),
			Paths: &types.Paths{
				Items: []string{
					input.URI,
				},
				Quantity: aws.Int32(1),
			},
		},
	})
	if err != nil {
		return
	}
	input.SendMessage(ctx, fmt.Sprintf(":spin: invalidation submitted <https://console.aws.amazon.com/cloudfront/v4/home?region=us-east-1#/distributions/%s/invalidations/details/%s>", input.DistributionID, aws.ToString(resp.Invalidation.Id)))

	waiter := cloudfront.NewInvalidationCompletedWaiter(cfn)
	err = waiter.Wait(ctx, &cloudfront.GetInvalidationInput{
		DistributionId: aws.String(input.DistributionID),
		Id:             resp.Invalidation.Id,
	}, time.Minute*5)
	if err != nil {
		return
	}
	res = &InvalidationResponse{
		DistributionID: input.DistributionID,
		InvalidationID: aws.ToString(resp.Invalidation.Id),
	}
	//
	input.SendMessage(ctx, fmt.Sprintf(":greentick1: invalidation completed <https://console.aws.amazon.com/cloudfront/v4/home?region=us-east-1#/distributions/%s/invalidations/details/%s>", input.DistributionID, aws.ToString(resp.Invalidation.Id)))
	return

}
