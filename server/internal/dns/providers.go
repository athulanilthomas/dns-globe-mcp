package dns

import (
	"fmt"
	"net/url"
)

var providers = []dohProvider{
	{
		name: "Cloudflare",
		lat:  37.77,
		lng:  -122.42,
		url: func(d, t string) string {
			return fmt.Sprintf("https://cloudflare-dns.com/dns-query?name=%s&type=%s", url.QueryEscape(d), t)
		},
		headers: map[string]string{"Accept": "application/dns-json"},
	},
	{
		name: "Google",
		lat:  37.42,
		lng:  -122.08,
		url: func(d, t string) string {
			return fmt.Sprintf("https://dns.google/resolve?name=%s&type=%s", url.QueryEscape(d), t)
		},
	},
	{
		name: "Quad9",
		lat:  47.37,
		lng:  8.55,
		url: func(d, t string) string {
			return fmt.Sprintf("https://dns.quad9.net/dns-query?name=%s&type=%s", url.QueryEscape(d), t)
		},
		headers: map[string]string{"Accept": "application/dns-json"},
	},
	{
		name: "OpenDNS",
		lat:  37.55,
		lng:  -122.27,
		url: func(d, t string) string {
			return fmt.Sprintf("https://doh.opendns.com/dns-query?name=%s&type=%s", url.QueryEscape(d), t)
		},
		headers: map[string]string{"Accept": "application/dns-json"},
	},
}
