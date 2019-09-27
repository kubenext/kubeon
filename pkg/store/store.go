/*
 * Copyright (c) 2019 Kubenext, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package store

import (
	"context"
	"github.com/kubenext/kubeon/internal/cluster"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/tools/cache"
)

// UpdateFn is a function that is called when
type UpdateFn func(store Store)

// Store stores kubernetes objects.
type Store interface {
	List(ctx context.Context, key Key) (list *unstructured.UnstructuredList, loading bool, err error)
	Get(ctx context.Context, key Key) (object *unstructured.Unstructured, found bool, err error)
	Delete(ctx context.Context, key Key) error
	Watch(ctx context.Context, key Key, handler cache.ResourceEventHandler) error
	Unwatch(ctx context.Context, gvk ...schema.GroupVersionKind) error
	UpdateClusterClient(ctx context.Context, client cluster.ClientInterface) error
	RegisterOnUpdate(fn UpdateFn)
	Update(ctx context.Context, key Key, updater func(*unstructured.Unstructured) error) error
	IsLoading(ctx context.Context, key Key) bool
}
