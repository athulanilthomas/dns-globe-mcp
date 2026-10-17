package mcp

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

//go:embed dist/index.html
var widgetHTML string

const (
	WidgetURI       = "ui://dns-globe/widget"
	dataPlaceholder = `"__DNS_DATA__"`
)

func RenderWidget(data any) (string, error) {
	if !strings.Contains(widgetHTML, dataPlaceholder) {
		return "", fmt.Errorf("widget HTML is missing the %s placeholder; rebuild the widget", dataPlaceholder)
	}
	payload, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("marshal widget data: %w", err)
	}
	return strings.Replace(widgetHTML, dataPlaceholder, string(payload), 1), nil
}

func UIResourceContent(uri string, data any) (*mcp.EmbeddedResource, error) {
	html, err := RenderWidget(data)

	if err != nil {
		return nil, err
	}

	resource := mcp.EmbeddedResource{Resource: &mcp.ResourceContents{
		URI:      uri,
		Text:     html,
		MIMEType: "text/html",
	}}

	return &resource, nil
}
