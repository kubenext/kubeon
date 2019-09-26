/*
 * Copyright (c) 2019 Kubenext, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package store

import (
	"context"
	"fmt"
	"github.com/kubenext/kubeon/internal/cluster"
	"github.com/kubenext/kubeon/pkg/action"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/tools/cache"
	"strings"
)

type Key struct {
	Namespace  string
	ApiVersion string
	Kind       string
	Name       string
	Selector   *labels.Set
}

type Store interface {
	List(ctx context.Context, key Key) (list *unstructured.UnstructuredList, loading bool, err error)
	Get(ctx context.Context, key Key) (object *unstructured.Unstructured, found bool, err error)
	Delete(ctx context.Context, key Key) error
	Watch(ctx context.Context, key Key, handler cache.ResourceEventHandler) error
	Unwatch(ctx context.Context, groupVersionKinds ...schema.GroupVersionKind) error
	UpdateClusterClient(ctx context.Context, client cluster.ClientInterface) error
	RegisterOnUpdate(fn UpdateFn)
	Update(ctx context.Context, key Key, updater func(*unstructured.Unstructured) error) error
	IsLoading(ctx context.Context, key Key) bool
}

type UpdateFn func(store Store)

func (k Key) String() string {
	var sb strings.Builder
	sb.WriteString("CacheKey[")
	if k.Namespace != "" {
		sb.WriteString(fmt.Sprintf("Namespace='%s', ", k.Namespace))
	}
	sb.WriteString(fmt.Sprintf("ApiVersion='%s', ", k.ApiVersion))
	sb.WriteString(fmt.Sprintf("Kind='%s'", k.Kind))

	if k.Name != "" {
		sb.WriteString(fmt.Sprintf(", Name='%s'", k.Name))
	}

	if k.Selector != nil && k.Selector.String() != "" {
		sb.WriteString(fmt.Sprintf(", Selector='%s'", k.Selector.String()))
	}

	sb.WriteString("]")

	return sb.String()
}

// GroupVersionKind converts the Key to a GroupVersionKind.
func (k Key) GroupVersionKind() schema.GroupVersionKind {
	return schema.FromAPIVersionAndKind(k.ApiVersion, k.Kind)
}

// KeyFromObject creates a key from a runtime object.
func KeyFromObject(object runtime.Object) (Key, error) {
	accessor := meta.NewAccessor()
	namespace, err := accessor.Namespace(object)
	if err != nil {
		return Key{}, err
	}

	apiVersion, err := accessor.APIVersion(object)
	if err != nil {
		return Key{}, err
	}

	kind, err := accessor.Kind(object)
	if err != nil {
		return Key{}, err
	}

	name, err := accessor.Name(object)
	if err != nil {
		return Key{}, err
	}

	return Key{
		Namespace:  namespace,
		ApiVersion: apiVersion,
		Kind:       kind,
		Name:       name,
	}, nil
}

// KeyFromGroupVersionKind creates a key from a group version kind.
func KeyFromGroupVersionKind(gvk schema.GroupVersionKind) Key {
	apiVersion, kind := gvk.ToAPIVersionAndKind()
	return Key{
		ApiVersion: apiVersion,
		Kind:       kind,
	}
}

// KeyFromPayload converts a payload into a Key.
func KeyFromPayload(payload action.Payload) (Key, error) {
	namespace, err := payload.OptionalString("namespace")
	if err != nil {
		return Key{}, err
	}
	apiVersion, err := payload.String("apiVersion")
	if err != nil {
		return Key{}, err
	}
	kind, err := payload.String("kind")
	if err != nil {
		return Key{}, err
	}
	name, err := payload.String("name")
	if err != nil {
		return Key{}, err
	}

	key := Key{
		Namespace:  namespace,
		ApiVersion: apiVersion,
		Kind:       kind,
		Name:       name,
	}

	return key, nil
}

// ToActionPayload converts the Key to a payload.
func (k Key) ToActionPayload() action.Payload {
	return action.Payload{
		"namespace":  k.Namespace,
		"apiVersion": k.ApiVersion,
		"kind":       k.Kind,
		"name":       k.Name,
	}
}
