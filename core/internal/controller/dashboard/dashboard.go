package dashboard

import (
	"context"

	v1 "billionmail-core/api/dashboard/v1"
	"billionmail-core/internal/service/dashboard"
)

type ControllerV1 struct{}

func NewV1() *ControllerV1 {
	return &ControllerV1{}
}

func (c *ControllerV1) GetStats(ctx context.Context, req *v1.GetStatsReq) (res *v1.GetStatsRes, err error) {
	res = &v1.GetStatsRes{}

	stats, err := dashboard.ColdDashboard().GetDashboardStats(ctx)
	if err != nil {
		res.SetError(err)
		return
	}

	res.Data = stats
	res.SetSuccess("")
	return
}

func (c *ControllerV1) GetActiveSequences(ctx context.Context, req *v1.GetActiveSequencesReq) (res *v1.GetActiveSequencesRes, err error) {
	res = &v1.GetActiveSequencesRes{}

	sequences, err := dashboard.ColdDashboard().GetActiveSequencesWithMetrics(ctx)
	if err != nil {
		res.SetError(err)
		return
	}

	res.Data = sequences
	res.SetSuccess("")
	return
}

func (c *ControllerV1) GetAlerts(ctx context.Context, req *v1.GetAlertsReq) (res *v1.GetAlertsRes, err error) {
	res = &v1.GetAlertsRes{}

	limit := req.Limit
	if limit <= 0 {
		limit = 20
	}

	alerts, err := dashboard.ColdDashboard().GetRecentAlerts(ctx, limit)
	if err != nil {
		res.SetError(err)
		return
	}

	res.Data = alerts
	res.SetSuccess("")
	return
}
