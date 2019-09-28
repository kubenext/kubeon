/*
 * Copyright (c) 2019 Kubenext, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package core

import (
	"context"
	"fmt"
	"github.com/kubenext/kubeon/internal/log"
	"github.com/kubenext/kubeon/pkg/action"
	"github.com/kubenext/kubeon/pkg/store"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// DeploymentConfigurationEditor edits a deployment's configuration.
type DeploymentConfigurationEditor struct {
	logger log.Logger
	store  store.Store
}

var _ action.Dispatcher = (*DeploymentConfigurationEditor)(nil)

// NewDeploymentConfigurationEditor edits a deployment.
func NewDeploymentConfigurationEditor(logger log.Logger, objectStore store.Store) *DeploymentConfigurationEditor {
	return &DeploymentConfigurationEditor{
		logger: logger,
		store:  objectStore,
	}
}

// ActionName returns the action name for this editor.
func (e *DeploymentConfigurationEditor) ActionName() string {
	return "deployment/configuration"
}

// Handle edits a deployment. Supported edits:
//   * replicas
func (e *DeploymentConfigurationEditor) Handle(ctx context.Context, alerter action.Alerter, payload action.Payload) error {
	e.logger.
		With("payload", payload, "actionName", e.ActionName()).
		Debugf("received action payload")

	replicaCountFloat, err := payload.Float64("replicas")
	if err != nil {
		return err
	}
	replicaCount := roundToInt(replicaCountFloat)

	key, err := store.KeyFromPayload(payload)
	if err != nil {
		return err
	}

	name, err := payload.String("name")
	if err != nil {
		return err
	}

	fn := func(object *unstructured.Unstructured) error {
		return unstructured.SetNestedField(object.Object, replicaCount, "spec", "replicas")
	}

	var alertType = action.AlertTypeInfo
	message := fmt.Sprintf("Updated Deployment %q", name)
	if err := e.store.Update(ctx, key, fn); err != nil {
		alertType = action.AlertTypeWarning
		message = fmt.Sprintf("Unable to update Deployment %q: %s", name, err)
	}
	alert := action.CreateAlert(alertType, message, action.DefaultAlertExpiration)
	alerter.SendAlert(alert)

	return nil
}

func roundToInt(val float64) int64 {
	if val < 0 {
		return int64(val - 0.5)
	}
	return int64(val + 0.5)
}
