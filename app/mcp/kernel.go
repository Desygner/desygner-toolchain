package mcp

import (
	"github.com/desygner/desygner-toolchain/app/handlers"
	"github.com/desygner/desygner-toolchain/app/mcp/controllers"
	"github.com/desygner/desygner-toolchain/packages/mcp/facades"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func Tools() {
	srv := facades.Mcp().Server()
	// register tool
	sourceInvalidateHandler := &handlers.Invalidate{}
	invalidationTool, invalidationHandler := controllers.NewCloudfrontServerTool(sourceInvalidateHandler)
	bucketTool, bucketHandler := controllers.NewBucketResolver(sourceInvalidateHandler)
	fileReplaceTool, fileReplaceHandler := controllers.NewFileReplace(&handlers.Replace{})

	mcp.AddTool(srv, invalidationTool, invalidationHandler)
	mcp.AddTool(srv, bucketTool, bucketHandler)
	mcp.AddTool(srv, fileReplaceTool, fileReplaceHandler)
}
