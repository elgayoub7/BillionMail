package v1

import (
	"github.com/gogf/gf/v2/frame/g"

	"billionmail-core/utility/types/api_v1"
)

// --- Create ---

type VariantInput struct {
	Subject  string `json:"subject"   dc:"Email subject"`
	BodyHtml string `json:"body_html" dc:"HTML body"`
	BodyText string `json:"body_text" dc:"Plain text body"`
}

type CreateAbTestReq struct {
	g.Meta           `path:"/abtest/create" method:"post" tags:"ABTest" summary:"Create AB test"`
	Authorization    string        `json:"authorization" in:"header"`
	SequenceId       int           `json:"sequence_id"    v:"required|min:1"  dc:"Sequence ID"`
	StepId           int           `json:"step_id"        v:"required|min:0"  dc:"Step ID"`
	Name             string        `json:"name"                                dc:"Test name"`
	SplitRatio       float64       `json:"split_ratio"     d:"0.5"             dc:"Split ratio (0.0-1.0)"`
	WinnerCriteria   string        `json:"winner_criteria" d:"open_rate"       dc:"open_rate, click_rate, reply_rate"`
	VariantA         VariantInput  `json:"variant_a"       v:"required"        dc:"Variant A content"`
	VariantB         VariantInput  `json:"variant_b"       v:"required"        dc:"Variant B content"`
}

type CreateAbTestRes struct {
	api_v1.StandardRes
	Data struct {
		Id int `json:"id"`
	} `json:"data"`
}

// --- Get ---

type GetAbTestReq struct {
	g.Meta        `path:"/abtest/get" method:"get" tags:"ABTest" summary:"Get AB test detail"`
	Authorization string `json:"authorization" in:"header"`
	Id            int    `json:"id" v:"required|min:1" dc:"AB Test ID"`
}

type GetAbTestRes struct {
	api_v1.StandardRes
}

// --- List ---

type ListAbTestsReq struct {
	g.Meta        `path:"/abtest/list" method:"get" tags:"ABTest" summary:"List AB tests for sequence"`
	Authorization string `json:"authorization" in:"header"`
	SequenceId    int    `json:"sequence_id" v:"required|min:1" dc:"Sequence ID"`
}

type ListAbTestsRes struct {
	api_v1.StandardRes
}

// --- Results ---

type GetAbTestResultsReq struct {
	g.Meta        `path:"/abtest/results" method:"get" tags:"ABTest" summary:"Get AB test results"`
	Authorization string `json:"authorization" in:"header"`
	Id            int    `json:"id" v:"required|min:1" dc:"AB Test ID"`
}

type GetAbTestResultsRes struct {
	api_v1.StandardRes
}

// --- Pick Winner ---

type PickWinnerReq struct {
	g.Meta         `path:"/abtest/pick_winner" method:"post" tags:"ABTest" summary:"Manually pick AB test winner"`
	Authorization  string `json:"authorization" in:"header"`
	Id             int    `json:"id"              v:"required|min:1"   dc:"AB Test ID"`
	WinnerVariant  int    `json:"winner_variant"  v:"required|in:0,1"  dc:"Winner variant: 0=A, 1=B"`
}

type PickWinnerRes struct {
	api_v1.StandardRes
}

// --- Delete ---

type DeleteAbTestReq struct {
	g.Meta        `path:"/abtest/delete" method:"post" tags:"ABTest" summary:"Delete AB test"`
	Authorization string `json:"authorization" in:"header"`
	Id            int    `json:"id" v:"required|min:1" dc:"AB Test ID"`
}

type DeleteAbTestRes struct {
	api_v1.StandardRes
}
