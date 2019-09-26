/*
 * Copyright (c) 2019 Kubenext, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package store

import (
	"context"
	kubeonstructured "github.com/kubenext/kubeon/thirdparty/unstructured"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// GetAs gets an object from the object store by key. If the object is not found,
// return false and a nil error.
func GetAs(ctx context.Context, o Store, key Key, as interface{}) (bool, error) {
	u, found, err := o.Get(ctx, key)
	if err != nil {
		return false, errors.Wrap(err, "get object from object store")
	}

	if !found {
		return false, nil
	}

	// NOTE: (bryanl) vendored converter can't convert from int64 to float64. Watching
	// https://github.com/kubernetes-sigs/yaml/pull/14 to see when it gets pulled into
	// a release so Kubeon can switch back.
	if err := kubeonstructured.DefaultUnstructuredConverter.FromUnstructured(u.Object, as); err != nil {
		return false, errors.Wrap(err, "unable to convert object to unstructured")
	}

	if err := copyObjectMeta(as, u); err != nil {
		return false, errors.Wrap(err, "copy object metadata")
	}

	return true, nil
}

func copyObjectMeta(to interface{}, from *unstructured.Unstructured) error {
	object, ok := to.(metav1.Object)
	if !ok {
		return errors.Errorf("%T is not an object", to)
	}

	t, err := meta.TypeAccessor(object)
	if err != nil {
		return errors.Wrapf(err, "accessing type meta")
	}

	t.SetAPIVersion(from.GetAPIVersion())
	t.SetKind(from.GetObjectKind().GroupVersionKind().Kind)

	object.SetNamespace(from.GetNamespace())
	object.SetName(from.GetName())
	object.SetGenerateName(from.GetGenerateName())
	object.SetUID(from.GetUID())
	object.SetResourceVersion(from.GetResourceVersion())
	object.SetGeneration(from.GetGeneration())
	object.SetSelfLink(from.GetSelfLink())
	object.SetCreationTimestamp(from.GetCreationTimestamp())
	object.SetDeletionTimestamp(from.GetDeletionTimestamp())
	object.SetDeletionGracePeriodSeconds(from.GetDeletionGracePeriodSeconds())
	object.SetLabels(from.GetLabels())
	object.SetAnnotations(from.GetAnnotations())
	object.SetInitializers(from.GetInitializers())
	object.SetOwnerReferences(from.GetOwnerReferences())
	object.SetClusterName(from.GetClusterName())
	object.SetFinalizers(from.GetFinalizers())

	return nil
}
