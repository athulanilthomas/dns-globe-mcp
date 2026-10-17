package mcp

import (
	"context"
	_ "embed"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

//go:embed dist/index.html
var widgetHTML string

const (
	WidgetURI      = "ui://dns-globe/widget"
	WidgetMIMEType = "text/html;profile=mcp-app"
)

func AddWidgetResource(server *mcp.Server) {
	server.AddResource(&mcp.Resource{
		Name:        "dns-globe-widget",
		Title:       "DNS Globe",
		Description: "Interactive globe showing DNS propagation results",
		URI:         WidgetURI,
		MIMEType:    WidgetMIMEType,
	}, func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		return &mcp.ReadResourceResult{
			Contents: []*mcp.ResourceContents{{
				URI:      WidgetURI,
				MIMEType: WidgetMIMEType,
				Text:     widgetHTML,
			}},
		}, nil
	})
}
