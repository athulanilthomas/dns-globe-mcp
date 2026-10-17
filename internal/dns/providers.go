package dns

type dohProvider struct {
	name     string
	lat, lng float64
	url      string
	headers  map[string]string
}

// lat and lng are dummy data

var providers = []dohProvider{
	{
		name: "Cloudflare",
		lat:  37.77,
		lng:  -10.42,
		url:  "https://cloudflare-dns.com/dns-query",
	},
	{
		name: "Google",
		lat:  67.42,
		lng:  107.08,
		url:  "https://dns.google/dns-query",
	},
	{
		name: "Quad9",
		lat:  27.37,
		lng:  67.55,
		url:  "https://dns.quad9.net/dns-query",
	},
	{
		name: "OpenDNS",
		lat:  57.55,
		lng:  -122.27,
		url:  "https://doh.opendns.com/dns-query",
	},
}
