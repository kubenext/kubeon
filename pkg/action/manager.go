/*
 * Copyright (c) 2019 Kubenext, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package action

import (
	"context"
	"github.com/kubenext/kubeon/internal/log"
	"sync"
)

// DispatcherFunc is a function that will be dispatched to handle a payload.
type DispatcherFunc func(ctx context.Context, alerter Alerter, payload Payload) error

// Dispatcher handles actions.
type Dispatcher interface {
	ActionName() string
	Handle(ctx context.Context, alerter Alerter, payload Payload) error
}

// Dispatchers is a slice of Dispatcher.
type Dispatchers []Dispatcher

// ToActionPaths converts Dispatchers to a map.
func (d Dispatchers) ToActionPaths() map[string]DispatcherFunc {
	m := make(map[string]DispatcherFunc)

	for i := range d {
		m[d[i].ActionName()] = d[i].Handle
	}

	return m
}

// Manager manages actions.
type Manager struct {
	logger     log.Logger
	dispatches map[string]DispatcherFunc
	mu         sync.Mutex
}

// NewManager creates an instance of Manager.
func NewManager(logger log.Logger) *Manager {
	return &Manager{
		logger:     logger.With("component", "action-manager"),
		dispatches: make(map[string]DispatcherFunc),
	}
}

// Register registers a dispatcher function to an action path.
func (m *Manager) Register(actionPath string, actionFunc DispatcherFunc) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.dispatches[actionPath] = actionFunc
	return nil
}

// Dispatch dispatches a payload to a path.
func (m *Manager) Dispatch(ctx context.Context, alerter Alerter, actionPath string, payload Payload) error {
	fn, ok := m.dispatches[actionPath]
	if !ok {
		return &NotFoundError{Path: actionPath}
	}
	return fn(ctx, alerter, payload)
}
