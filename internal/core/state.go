/*
 * Copyright (c) 2019 Kubenext, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package core

import (
	"context"
	"github.com/kubenext/kubeon/pkg/action"
)

// UpdateCancelFunc cancels the update.
type UpdateCancelFunc func()

// State represents view state.
type State interface {
	// SetContentPath sets the content path.
	SetContentPath(string)
	// GetContentPath returns the content path.
	GetContentPath() string
	// OnNamespaceUpdate registers a function to be called with the content path
	// is changed.
	OnContentPathUpdate(fn ContentPathUpdateFunc) UpdateCancelFunc
	// GetQueryParams returns the query params.
	GetQueryParams() map[string][]string
	// SetNamespace sets the namespace.
	SetNamespace(namespace string)
	// GetNamespace returns the namespace.
	GetNamespace() string
	// OnNamespaceUpdate returns a function to be called when the namespace
	// is changed.
	OnNamespaceUpdate(fun NamespaceUpdateFunc) UpdateCancelFunc
	// AddFilter adds a label to filtered.
	AddFilter(filter Filter)
	// RemoveFilter removes a filter.
	RemoveFilter(filter Filter)
	// GetFilters returns a slice of filters.
	GetFilters() []Filter
	// SetFilters replaces the current filters with a slice of filters.
	// The slice can be empty.
	SetFilters(filters []Filter)
	// SetContext sets the current context.
	SetContext(requestedContext string)
	// Dispatch dispatches a payload for an action.
	Dispatch(ctx context.Context, actionName string, payload action.Payload) error
	// SendAlert sends an alert.
	SendAlert(alert action.Alert)
}

// ContentPathUpdateFunc is a function that is called when content path is updated.
type ContentPathUpdateFunc func(contentPath string)

// NamespaceUpdateFunc is a function that is called when namespace is updated.
type NamespaceUpdateFunc func(namespace string)
