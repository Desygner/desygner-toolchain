package controllers

import (
	"context"
	"encoding/json"
	"github.com/desygner/desygner-toolchain/app/handlers"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func NewBucketResolver(handler *handlers.Invalidate) (*mcp.Tool, mcp.ToolHandlerFor[CloudfrontInvalidateArguments, any]) {
	cfInvalidate := NewBucketResolverHandler(handler)
	return &mcp.Tool{
		Name:        "bucket-resolver",
		Description: "Resolve the source bucket of the CloudFront URL",
	}, cfInvalidate.Handle
}

type BucketResolver struct {
	handler *handlers.Invalidate
}

func NewBucketResolverHandler(handler *handlers.Invalidate) *BucketResolver {
	return &BucketResolver{
		handler: handler,
	}
}

func (c *BucketResolver) Handle(ctx context.Context, req *mcp.CallToolRequest, params CloudfrontInvalidateArguments) (resp *mcp.CallToolResult, response any, err error) {
	defer func() {
		if err != nil {
			req.Session.NotifyProgress(ctx, &mcp.ProgressNotificationParams{
				Message:  "request failed. error: " + err.Error(),
				Total:    3,
				Progress: 3,
			})
			resp = &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: "request failed. error: " + err.Error()},
				},
			}
		}
	}()
	input := &handlers.InvalidationInput{}
	if params.Slack != nil && params.Slack.Channel != "" && params.Slack.TS != "" {
		input.Input.SlackChannel = params.Slack.Channel
		input.Input.SlackTS = params.Slack.TS
	}
	bucketResult, err := c.handler.ResolveBucket(context.Background(), params.URL)
	if err == nil {
		enc, _ := json.MarshalIndent(bucketResult, "", "  ")
		resp = &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: string(enc)},
			},
		}
		return
	}
	return
}
