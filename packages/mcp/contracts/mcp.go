package contracts

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type Mcp interface {
	Server() *mcp.Server
	Transport() *mcp.StreamableServerTransport
}
