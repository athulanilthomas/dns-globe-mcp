package dns

type RegionResult struct {
	Region  string   `json:"region"`
	Lat     float64  `json:"lat"`
	Lng     float64  `json:"lng"`
	Status  string   `json:"status"` // "resolved" | "stale" | "pending"
	Records []string `json:"records"`
}
