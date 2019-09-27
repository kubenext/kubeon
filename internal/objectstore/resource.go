/*
 * Copyright (c) 2019 Kubenext, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package objectstore

import (
	"context"
	"github.com/kubenext/kubeon/internal/cluster"
	"github.com/kubenext/kubeon/pkg/store"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
	authorizationv1 "k8s.io/api/authorization/v1"
	"sync"
)

type ResourceAccess interface {
	HasAccess(ctx context.Context, key store.Key, verb string) error
	Reset()
	Get(key AccessKey) (bool, bool)
	Set(key AccessKey, value bool)
	UpdateClient(client cluster.ClientInterface)
}

type resourceAccess struct {
	client cluster.ClientInterface
	cache  *accessCache
	mu     sync.RWMutex
}

var _ ResourceAccess = (*resourceAccess)(nil)

func NewResourceAccess(client cluster.ClientInterface) ResourceAccess {
	return &resourceAccess{
		client: client,
		cache:  newAccessCache(),
	}
}

func (r *resourceAccess) HasAccess(ctx context.Context, key store.Key, verb string) error {
	_, span := trace.StartSpan(ctx, "resourceAccessHasAccess")
	defer span.End()

	ak, err := r.keyToAccessKey(key, verb)
	if err != nil {
		return err
	}

	access, ok := r.cache.get(ak)
	if !ok {
		span.Annotate([]trace.Attribute{}, "fetch access start")
		val, err := r.fetchAccess(ak, verb)
		if err != nil {
			return errors.Wrapf(err, "fetch access: %+v", ak)
		}

		r.cache.set(ak, val)
		access = val
		span.Annotate([]trace.Attribute{}, "fetch access finish")
	}

	if !access {
		return &AccessError{Key: ak}
	}

	return nil
}

func (r *resourceAccess) Reset() {
	r.cache.reset()
}

func (r *resourceAccess) Get(key AccessKey) (bool, bool) {
	return r.cache.get(key)
}

func (r *resourceAccess) Set(key AccessKey, value bool) {
	r.cache.set(key, value)
}

func (r *resourceAccess) UpdateClient(client cluster.ClientInterface) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.client = client
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

	ak := AccessKey{
		Namespace: key.Namespace,
		Group:     gvr.Group,
		Resource:  gvr.Resource,
		Verb:      verb,
	}
	return ak, nil
}

func (r *resourceAccess) fetchAccess(key AccessKey, verb string) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	kubernetesClient, err := r.client.KubernetesClient()
	if err != nil {
		return false, errors.Wrap(err, "client kubernetes")
	}

	authClient := kubernetesClient.AuthorizationV1()
	sar := &authorizationv1.SelfSubjectAccessReview{
		Spec: authorizationv1.SelfSubjectAccessReviewSpec{
			ResourceAttributes: &authorizationv1.ResourceAttributes{
				Namespace: key.Namespace,
				Verb:      verb,
				Group:     key.Group,
				Resource:  key.Resource,
			},
		},
	}

	review, err := authClient.SelfSubjectAccessReviews().Create(sar)
	if err != nil {
		return false, errors.Wrap(err, "client auth")
	}
	return review.Status.Allowed, nil
}
