package domainhealth

import (
	"context"

	v1 "billionmail-core/api/domainhealth/v1"
	"billionmail-core/internal/service/domainhealth"
)

type ControllerV1 struct{}

func NewV1() *ControllerV1 {
	return &ControllerV1{}
}

func (c *ControllerV1) ListDomainHealth(ctx context.Context, req *v1.ListDomainHealthReq) (res *v1.ListDomainHealthRes, err error) {
	res = &v1.ListDomainHealthRes{}

	results, err := domainhealth.DomainHealth().GetAllDomainHealth(ctx)
	if err != nil {
		res.SetError(err)
		return
	}

	res.Data = results
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
