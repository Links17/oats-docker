package db

import (
	"gorm.io/gorm"
	"oats-docker/pkg/db/order"
)

type ShareDaoFactory interface {
	Order() order.OrderInterface
}

type shareDaoFactory struct {
	db *gorm.DB
}

func (f *shareDaoFactory) Order() order.OrderInterface {
	return order.NewOrder(f.db)
}
func NewDaoFactory(db *gorm.DB) ShareDaoFactory {
	return &shareDaoFactory{
		db: db,
	}
}
