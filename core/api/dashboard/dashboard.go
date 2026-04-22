package dashboard

import (
	"context"

	v1 "billionmail-core/api/dashboard/v1"
)

type IDashboardV1 interface {
	GetStats(ctx context.Context, req *v1.GetStatsReq) (res *v1.GetStatsRes, err error)
	GetActiveSequences(ctx context.Context, req *v1.GetActiveSequencesReq) (res *v1.GetActiveSequencesRes, err error)
	GetAlerts(ctx context.Context, req *v1.GetAlertsReq) (res *v1.GetAlertsRes, err error)
}
