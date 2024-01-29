package types

import "time"

type IdUriMeta struct {
	Id int64 `uri:"id" binding:"required"`
}

type CloudUriMeta struct {
	IdUriMeta `json:",inline"`

	CloudName string `uri:"cloud_name" binding:"required"`
}

type IdMeta struct {
	Id              int64 `json:"id"`
	ResourceVersion int64 `json:"resource_version"`
}

// PageOptions 分页选项
type PageOptions struct {
	Limit int `form:"limit"`
	Page  int `form:"page"`
}

// TimeOption 通用时间规格
type TimeOption struct {
	GmtCreate   interface{} `json:"gmt_create,omitempty"`
	GmtModified interface{} `json:"gmt_modified,omitempty"`
}

const (
	timeLayout = "2006-01-02 15:04:05.999999999"
)

func NewTypeTime(GmtCreate time.Time, GmtModified time.Time) TimeOption {
	return TimeOption{
		GmtCreate:   GmtCreate.Format(timeLayout),
		GmtModified: GmtModified.Format(timeLayout),
	}
}
