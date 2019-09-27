/*
 * Copyright (c) 2019 Kubenext, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package store

import (
	"fmt"
	"github.com/kubenext/kubeon/pkg/action"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"strings"
)

// Key is a key for the object store.
type Key struct {
	Namespace  string
	ApiVersion string
	Kind       string
	Name       string
	Selector   *labels.Set
}

// Convert Key to a string
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

// Converts the Key to a GroupVersionKind.
func (k Key) GroupVersionKind() schema.GroupVersionKind {
	return schema.FromAPIVersionAndKind(k.ApiVersion, k.Kind)
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

	return Key{
		Namespace:  namespace,
		ApiVersion: apiVersion,
		Kind:       kind,
		Name:       name,
	}, nil
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
