package v1

import (
	"github.com/gogf/gf/v2/frame/g"

	"billionmail-core/utility/types/api_v1"
)

type GetLeadsReq struct {
	g.Meta        `path:"/scoring/leads" method:"get" tags:"Scoring" summary:"Get leads by engagement level"`
	Authorization string `json:"authorization" in:"header"`
	Level         string `json:"level"      d:"all"      dc:"Filter by level: cold, warm, hot, converted, all"`
	Page          int    `json:"page"       d:"1"        dc:"Page number"`
	PageSize      int    `json:"page_size"  d:"20"       dc:"Items per page"`
}

type GetLeadsRes struct {
	api_v1.StandardRes
}

type GetScoringStatsReq struct {
	g.Meta        `path:"/scoring/stats" method:"get" tags:"Scoring" summary:"Get scoring statistics"`
	Authorization string `json:"authorization" in:"header"`
}

type GetScoringStatsRes struct {
	api_v1.StandardRes
}

type GetLeadReq struct {
	g.Meta        `path:"/scoring/lead" method:"get" tags:"Scoring" summary:"Get lead score by email"`
	Authorization string `json:"authorization" in:"header"`
	Email         string `json:"email" v:"required|email" dc:"Contact email"`
}

type GetLeadRes struct {
	api_v1.StandardRes
}

type RecalculateScoresReq struct {
	g.Meta        `path:"/scoring/recalculate" method:"post" tags:"Scoring" summary:"Recalculate all lead scores"`
	Authorization string `json:"authorization" in:"header"`
}

type RecalculateScoresRes struct {
	api_v1.StandardRes
}
