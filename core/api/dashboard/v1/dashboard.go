package v1

import (
	"github.com/gogf/gf/v2/frame/g"

	"billionmail-core/utility/types/api_v1"
)

type GetStatsReq struct {
	g.Meta        `path:"/cold-dashboard/stats" method:"get" tags:"ColdDashboard" summary:"Get dashboard statistics"`
	Authorization string `json:"authorization" in:"header"`
}

type GetStatsRes struct {
	api_v1.StandardRes
}

type GetActiveSequencesReq struct {
	g.Meta        `path:"/cold-dashboard/active-sequences" method:"get" tags:"ColdDashboard" summary:"Get active sequences with metrics"`
	Authorization string `json:"authorization" in:"header"`
}

type GetActiveSequencesRes struct {
	api_v1.StandardRes
}

type GetAlertsReq struct {
	g.Meta        `path:"/cold-dashboard/alerts" method:"get" tags:"ColdDashboard" summary:"Get recent alerts"`
	Authorization string `json:"authorization" in:"header"`
	Limit         int `json:"limit" d:"20" dc:"Max alerts to return"`
}

type GetAlertsRes struct {
	api_v1.StandardRes
}
