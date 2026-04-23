package domainhealth

import (
	"context"
	"net"
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

type DomainHealthService struct{}

var service = &DomainHealthService{}

func DomainHealth() *DomainHealthService {
	return service
}

type DomainHealthResult struct {
	Domain      string   `json:"domain"`
	SPFStatus   string   `json:"spf_status"`
	DKIMStatus  string   `json:"dkim_status"`
	DMARCStatus string   `json:"dmarc_status"`
	MXStatus    string   `json:"mx_status"`
	Issues      []string `json:"issues"`
}

// DomainHealthRow is a typed struct for DB row serialization (fixes gdb.Record null JSON issue)
type DomainHealthRow struct {
	Id          int    `json:"id"`
	Domain      string `json:"domain"`
	SPFStatus   string `json:"spf_status"`
	DKIMStatus  string   `json:"dkim_status"`
	DMARCStatus string   `json:"dmarc_status"`
	MXStatus    string   `json:"mx_status"`
	SPFRecord   string `json:"spf_record"`
	DKIMRecord  string `json:"dkim_record"`
	DMARCRecord string `json:"dmarc_record"`
	MXRecords   string `json:"mx_records"`
	Issues      string `json:"issues"`
	LastChecked int64  `json:"last_checked"`
}

// CheckDomainHealth performs DNS lookups for SPF, DKIM, DMARC, MX
func (s *DomainHealthService) CheckDomainHealth(ctx context.Context, domain string) (*DomainHealthResult, error) {
	result := &DomainHealthResult{
		Domain: domain,
		Issues: []string{},
	}

	// SPF check
	txtRecords, _ := net.DefaultResolver.LookupTXT(ctx, domain)
	result.SPFStatus = "missing"
	spfRecord := ""
	for _, txt := range txtRecords {
		if strings.HasPrefix(txt, "v=spf1") {
			result.SPFStatus = "pass"
			spfRecord = txt
			break
		}
	}
	if result.SPFStatus == "missing" {
		result.Issues = append(result.Issues, "SPF record not found — emails may be rejected")
	}

	// DKIM check (default selector)
	dkimRecords, _ := net.DefaultResolver.LookupTXT(ctx, "default._domainkey."+domain)
	result.DKIMStatus = "missing"
	dkimRecord := ""
	for _, txt := range dkimRecords {
		if strings.Contains(txt, "v=DKIM1") || strings.Contains(txt, "p=") {
			result.DKIMStatus = "pass"
			dkimRecord = txt
			break
		}
	}
	if result.DKIMStatus == "missing" {
		result.Issues = append(result.Issues, "DKIM record not found — emails may land in spam")
	}

	// DMARC check
	dmarcRecords, _ := net.DefaultResolver.LookupTXT(ctx, "_dmarc."+domain)
	result.DMARCStatus = "missing"
	dmarcRecord := ""
	for _, txt := range dmarcRecords {
		if strings.HasPrefix(txt, "v=DMARC1") {
			result.DMARCStatus = "pass"
			dmarcRecord = txt
			break
		}
	}
	if result.DMARCStatus == "missing" {
		result.Issues = append(result.Issues, "DMARC record not found — deliverability risk")
	}

	// MX check
	mxRecords, _ := net.DefaultResolver.LookupMX(ctx, domain)
	result.MXStatus = "missing"
	mxList := []string{}
	if len(mxRecords) > 0 {
		result.MXStatus = "pass"
		for _, mx := range mxRecords {
			mxList = append(mxList, mx.Host)
		}
	}
	if result.MXStatus == "missing" {
		result.Issues = append(result.Issues, "MX records not found — cannot receive emails")
	}

	// Save to DB
	now := time.Now().Unix()
	issuesJSON := "[]"
	if len(result.Issues) > 0 {
		issuesJSON = `["` + strings.Join(result.Issues, `","`) + `"]`
	}

	g.DB().Model("bm_domain_health").Save(g.Map{
		"domain":       domain,
		"spf_status":   result.SPFStatus,
		"dkim_status":  result.DKIMStatus,
		"dmarc_status": result.DMARCStatus,
		"mx_status":    result.MXStatus,
		"spf_record":   spfRecord,
		"dkim_record":  dkimRecord,
		"dmarc_record": dmarcRecord,
		"mx_records":   strings.Join(mxList, ","),
		"issues":       issuesJSON,
		"last_checked": now,
	})

	return result, nil
}

// CheckAllDomains checks health for all configured domains
func (s *DomainHealthService) CheckAllDomains(ctx context.Context) ([]DomainHealthResult, error) {
	var domains []struct {
		ARecord string `json:"a_record"`
	}
	err := g.DB().Model("domain").Fields("a_record").Scan(&domains)
	if err != nil || len(domains) == 0 {
		return nil, err
	}

	results := []DomainHealthResult{}
	for _, d := range domains {
		domain := strings.TrimSuffix(d.ARecord, ".")
		if domain == "" {
			continue
		}
		result, err := s.CheckDomainHealth(ctx, domain)
		if err != nil {
			g.Log().Warningf(ctx, "Domain health check failed for %s: %v", domain, err)
			continue
		}
		results = append(results, *result)
	}

	g.Log().Infof(ctx, "Domain health check completed: %d domains", len(results))
	return results, nil
}

// GetDomainHealth returns cached health for a domain
func (s *DomainHealthService) GetDomainHealth(ctx context.Context, domain string) (*DomainHealthRow, error) {
	var row DomainHealthRow
	err := g.DB().Model("bm_domain_health").Where("domain = ?", domain).Scan(&row)
	if err != nil {
		return nil, err
	}
	if row.Id == 0 {
		return nil, nil
	}
	return &row, nil
}

// GetAllDomainHealth returns cached health for all domains
func (s *DomainHealthService) GetAllDomainHealth(ctx context.Context) ([]DomainHealthRow, error) {
	var rows []DomainHealthRow
	err := g.DB().Model("bm_domain_health").Order("domain").Scan(&rows)
	if err != nil {
		g.Log().Errorf(ctx, "DomainHealth query error: %v", err)
		return nil, err
	}
	g.Log().Infof(ctx, "DomainHealth rows: %d", len(rows))
	return rows, nil
}
