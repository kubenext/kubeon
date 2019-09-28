package configuration

import (
	"context"
	"fmt"

	"github.com/vmware/octant/internal/log"
	"github.com/vmware/octant/internal/octant"
	"github.com/vmware/octant/pkg/action"
	"github.com/vmware/octant/pkg/store"
)

type ObjectDeleter struct {
	logger log.Logger
	store  store.Store
}

func NewObjectDeleter(logger log.Logger, clusterClient store.Store) *ObjectDeleter {
	return &ObjectDeleter{
		logger: logger.With("action", octant.ActionDeleteObject),
		store:  clusterClient,
	}
}

func (d *ObjectDeleter) ActionName() string {
	return octant.ActionDeleteObject
}

func (d *ObjectDeleter) Handle(ctx context.Context, alerter action.Alerter, payload action.Payload) error {
	d.logger.With("payload", payload).Debugf("deleting object")

	key, err := store.KeyFromPayload(payload)
	if err != nil {
		return err
	}

	alertType := action.AlertTypeInfo
	message := fmt.Sprintf("Deleted %s %q", key.Kind, key.Name)
	if err := d.store.Delete(ctx, key); err != nil {
		alertType = action.AlertTypeWarning
		message = fmt.Sprintf("Unable to deleted %s %q: %s", key.Kind, key.Name, err)
	}
	alert := action.CreateAlert(alertType, message, action.DefaultAlertExpiration)
	alerter.SendAlert(alert)

	return nil
}
