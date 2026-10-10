package dns

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/dns/dnsmessage"
)

type DoHRequest struct {
	provider *dohProvider
	query    []byte
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
	default:
		return 0, fmt.Errorf("unsupported DNS record type: %q", recordType)
	}
}

func formatRecord(rr *dnsmessage.Resource) string {
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

var client = &http.Client{}

func buildQuery(domain string, recordType dnsmessage.Type) ([]byte, error) {
	if !strings.HasSuffix(domain, ".") {
		domain += "."
	}

	domain = strings.ToLower(strings.TrimSpace(domain))

	name, err := dnsmessage.NewName(domain)
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

func resolveDNS(ctx context.Context, reqParams DoHRequest) RegionResult {
	base := RegionResult{
		Resolver: reqParams.provider.name,
		Lat:      reqParams.provider.lat,
		Lng:      reqParams.provider.lng,
	}

	reqCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(reqCtx, http.MethodPost, reqParams.provider.url, bytes.NewReader(reqParams.query))
	if err != nil {
		return fail(base, err)
	}

	req.Header.Add("accept", "application/dns-message")
	req.Header.Add("content-type", "application/dns-message")

	res, err := client.Do(req)
	if err != nil || res.StatusCode >= http.StatusBadRequest {
		return fail(base, err)
	}

	defer res.Body.Close()

	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return fail(base, err)
	}

	msg := new(dnsmessage.Message)

	unpackError := msg.Unpack(resBytes)
	if unpackError != nil {
		return fail(base, err)
	}

	switch msg.RCode {
	case dnsmessage.RCodeSuccess:
	case dnsmessage.RCodeNameError:
		base.Status = StatusNXDomain
		return base
	default:
		return fail(base, fmt.Errorf("DNS error: %v", msg.RCode))
	}

	records := make([]string, 0, len(msg.Answers))
	for _, rec := range msg.Answers {
		records = append(records, formatRecord(&rec))
	}

	if len(records) == 0 {
		base.Status = StatusPending
		return base
	}

	base.Status = StatusResolved
	base.Records = records

	return base
}

func fail(r RegionResult, err error) RegionResult {
	r.Status = StatusError
	r.Error = err.Error()
	return r
}

func CheckDNSPropagation(ctx context.Context, domain string, recordType string) ([]RegionResult, error) {
	wg := &sync.WaitGroup{}

	results := make([]RegionResult, len(providers))

	mappedType, err := mapRecordType(recordType)
	if err != nil {
		return nil, err
	}

	query, err := buildQuery(domain, mappedType)
	if err != nil {
		return nil, err
	}

	for i, r := range providers {
		reqParams := DoHRequest{provider: &r, query: query}
		wg.Go(func() {
			results[i] = resolveDNS(ctx, reqParams)
		})
	}

	wg.Wait()

	return results, nil
}
