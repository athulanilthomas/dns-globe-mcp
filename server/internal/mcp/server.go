package mcp

import (
	"github.com/athulanilthomas/dns-globe-mcp/server/internal/config"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func NewServer(cfg *config.Config) *mcp.Server {
	impl := &mcp.Implementation{
		Name:    cfg.MCP.Name,
		Version: cfg.MCP.Version,
	}

	server := mcp.NewServer(impl, nil)

	return server
}
