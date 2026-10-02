package config

type Config struct {
	Name    string
	Version string
}

func NewConfig() Config {
	return Config{
		Name:    "dns-globe-mcp",
		Version: "v0.1.0",
	}
}
