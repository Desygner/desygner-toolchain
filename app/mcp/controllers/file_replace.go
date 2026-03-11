package controllers

import (
	"context"
	"github.com/desygner/desygner-toolchain/app/handlers"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func NewFileReplace(handler *handlers.Replace) (*mcp.Tool, mcp.ToolHandlerFor[FileReplaceArguments, any]) {
	cfInvalidate := NewFileReplaceHandler(handler)
	return &mcp.Tool{
		Name:        "file-replace",
		Description: "Replace object from slack private file to s3 desygner-owned buckets. no need to trigger other tools",
	}, cfInvalidate.Handle
}

type FileReplace struct {
	handler *handlers.Replace
}

type FileReplaceArguments struct {
	Source string         `json:"source" mcp:"only slack private url allowed"`
	Target string         `json:"target" mcp:"target location the file will be replacing. it can be s3 url or cloudfront url"`
	Slack  *SlackArgument `json:"slack" mcp:"Slack definition"`
}

func NewFileReplaceHandler(handler *handlers.Replace) *FileReplace {
	return &FileReplace{
		handler: handler,
	}
}

func (c *FileReplace) Handle(ctx context.Context, req *mcp.CallToolRequest, params FileReplaceArguments) (resp *mcp.CallToolResult, response any, err error) {
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
	input := &handlers.ReplaceInput{
		Source: params.Source,
		Target: params.Target,
	}
	if params.Slack != nil && params.Slack.Channel != "" && params.Slack.TS != "" {
		input.Input.SlackChannel = params.Slack.Channel
		input.Input.SlackTS = params.Slack.TS
	}
	err = c.handler.Handle(context.Background(), input)
	if err == nil {
		resp = &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "replacement done"},
			},
		}
		return
	}
	return
}
