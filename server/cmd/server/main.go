package main

import (
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	internalconfig "github.com/athulanilthomas/dns-globe-mcp/server/internal/config"
	internalmcp "github.com/athulanilthomas/dns-globe-mcp/server/internal/mcp"
)

func main() {
	config := internalconfig.NewConfig()
	server := internalmcp.NewServer(config)

	if server == nil {
		log.Fatal("failed to construct MCP server")
	}

	pingTool := internalmcp.PingTool{}
	mcp.AddTool(server, pingTool.Meta(), pingTool.Implementation)

	log.Println("MCP server instance created:", "dns-globe-mcp v0.1.0")
}
