package dns

type RegionResult struct {
	Resolver string   `json:"resolver"`
	Lat      float64  `json:"lat"`
	Lng      float64  `json:"lng"`
	Status   string   `json:"status"` // "resolved" | "stale" | "pending"
	Records  []string `json:"records"`
}
