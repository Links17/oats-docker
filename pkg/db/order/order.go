package order

import (
	"context"
	dberrors "oats-docker/pkg/db/errors"
	"oats-docker/pkg/db/model"

	"gorm.io/gorm"
)

type OrderInterface interface {
	Update(ctx context.Context, requestId string, updates map[string]interface{}) error
	Get(ctx context.Context, requestId string) (*model.Order, error)
}

type order struct {
	db *gorm.DB
}

func NewOrder(db *gorm.DB) *order {
	return &order{db}
}

func (u *order) Update(ctx context.Context, requestId string, updates map[string]interface{}) error {
	f := u.db.Model(&model.Order{}).
		Where("request_id = ? ", requestId).
		Updates(updates)
	if f.Error != nil {
		return f.Error
	}

	if f.RowsAffected == 0 {
		return dberrors.ErrRecordNotUpdate
	}

	return nil
}

func (u *order) Get(ctx context.Context, requestId string) (*model.Order, error) {
	var obj model.Order
	if err := u.db.Where("request_id = ?", requestId).First(&obj).Error; err != nil {
		return nil, err
	}
	return &obj, nil
}
