package main

import (
	"log"
	"net/http"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	internalconfig "github.com/athulanilthomas/dns-globe-mcp/internal/config"
	internalmcp "github.com/athulanilthomas/dns-globe-mcp/internal/mcp"
	internalweb "github.com/athulanilthomas/dns-globe-mcp/internal/web"
)

func main() {
	config := internalconfig.NewConfig()
	server := internalmcp.NewServer(config)

	if server == nil {
		log.Fatal("failed to construct MCP server")
	}

	dnsTool := internalmcp.DnsPropagationTool{}
	mcp.AddTool(server, dnsTool.Meta(), dnsTool.Implementation)
	internalmcp.AddWidgetResource(server)

	mux := http.NewServeMux()
	mux.Handle("/", internalweb.Handler())
	mux.Handle("/mcp", internalmcp.Handler(server))

	log.Println("MCP server listening on", config.Port, "at /mcp")

	if err := http.ListenAndServe(config.Port, mux); err != nil {
		log.Fatal("server failed:", err)
	}
}
