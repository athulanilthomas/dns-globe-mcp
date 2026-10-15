export const SAMPLE = {
  domain: "example.com",
  recordType: "A",
  results: [
    { resolver: "Cloudflare", lat: 51.51, lng: -0.13, status: "resolved", records: ["93.184.216.34"], error: "" },
    { resolver: "Google", lat: 37.42, lng: -122.08, status: "resolved", records: ["93.184.216.34"], error: "" },
    { resolver: "Quad9", lat: 47.37, lng: 8.54, status: "stale", records: ["93.184.216.33"], error: "" },
    { resolver: "OpenDNS", lat: 1.35, lng: 103.82, status: "pending", records: null, error: "" },
  ],
};
