package mcp

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type Tool[I any, O any] interface {
	Meta() *mcp.Tool
	Implementation(ctx context.Context, req *mcp.CallToolRequest, input I) (*mcp.CallToolResult, O, error)
}

type PingTool struct{}

type PingToolInput struct {
	Name string `json:"name" jsonschema:"the name of the person to greet"`
}

type PingToolOutput struct {
	Greeting string `json:"greeting" jsonschema:"the greeting to tell to the user"`
}

func (ping *PingTool) Meta() *mcp.Tool {
	return &mcp.Tool{
		Name:        "ping",
		Description: "Ping Tool",
	}
}

func (ping *PingTool) Implementation(ctx context.Context, req *mcp.CallToolRequest, input PingToolInput) (
	*mcp.CallToolResult,
	PingToolOutput,
	error,
) {
	return nil, PingToolOutput{Greeting: "Hi " + input.Name}, nil
}

// Comil time check
var _ Tool[PingToolInput, PingToolOutput] = (*PingTool)(nil)
