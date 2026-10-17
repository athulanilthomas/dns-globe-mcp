package mcp

import (
	"context"
	"fmt"

	"github.com/athulanilthomas/dns-globe-mcp/internal/dns"
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
	Domain     string             `json:"domain"`
	RecordType string             `json:"recordType"`
	Results    []dns.RegionResult `json:"results"`
}

func (t *DnsPropagationTool) Meta() *mcp.Tool {
	return &mcp.Tool{
		Name:        "check_dns_propagation",
		Description: "Queries public DNS resolvers (Cloudflare, Google, Quad9, OpenDNS) over DNS-over-HTTPS for a domain and record type, and reports each resolver's status (resolved, stale, pending, error, nxdomain) and returned records. Use it to check whether a DNS change has propagated. Renders an interactive globe of the results.",
		Meta: mcp.Meta{
			"ui":             map[string]any{"resourceUri": WidgetURI},
			"ui/resourceUri": WidgetURI,
		},
	}
}

func (t *DnsPropagationTool) Implementation(ctx context.Context, req *mcp.CallToolRequest, input DnsPropagationInput) (
	*mcp.CallToolResult,
	DnsPropagationOutput,
	error,
) {
	results, err := dns.CheckDNSPropagation(ctx, input.Domain, input.RecordType)
	if err != nil {
		return nil, DnsPropagationOutput{}, err
	}

	output := DnsPropagationOutput{
		Domain:     input.Domain,
		RecordType: input.RecordType,
		Results:    results,
	}

	summary := fmt.Sprintf("Checked %s %s records across %d resolvers.", input.Domain, input.RecordType, len(results))
	result := &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: summary}},
	}

	return result, output, nil
}

var _ Tool[DnsPropagationInput, DnsPropagationOutput] = (*DnsPropagationTool)(nil)
