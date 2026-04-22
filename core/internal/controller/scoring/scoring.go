package scoring

import (
	"context"

	v1 "billionmail-core/api/scoring/v1"
	"billionmail-core/internal/service/scoring"
)

type ControllerV1 struct{}

func NewV1() *ControllerV1 {
	return &ControllerV1{}
}

func (c *ControllerV1) GetLeads(ctx context.Context, req *v1.GetLeadsReq) (res *v1.GetLeadsRes, err error) {
	res = &v1.GetLeadsRes{}

	total, leads, err := scoring.Scoring().GetLeadsByLevel(ctx, req.Level, req.Page, req.PageSize)
	if err != nil {
		res.SetError(err)
		return
	}

	res.SetSuccess("")
	_ = total
	_ = leads
	return
}

func (c *ControllerV1) GetScoringStats(ctx context.Context, req *v1.GetScoringStatsReq) (res *v1.GetScoringStatsRes, err error) {
	res = &v1.GetScoringStatsRes{}

	stats, err := scoring.Scoring().GetStats(ctx)
	if err != nil {
		res.SetError(err)
		return
	}

	res.SetSuccess("")
	_ = stats
	return
}

func (c *ControllerV1) GetLead(ctx context.Context, req *v1.GetLeadReq) (res *v1.GetLeadRes, err error) {
	res = &v1.GetLeadRes{}

	lead, err := scoring.Scoring().GetOrCreateLeadScore(ctx, req.Email)
	if err != nil {
		res.SetError(err)
		return
	}

	res.SetSuccess("")
	_ = lead
	return
}

func (c *ControllerV1) RecalculateScores(ctx context.Context, req *v1.RecalculateScoresReq) (res *v1.RecalculateScoresRes, err error) {
	res = &v1.RecalculateScoresRes{}

	count, err := scoring.Scoring().RecalculateAll(ctx)
	if err != nil {
		res.SetError(err)
		return
	}

	res.SetSuccess("")
	_ = count
	return
}
