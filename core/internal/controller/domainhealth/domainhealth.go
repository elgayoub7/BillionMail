package domainhealth

import (
	"context"
	"encoding/json"

	v1 "billionmail-core/api/domainhealth/v1"
	"billionmail-core/internal/service/domainhealth"
)

type ControllerV1 struct{}

func NewV1() *ControllerV1 {
	return &ControllerV1{}
}

// domainHealthAPIResponse is the API-safe version with Issues as []string instead of raw JSON string
type domainHealthAPIResponse struct {
	Id          int      `json:"id"`
	Domain      string   `json:"domain"`
	SPFStatus   string   `json:"spf_status"`
	DKIMStatus  string   `json:"dkim_status"`
	DMARCStatus string   `json:"dmarc_status"`
	MXStatus    string   `json:"mx_status"`
	SPFRecord   string   `json:"spf_record"`
	DKIMRecord  string   `json:"dkim_record"`
	DMARCRecord string   `json:"dmarc_record"`
	MXRecords   string   `json:"mx_records"`
	Issues      []string `json:"issues"`
	LastChecked int64    `json:"last_checked"`
}

// convertRow parses the DB row's Issues JSON string into a proper []string for the API
func convertRow(row *domainhealth.DomainHealthRow) domainHealthAPIResponse {
	resp := domainHealthAPIResponse{
		Id:          row.Id,
		Domain:      row.Domain,
		SPFStatus:   row.SPFStatus,
		DKIMStatus:  row.DKIMStatus,
		DMARCStatus: row.DMARCStatus,
		MXStatus:    row.MXStatus,
		SPFRecord:   row.SPFRecord,
		DKIMRecord:  row.DKIMRecord,
		DMARCRecord: row.DMARCRecord,
		MXRecords:   row.MXRecords,
		Issues:      []string{},
		LastChecked: row.LastChecked,
	}

	if row.Issues != "" && row.Issues != "[]" {
		var issues []string
		if err := json.Unmarshal([]byte(row.Issues), &issues); err == nil {
			resp.Issues = issues
		}
	}

	return resp
}

func (c *ControllerV1) ListDomainHealth(ctx context.Context, req *v1.ListDomainHealthReq) (res *v1.ListDomainHealthRes, err error) {
	res = &v1.ListDomainHealthRes{}

	results, err := domainhealth.DomainHealth().GetAllDomainHealth(ctx)
	if err != nil {
		res.SetError(err)
		return
	}

	// Convert DB rows to API response (parse Issues JSON string -> []string)
	apiResults := make([]domainHealthAPIResponse, len(results))
	for i := range results {
		apiResults[i] = convertRow(&results[i])
	}

	res.Data = apiResults
	res.SetSuccess("")
	return
}

func (c *ControllerV1) CheckDomain(ctx context.Context, req *v1.CheckDomainReq) (res *v1.CheckDomainRes, err error) {
	res = &v1.CheckDomainRes{}

	result, err := domainhealth.DomainHealth().CheckDomainHealth(ctx, req.Domain)
	if err != nil {
		res.SetError(err)
		return
	}

	res.Data = result
	res.SetSuccess("")
	return
}

func (c *ControllerV1) CheckAllDomains(ctx context.Context, req *v1.CheckAllDomainsReq) (res *v1.CheckAllDomainsRes, err error) {
	res = &v1.CheckAllDomainsRes{}

	results, err := domainhealth.DomainHealth().CheckAllDomains(ctx)
	if err != nil {
		res.SetError(err)
		return
	}

	res.Data = results
	res.SetSuccess("")
	return
}
