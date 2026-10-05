export const SAMPLE = {
  domain: "example.com",
  recordType: "A",
  results: [
    { region: "Mumbai", lat: 19.07, lng: 72.87, status: "resolved", records: ["93.184.216.34"] },
    { region: "Singapore", lat: 1.35, lng: 103.81, status: "resolved", records: ["93.184.216.34"] },
    { region: "Frankfurt", lat: 50.11, lng: 8.68, status: "stale", records: ["93.184.216.33"] },
    { region: "Virginia", lat: 38.95, lng: -77.45, status: "pending", records: null },
  ],
};
