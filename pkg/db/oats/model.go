package oats

import "time"

type Model struct {
	Id              int64     `gorm:"column:id;primary_key;AUTO_INCREMENT;not null" json:"id"`
	GmtCreate       time.Time `json:"gmt_create"`
	GmtModified     time.Time `json:"gmt_modified"`
	ResourceVersion int64     `json:"resource_version"`
}
