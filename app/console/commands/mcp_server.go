package commands

import (
	"fmt"
	"github.com/desygner/desygner-toolchain/packages/mcp/facades"
	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	facades2 "github.com/goravel/framework/facades"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"net/http"
)

type McpServer struct {
}

// Signature The name and signature of the console command.
func (receiver *McpServer) Signature() string {
	return "mcp:start"
}

// Description The console command description.
func (receiver *McpServer) Description() string {
	return "Start stdio server for mcp"
}

// Extend The console command extend.
func (receiver *McpServer) Extend() command.Extend {
	return command.Extend{}
}

// Handle Execute the console command.
func (receiver *McpServer) Handle(ctx console.Context) error {
	srv := facades.Mcp().Server()
	handler := mcp.NewStreamableHTTPHandler(func(request *http.Request) *mcp.Server {
		return srv
	}, nil)
	mcpPort := facades2.Config().GetInt("mcp.port")
	http.HandleFunc("/mcp", handler.ServeHTTP)
	http.HandleFunc("/mcp/healthcheck", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(200)
		writer.Write([]byte("ok"))
	})
	return http.ListenAndServe(fmt.Sprintf(":%d", mcpPort), nil)
}
