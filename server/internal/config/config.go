package config

type MCPConfig struct {
	Name    string
	Version string
}

type Config struct {
	Port string
	MCP  MCPConfig
}

func NewConfig() *Config {
	return &Config{
		Port: ":8080",
		MCP: MCPConfig{
			Name:    "dns-globe-mcp",
			Version: "v0.1.0",
		},
	}
}
