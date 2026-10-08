package dns

type dohProvider struct {
	name     string
	lat, lng float64
	url      string
	headers  map[string]string
}

var providers = []dohProvider{
	{
		name: "Cloudflare",
		lat:  37.77,
		lng:  -122.42,
		url:  "https://cloudflare-dns.com/dns-query",
	},
	{
		name: "Google",
		lat:  37.42,
		lng:  -122.08,
		url:  "https://dns.google/dns-query",
	},
	{
		name: "Quad9",
		lat:  47.37,
		lng:  8.55,
		url:  "https://dns.quad9.net/dns-query",
	},
	{
		name: "OpenDNS",
		lat:  37.55,
		lng:  -122.27,
		url:  "https://doh.opendns.com/dns-query",
	},
}
