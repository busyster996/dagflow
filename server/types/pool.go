package types

type SPoolReq struct {
	Size int `json:"size" form:"size" yaml:"size" binding:"required"`
}

type SPoolRes struct {
	Size    int   `json:"size" yaml:"size"`
	Total   int64 `json:"total" yaml:"total"`
	Running int64 `json:"running" yaml:"running"`
	Waiting int64 `json:"waiting" yaml:"waiting"`
}
