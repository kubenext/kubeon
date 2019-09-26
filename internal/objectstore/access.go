/*
 * Copyright (c) 2019 Kubenext, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package objectstore

import (
	"context"
	"fmt"
	"github.com/kubenext/kubeon/internal/cluster"
	"github.com/kubenext/kubeon/pkg/store"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
	authorizationv1 "k8s.io/api/authorization/v1"
	"sync"
)

// AccessKey is used at a key in an access map. It is made up of a Namespace, Group, Resource, and Verb.
type AccessKey struct {
	Namespace string
	Group     string
	Resource  string
	Verb      string
}

type AccessError struct {
	Key AccessKey
}

type accessMap map[AccessKey]bool

type accessCache struct {
	access accessMap
	mu     sync.RWMutex
}

type ResourceAccess interface {
	HasAccess(context.Context, store.Key, string) error
	Reset()
	Get(key AccessKey) (bool, bool)
	Set(AccessKey, bool)
	UpdateClient(cluster.ClientInterface)
}

type resourceAccess struct {
	client cluster.ClientInterface
	cache  *accessCache

	mu sync.RWMutex
}

func newAccessCache() *accessCache {
	return &accessCache{
		access: accessMap{},
	}
}

func NewResourceAccess(client cluster.ClientInterface) ResourceAccess {
	return &resourceAccess{
		client: client,
		cache:  newAccessCache(),
	}
}

func (r *resourceAccess) HasAccess(ctx context.Context, key store.Key, verb string) error {
	_, span := trace.StartSpan(ctx, "resourceAccessHasAccess")
	defer span.End()

	aKey, err := r.keyToAccessKey(key, verb)
	if err != nil {
		return err
	}

	access, ok := r.cache.get(aKey)

	if !ok {
		span.Annotate([]trace.Attribute{}, "fetch access start")
		val, err := r.fetchAccess(aKey, verb)
		if err != nil {
			return errors.Wrapf(err, "fetch access: %+v", aKey)
		}

		r.cache.set(aKey, val)
		access = val
		span.Annotate([]trace.Attribute{}, "fetch access finish")
	}

	if !access {
		return &AccessError{Key: aKey}
	}

	return nil
}

func (r *resourceAccess) Reset() {
	r.cache.reset()
}

func (r *resourceAccess) Get(key AccessKey) (bool, bool) {
	return r.cache.get(key)
}

func (r *resourceAccess) Set(key AccessKey, v bool) {
	r.cache.set(key, v)
}

func (r *resourceAccess) UpdateClient(cluster.ClientInterface) {
	panic("implement me")
}

func (r *resourceAccess) keyToAccessKey(key store.Key, verb string) (AccessKey, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	gvk := key.GroupVersionKind()

	if gvk.GroupKind().Empty() {
		return AccessKey{}, errors.Errorf("unable to check access for key %s", key.String())
	}

	gvr, err := r.client.Resource(gvk.GroupKind())
	if err != nil {
		return AccessKey{}, errors.Wrap(err, "client resource")
	}

	aKey := AccessKey{
		Namespace: key.Namespace,
		Group:     gvr.Group,
		Resource:  gvr.Resource,
		Verb:      verb,
	}
	return aKey, nil
}

func (r *resourceAccess) fetchAccess(key AccessKey, verb string) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	k8sClient, err := r.client.KubernetesClient()
	if err != nil {
		return false, errors.Wrap(err, "client kubernetes")
	}

	authClient := k8sClient.AuthorizationV1()
	sar := &authorizationv1.SelfSubjectAccessReview{
		Spec: authorizationv1.SelfSubjectAccessReviewSpec{
			ResourceAttributes: &authorizationv1.ResourceAttributes{
				Namespace: key.Namespace,
				Group:     key.Group,
				Resource:  key.Resource,
				Verb:      verb,
			},
		},
	}

	review, err := authClient.SelfSubjectAccessReviews().Create(sar)
	if err != nil {
		return false, errors.Wrap(err, "client auth")
	}
	return review.Status.Allowed, nil
}

func (ae *AccessError) Error() string {
	return fmt.Sprintf("access denied: no %s access in %s to %s/%s", ae.Key.Verb, ae.Key.Namespace, ae.Key.Group, ae.Key.Resource)
}

func (c *accessCache) set(key AccessKey, value bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.access[key] = value
}

func (c *accessCache) get(key AccessKey) (v, ok bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	v, ok = c.access[key]
	return v, ok
}

func (c *accessCache) reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.access = accessMap{}
}

var _ ResourceAccess = (*resourceAccess)(nil)
