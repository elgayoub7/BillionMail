package scoring

import (
	"context"

	v1 "billionmail-core/api/scoring/v1"
)

type IScoringV1 interface {
	GetLeads(ctx context.Context, req *v1.GetLeadsReq) (res *v1.GetLeadsRes, err error)
	GetScoringStats(ctx context.Context, req *v1.GetScoringStatsReq) (res *v1.GetScoringStatsRes, err error)
	GetLead(ctx context.Context, req *v1.GetLeadReq) (res *v1.GetLeadRes, err error)
	RecalculateScores(ctx context.Context, req *v1.RecalculateScoresReq) (res *v1.RecalculateScoresRes, err error)
}
