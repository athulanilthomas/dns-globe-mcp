package mcp

import "github.com/modelcontextprotocol/go-sdk/mcp"

func UIResourceContent(uri, widgetURL string) *mcp.EmbeddedResource {
	return &mcp.EmbeddedResource{
		Resource: &mcp.ResourceContents{
			URI:      uri,
			MIMEType: "text/uri-list",
			Text:     widgetURL,
		},
	}
}
