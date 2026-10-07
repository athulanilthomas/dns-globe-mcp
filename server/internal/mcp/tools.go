package mcp

import (
	"context"

	"github.com/athulanilthomas/dns-globe-mcp/server/internal/dns"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type Tool[I any, O any] interface {
	Meta() *mcp.Tool
	Implementation(ctx context.Context, req *mcp.CallToolRequest, input I) (*mcp.CallToolResult, O, error)
}
type DnsPropagationTool struct{}

type DnsPropagationInput struct {
	Domain     string `json:"domain" jsonschema:"the domain name to check, e.g. example.com"`
	RecordType string `json:"recordType" jsonschema:"DNS record type, e.g. A, AAAA, MX, TXT"`
}

type DnsPropagationOutput struct {
	Results []dns.RegionResult `json:"results"`
}

func (t *DnsPropagationTool) Meta() *mcp.Tool {
	return &mcp.Tool{
		Name:        "check_dns_propagation",
		Description: "Checks DNS propagation for a domain across multiple global regions",
	}
}

func (t *DnsPropagationTool) Implementation(ctx context.Context, req *mcp.CallToolRequest, input DnsPropagationInput) (
	*mcp.CallToolResult,
	DnsPropagationOutput,
	error,
) {
	results, err := dns.CheckDNSPropagation(input.Domain, input.RecordType)
	if err != nil {
		return nil, DnsPropagationOutput{}, err
	}

	result := &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: "DNS propagation checked across regions. See embedded widget for details."},
			UIResourceContent("ui://dns-globe/widget", "http://localhost:5172/index.html"),
		},
	}

	return result, DnsPropagationOutput{Results: results}, nil
}

var _ Tool[DnsPropagationInput, DnsPropagationOutput] = (*DnsPropagationTool)(nil)
