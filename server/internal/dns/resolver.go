package dns

func CheckPropagation(domain string, recordType string) ([]RegionResult, error) {
	return []RegionResult{
		{Region: "Mumbai", Lat: 19.07, Lng: 72.87, Status: "resolved", Records: []string{"93.184.216.34"}},
		{Region: "Singapore", Lat: 1.35, Lng: 103.81, Status: "resolved", Records: []string{"93.184.216.34"}},
		{Region: "Frankfurt", Lat: 50.11, Lng: 8.68, Status: "stale", Records: []string{"93.184.216.33"}},
		{Region: "Virginia", Lat: 38.95, Lng: -77.45, Status: "pending", Records: nil},
	}, nil
}
