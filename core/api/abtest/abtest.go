package abtest

import (
	"context"

	v1 "billionmail-core/api/abtest/v1"
)

type IAbTestV1 interface {
	CreateAbTest(ctx context.Context, req *v1.CreateAbTestReq) (res *v1.CreateAbTestRes, err error)
	GetAbTest(ctx context.Context, req *v1.GetAbTestReq) (res *v1.GetAbTestRes, err error)
	ListAbTests(ctx context.Context, req *v1.ListAbTestsReq) (res *v1.ListAbTestsRes, err error)
	GetAbTestResults(ctx context.Context, req *v1.GetAbTestResultsReq) (res *v1.GetAbTestResultsRes, err error)
	PickWinner(ctx context.Context, req *v1.PickWinnerReq) (res *v1.PickWinnerRes, err error)
	DeleteAbTest(ctx context.Context, req *v1.DeleteAbTestReq) (res *v1.DeleteAbTestRes, err error)
}
