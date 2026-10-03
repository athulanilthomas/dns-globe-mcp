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

	dnsTool := internalmcp.DnsPropagationTool{}
	mcp.AddTool(server, dnsTool.Meta(), dnsTool.Implementation)

	mux := http.NewServeMux()
	mux.Handle("/mcp", internalmcp.Handler(server))

	antiCSRF := http.NewCrossOriginProtection()
	antiCSRF.AddTrustedOrigin("http://localhost:5173")

	log.Println("MCP server listening on", config.Port, "at /mcp")

	if err := http.ListenAndServe(config.Port, corsMiddleware(antiCSRF.Handler(mux))); err != nil {
		log.Fatal("server failed:", err)
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, mcp-protocol-version")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
