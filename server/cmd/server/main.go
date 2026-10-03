package main

import (
	"log"
	"net/http"

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

	mux := http.NewServeMux()
	mux.Handle("/mcp", internalmcp.Handler(server))

	log.Println("MCP server listening on", config.Port, "at /mcp")

	if err := http.ListenAndServe(config.Port, mux); err != nil {
		log.Fatal("server failed:", err)
	}
}
