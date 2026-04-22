package v1

import (
	"github.com/gogf/gf/v2/frame/g"

	"billionmail-core/utility/types/api_v1"
)

type ListDomainHealthReq struct {
	g.Meta        `path:"/domainhealth/list" method:"get" tags:"DomainHealth" summary:"List domain health status"`
	Authorization string `json:"authorization" in:"header"`
}

type ListDomainHealthRes struct {
	api_v1.StandardRes
}

type CheckDomainReq struct {
	g.Meta        `path:"/domainhealth/check" method:"post" tags:"DomainHealth" summary:"Check domain health"`
	Authorization string `json:"authorization" in:"header"`
	Domain        string `json:"domain" v:"required" dc:"Domain to check"`
}

type CheckDomainRes struct {
	api_v1.StandardRes
}

type CheckAllDomainsReq struct {
	g.Meta        `path:"/domainhealth/check_all" method:"post" tags:"DomainHealth" summary:"Check all domains"`
	Authorization string `json:"authorization" in:"header"`
}

type CheckAllDomainsRes struct {
	api_v1.StandardRes
}
