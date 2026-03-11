package service

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type Mcp struct {
	srv       *mcp.Server
	transport *mcp.StreamableServerTransport
}

func (m *Mcp) Server() *mcp.Server {
	return m.srv
}

func (m *Mcp) Transport() *mcp.StreamableServerTransport {
	return m.transport
}

func NewMcp(srv *mcp.Server, transport *mcp.StreamableServerTransport) *Mcp {
	return &Mcp{
		srv:       srv,
		transport: transport,
	}
}
