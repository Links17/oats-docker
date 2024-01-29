package core

import (
	"context"
	"oats-docker/api/types"
	"oats-docker/cmd/app/config"
	"oats-docker/pkg/db"
	"oats-docker/pkg/db/model"
	"oats-docker/pkg/log"
)

const defaultJWTKey string = "oats"

type OrderGetter interface {
	Order() OrderInterface
}

type OrderInterface interface {
	Update(ctx context.Context, obj *types.Order) error
}

type order struct {
	ComponentConfig config.Config
	app             *oats
	factory         db.ShareDaoFactory
}

func newOrder(c *oats) OrderInterface {
	return &order{
		ComponentConfig: c.cfg,
		app:             c,
		factory:         c.factory,
	}
}

func (u *order) Update(ctx context.Context, obj *types.Order) error {
	oldOrder, err := u.factory.Order().Get(ctx, obj.RequestId)
	if err != nil {
		log.Logger.Errorf("failed to get order %d: %v", obj.RequestId)
		return err
	}
	updates := u.parseOrderUpdates(oldOrder, obj)
	if len(updates) == 0 {
		return nil
	}
	if err = u.factory.Order().Update(ctx, obj.RequestId, updates); err != nil {
		log.Logger.Errorf("failed to order user %d: %v", obj.RequestId, err)
		return err
	}

	return nil
}

func (u *order) parseOrderUpdates(oldObj *model.Order, newObj *types.Order) map[string]interface{} {
	updates := make(map[string]interface{})
	// 更新response
	if oldObj.ResponseStatus != newObj.Event.Value.ResponseStatus {
		updates["response_status"] = newObj.Event.Value.ResponseStatus
	}
	if oldObj.ResponseMsg != newObj.Event.Value.ResponseMsg {
		updates["response_msg"] = newObj.Event.Value.ResponseMsg
	}
	if oldObj.ExecStatus != newObj.Event.Value.ExecStatus {
		updates["exec_status"] = newObj.Event.Value.ExecStatus
	}
	if oldObj.ExecMsg != newObj.Event.Value.ExecMsg {
		updates["exec_msg"] = newObj.Event.Value.ExecMsg
	}
	return updates
}
