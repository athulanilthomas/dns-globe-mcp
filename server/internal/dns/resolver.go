package dns

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

type ResolverResults struct {
	results []RegionResult
	mu      sync.RWMutex
}

type DoHRequest struct {
	provider   *dohProvider
	domain     string
	recordType string
}

type dohResponse struct {
	Status int `json:"Status"`
	Answer []struct {
		Data string `json:"data"`
	} `json:"Answer"`
}

var client = &http.Client{Timeout: 10 * time.Second}

func appendResolverResult(res *RegionResult, results *ResolverResults) {
	results.mu.Lock()
	results.results = append(results.results, *res)
	results.mu.Unlock()
}

func resolveDNS(reqParams DoHRequest, results *ResolverResults) {
	base := RegionResult{Resolver: reqParams.provider.name, Lat: reqParams.provider.lat, Lng: reqParams.provider.lng}

	req, err := http.NewRequest("GET", reqParams.provider.url(reqParams.domain, reqParams.recordType), nil)
	if err != nil {
		base.Status = "pending"
		appendResolverResult(&base, results)
		return
	}

	for k, v := range reqParams.provider.headers {
		req.Header.Add(k, v)
	}

	res, err := client.Do(req)
	if err != nil {
		base.Status = "pending"
		appendResolverResult(&base, results)
		return
	}

	defer res.Body.Close()

	if res.StatusCode >= http.StatusBadRequest {
		base.Status = "pending"
		appendResolverResult(&base, results)
		return
	}

	var parsed dohResponse
	dec := json.NewDecoder(res.Body)

	if err := dec.Decode(&parsed); err != nil {
		base.Status = "stale"
		appendResolverResult(&base, results)
		return
	}

	records := make([]string, 0, len(parsed.Answer))
	for _, rec := range parsed.Answer {
		records = append(records, rec.Data)
	}

	if len(records) == 0 {
		base.Status = "pending"
		appendResolverResult(&base, results)
		return
	}

	base.Status = "resolved"
	base.Records = records

	appendResolverResult(&base, results)
}

func CheckDNSPropagation(domain string, recordType string) ([]RegionResult, error) {
	wg := &sync.WaitGroup{}

	results := ResolverResults{
		results: make([]RegionResult, 0, len(providers)),
	}

	for _, r := range providers {
		reqParams := DoHRequest{provider: &r, domain: domain, recordType: recordType}
		wg.Go(func() {
			resolveDNS(reqParams, &results)
		})
	}

	wg.Wait()

	return results.results, nil
}
