package dns

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/dns/dnsmessage"
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

func mapRecordType(recordType string) (dnsmessage.Type, error) {
	switch strings.ToUpper(strings.TrimSpace(recordType)) {
	case "A":
		return dnsmessage.TypeA, nil
	case "AAAA":
		return dnsmessage.TypeAAAA, nil
	case "CNAME":
		return dnsmessage.TypeCNAME, nil
	case "MX":
		return dnsmessage.TypeMX, nil
	case "NS":
		return dnsmessage.TypeNS, nil
	case "TXT":
		return dnsmessage.TypeTXT, nil
	case "SOA":
		return dnsmessage.TypeSOA, nil
	case "PTR":
		return dnsmessage.TypePTR, nil
	case "SRV":
		return dnsmessage.TypeSRV, nil
	default:
		return 0, fmt.Errorf("unsupported DNS record type: %q", recordType)
	}
}

func assertAnswer(rr *dnsmessage.Resource) string {
	var record string

	switch body := rr.Body.(type) {
	case *dnsmessage.AResource:
		record = net.IP(body.A[:]).String()
	case *dnsmessage.AAAAResource:
		record = net.IP(body.AAAA[:]).String()
	case *dnsmessage.CNAMEResource:
		record = body.CNAME.String()
	case *dnsmessage.MXResource:
		record = fmt.Sprintf("%d %s", body.Pref, body.MX.String())
	case *dnsmessage.TXTResource:
		record = strings.Join(body.TXT, " ")
	case *dnsmessage.NSResource:
		record = body.NS.String()
	case *dnsmessage.SOAResource:
		record = fmt.Sprintf("%s %s", body.MBox.String(), body.NS.String())

	default:
		record = rr.GoString()
	}

	return record
}

var client = &http.Client{Timeout: 10 * time.Second}

func appendResolverResult(res *RegionResult, results *ResolverResults) {
	results.mu.Lock()
	results.results = append(results.results, *res)
	results.mu.Unlock()
}

func buildQuery(domain string, recordType dnsmessage.Type) ([]byte, error) {
	name, err := dnsmessage.NewName(domain + ".")
	if err != nil {
		return nil, err
	}

	msg := dnsmessage.Message{
		ID: 0, RecursionDesired: true,
		Questions: []dnsmessage.Question{
			{Name: name, Type: recordType, Class: dnsmessage.ClassINET},
		},
	}

	return msg.Pack()
}

func resolveDNS(reqParams DoHRequest, results *ResolverResults) {
	base := RegionResult{
		Resolver: reqParams.provider.name,
		Lat:      reqParams.provider.lat,
		Lng:      reqParams.provider.lng,
	}

	recordType, err := mapRecordType(reqParams.recordType)
	if err != nil {
		base.Status = "pending"
		appendResolverResult(&base, results)
		return
	}

	query, err := buildQuery(reqParams.domain, recordType)
	if err != nil {
		base.Status = "pending"
		appendResolverResult(&base, results)
		return
	}

	req, err := http.NewRequest(http.MethodPost, reqParams.provider.url, bytes.NewReader(query))
	req.Header.Add("accept", "application/dns-message")
	req.Header.Add("content-type", "application/dns-message")

	res, err := client.Do(req)
	if err != nil {
		base.Status = "pending"
		appendResolverResult(&base, results)
		return
	}

	defer res.Body.Close()

	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		base.Status = "pending"
		appendResolverResult(&base, results)
		return
	}

	msg := new(dnsmessage.Message)
	msg.Unpack(resBytes)

	records := make([]string, 0, len(msg.Answers))
	for _, rec := range msg.Answers {
		records = append(records, assertAnswer(&rec))
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
