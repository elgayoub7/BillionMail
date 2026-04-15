package batch_mail

import (
	"billionmail-core/api/batch_mail/v1"
	"billionmail-core/internal/service/batch_mail"
	"billionmail-core/internal/service/public"
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerV1) GetTaskRecipients(ctx context.Context, req *v1.GetTaskRecipientsReq) (res *v1.GetTaskRecipientsRes, err error) {
	res = &v1.GetTaskRecipientsRes{}

	// Verify task exists
	taskInfo, err := batch_mail.GetTaskInfo(ctx, req.TaskId)
	if err != nil {
		res.Code = 500
		res.SetError(gerror.New(public.LangCtx(ctx, "Failed to get task information: {}", err.Error())))
		return
	}

	if taskInfo == nil || taskInfo.Id == 0 {
		res.Code = 404
		res.SetError(gerror.New(public.LangCtx(ctx, "Task not found: {}", req.TaskId)))
		return
	}

	statService := batch_mail.NewTaskStatService()
	total, items, err := statService.GetTaskRecipients(req.TaskId, req.Status, req.Search, req.Page, req.PageSize)
	if err != nil {
		g.Log().Errorf(ctx, "Failed to get recipient analytics: %v", err)
		res.Code = 500
		res.SetError(gerror.New(public.LangCtx(ctx, "Failed to query recipient analytics: {}", err.Error())))
		return
	}

	res.Data.Total = total
	res.Data.List = items
	res.SetSuccess(public.LangCtx(ctx, "Recipient analytics retrieved successfully"))
	return
}
