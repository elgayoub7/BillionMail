package contact

import (
	"context"

	v1 "billionmail-core/api/contact/v1"
	"billionmail-core/internal/service/contact"
	"billionmail-core/internal/service/public"

	"github.com/gogf/gf/v2/errors/gerror"
)

func (c *ControllerV1) PreviewImport(ctx context.Context, req *v1.PreviewImportReq) (res *v1.PreviewImportRes, err error) {
	res = &v1.PreviewImportRes{}

	content := req.FileData
	if content == "" {
		res.Code = 400
		res.SetError(gerror.New(public.LangCtx(ctx, "No data provided")))
		return res, nil
	}

	result := contact.ParseCSV(content, 5)

	res.Data.Headers = result.Headers
	res.Data.Mapping = result.Mapping
	res.Data.Unmapped = result.Unmapped
	res.Data.PreviewRows = result.Rows
	res.Data.TotalRows = result.TotalRows
	res.Data.Duplicates = result.Duplicates
	res.Data.InvalidRows = result.InvalidRows

	res.SetSuccess(public.LangCtx(ctx, "Preview generated successfully"))
	return res, nil
}
