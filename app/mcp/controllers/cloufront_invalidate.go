package controllers

import (
	"context"
	"github.com/desygner/desygner-toolchain/app/handlers"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"net/url"
)

func NewCloudfrontServerTool(handler *handlers.Invalidate) (*mcp.Tool, mcp.ToolHandlerFor[CloudfrontInvalidateArguments, any]) {
	cfInvalidate := NewCloudfrontInvalidate(handler)
	return &mcp.Tool{
		Name:        "cloudfront-invalidate",
		Description: "Trigger Cloudfront invalidation based on the given URL. no need to trigger other tools",
	}, cfInvalidate.Handle
}

type CloudfrontInvalidate struct {
	handler *handlers.Invalidate
}

type SlackArgument struct {
	Channel string `json:"channel" mcp:"Slack channel to post the message"`
	TS      string `json:"ts" mcp:"Slack timestamp to post the message"`
}

type CloudfrontInvalidateArguments struct {
	URL   string         `json:"url" mcp:"URL to invalidate"`
	Slack *SlackArgument `json:"slack" mcp:"Slack definition"`
}

func NewCloudfrontInvalidate(handler *handlers.Invalidate) *CloudfrontInvalidate {
	return &CloudfrontInvalidate{
		handler: handler,
	}
}

func (c *CloudfrontInvalidate) Handle(ctx context.Context, req *mcp.CallToolRequest, params CloudfrontInvalidateArguments) (resp *mcp.CallToolResult, response any, err error) {
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
	progressToken := req.Params.GetProgressToken()
	req.Session.NotifyProgress(ctx, &mcp.ProgressNotificationParams{
		Total:         3,
		Progress:      0,
		ProgressToken: progressToken,
		Message:       "request queued. please wait",
	})
	input := &handlers.InvalidationInput{}
	if params.Slack != nil {
		input.Input.SlackChannel = params.Slack.Channel
		input.Input.SlackTS = params.Slack.TS
	}
	req.Session.NotifyProgress(ctx, &mcp.ProgressNotificationParams{
		Total:         3,
		Progress:      1,
		ProgressToken: progressToken,
		Message:       "resolving distribution",
	})
	distID, err := c.handler.ResolveDistribution(context.Background(), params.URL)
	if err != nil {
		return
	}
	req.Session.NotifyProgress(ctx, &mcp.ProgressNotificationParams{
		Total:         3,
		Progress:      2,
		ProgressToken: progressToken,
		Message:       "running invalidation",
	})
	parsedURL, _ := url.Parse(params.URL)
	res, err := c.handler.Handle(context.Background(), &handlers.InvalidationInput{
		DistributionID: distID.DistributionID,
		Input:          input.Input,
		URI:            parsedURL.EscapedPath(),
	})
	if err == nil {
		resp = &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "invalidation done. " + res.InvalidationID},
			},
		}
		return
	}
	return
}
