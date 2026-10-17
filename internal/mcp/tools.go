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
type LookupDnsRecordsTool struct{}

type LookupDnsRecordsInput struct {
	Domain     string `json:"domain" jsonschema:"the domain name to check, e.g. example.com"`
	RecordType string `json:"recordType" jsonschema:"DNS record type, e.g. A, AAAA, MX, TXT"`
}

type LookupDnsRecordsOutput struct {
	Domain     string             `json:"domain"`
	RecordType string             `json:"recordType"`
	Results    []dns.RegionResult `json:"results"`
}

func (t *LookupDnsRecordsTool) Meta() *mcp.Tool {
	return &mcp.Tool{
		Name:        "lookup_dns_records",
		Description: "Queries public DNS resolvers (Cloudflare, Google, Quad9, OpenDNS) over DNS-over-HTTPS for a domain and record type, and reports each resolver's status (resolved, pending, error, nxdomain) and returned records. Use it to check whether a DNS change has propagated. Renders an interactive globe of the results.",
		Meta: mcp.Meta{
			"ui":             map[string]any{"resourceUri": WidgetURI},
			"ui/resourceUri": WidgetURI,
		},
	}
}

func (t *LookupDnsRecordsTool) Implementation(ctx context.Context, req *mcp.CallToolRequest, input LookupDnsRecordsInput) (
	*mcp.CallToolResult,
	LookupDnsRecordsOutput,
	error,
) {
	results, err := dns.LookupRecords(ctx, input.Domain, input.RecordType)
	if err != nil {
		return nil, LookupDnsRecordsOutput{}, err
	}

	output := LookupDnsRecordsOutput{
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

var _ Tool[LookupDnsRecordsInput, LookupDnsRecordsOutput] = (*LookupDnsRecordsTool)(nil)
