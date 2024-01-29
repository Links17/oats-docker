package model

import (
	"oats-docker/pkg/db/oats"
	"time"
)

type Order struct {
	oats.Model
	RequestID      string    `gorm:"type:STRING(64);not null;primaryKey;column:request_id" json:"requestId"`
	DeviceUUID     string    `gorm:"type:STRING(64);not null;column:device_uuid" json:"deviceUUID"`
	FleetID        int       `gorm:"type:INTEGER;not null;column:fleet_id" json:"fleetId"`
	OrderName      int       `gorm:"type:INTEGER;not null;column:order_name" json:"orderName"`
	OrderData      JSONB     `gorm:"type:JSONB;not null;column:order_data" json:"orderData"`
	PublishTime    int64     `gorm:"type:BIGINT;not null;column:publish_time" json:"publishTime"`
	ResponseTime   int64     `gorm:"type:BIGINT;column:response_time" json:"responseTime"`
	ExecTime       int64     `gorm:"type:BIGINT;column:exec_time" json:"execTime"`
	ResponseStatus int       `gorm:"type:INTEGER;column:response_status" json:"responseStatus"`
	ResponseMsg    string    `gorm:"type:STRING;column:response_msg" json:"responseMsg"`
	ExecStatus     int       `gorm:"type:INTEGER;column:exec_status" json:"execStatus"`
	ExecMsg        string    `gorm:"type:STRING;column:exec_msg" json:"execMsg"`
	CreatedAt      time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt      time.Time `gorm:"column:updated_at" json:"updatedAt"`
}

func (order *Order) TableName() string {
	return "order_record"
}

type JSONB map[string]interface{}
