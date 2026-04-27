package v1

import (
	"billionmail-core/utility/types/api_v1"

	"github.com/gogf/gf/v2/frame/g"
)

// PreviewImportReq Preview CSV import with auto-mapping
type PreviewImportReq struct {
	g.Meta     `path:"/contact/group/preview_import" method:"post" tags:"Contact" summary:"Preview CSV import with auto-mapping"`
	Authorization string `json:"authorization" dc:"Authorization" in:"header"`
	FileData   string `json:"file_data" dc:"CSV file content"`
	ImportType int    `json:"import_type" v:"required|in:1,2" dc:"Import type (1: file, 2: paste)"`
}

// PreviewImportRes Preview CSV import response
type PreviewImportRes struct {
	api_v1.StandardRes
	Data struct {
		Headers     []string            `json:"headers"`
		Mapping     map[string]string   `json:"mapping"`
		Unmapped    []string            `json:"unmapped"`
		PreviewRows []map[string]string `json:"preview_rows"`
		TotalRows   int                 `json:"total_rows"`
		Duplicates  int                 `json:"duplicates"`
		InvalidRows int                 `json:"invalid_rows"`
	} `json:"data"`
}
