package batch_mail

import (
	"context"

	"billionmail-core/api/batch_mail/v1"
	"billionmail-core/internal/service/batch_mail"
)

func (c *ControllerV1) ScoreSubject(ctx context.Context, req *v1.ScoreSubjectReq) (res *v1.ScoreSubjectRes, err error) {
	result := batch_mail.ScoreSubject(req.Subject)
	res = &v1.ScoreSubjectRes{
		Score:    result.Score,
		Warnings: result.Warnings,
	}
	return
}
