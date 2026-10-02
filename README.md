# DNS Globe MCP

A small experimental Model Context Protocol (MCP) project for DNS-related tooling and browser integrations.

This repository currently contains:
- a Go-based MCP server under `server/`
- a minimal `ping` tool exposed by the server
- a browser widget scaffold under `widget/`

## Project layout

- `server/` — Go server code and MCP setup
- `server/cmd/server` — entry point for the MCP server
- `server/internal/` — configuration, DNS logic, and MCP handlers
- `widget/` — static front-end assets for experimenting with the MCP client

## Run the server

```bash
cd server
go run ./cmd/server
```