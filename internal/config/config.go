package config

import "os"

type MCPConfig struct {
	Name    string
	Version string
}

type Config struct {
	Port string
	MCP  MCPConfig
}

func NewConfig() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return &Config{
		Port: ":" + port,
		MCP: MCPConfig{
			Name:    "dns-globe-mcp",
			Version: "v0.1.0",
		},
	}
}
