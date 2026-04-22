package domainhealth

import (
	"context"

	v1 "billionmail-core/api/domainhealth/v1"
)

type IDomainHealthV1 interface {
	ListDomainHealth(ctx context.Context, req *v1.ListDomainHealthReq) (res *v1.ListDomainHealthRes, err error)
	CheckDomain(ctx context.Context, req *v1.CheckDomainReq) (res *v1.CheckDomainRes, err error)
	CheckAllDomains(ctx context.Context, req *v1.CheckAllDomainsReq) (res *v1.CheckAllDomainsRes, err error)
}
