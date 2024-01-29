package core

import (
	"github.com/bndr/gojenkins"

	"oats-docker/cmd/app/config"
	"oats-docker/pkg/db"
)

type CoreV1Interface interface {
	OrderGetter
}

type oats struct {
	cfg        config.Config
	factory    db.ShareDaoFactory
	cicdDriver *gojenkins.Jenkins
}

func New(cfg config.Config, factory db.ShareDaoFactory, cicdDriver *gojenkins.Jenkins) CoreV1Interface {
	return &oats{
		cfg:        cfg,
		factory:    factory,
		cicdDriver: cicdDriver,
	}
}
func (oats *oats) Order() OrderInterface {
	return newOrder(oats)
}
