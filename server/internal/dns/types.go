package dns

type status string

const (
	StatusResolved status = "resolved"
	StatusStale    status = "stale"
	StatusPending  status = "pending"
	StatusError    status = "error"
)

type RegionResult struct {
	Resolver string   `json:"resolver"`
	Lat      float64  `json:"lat"`
	Lng      float64  `json:"lng"`
	Status   status   `json:"status"`
	Records  []string `json:"records"`
	Error    string   `json:"error"`
}
