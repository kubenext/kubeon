/*
 * Copyright (c) 2019 Kubenext, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package store

import (
	"context"
	tpunstructured "github.com/kubenext/kubeon/thirdparty/unstructured"
	"github.com/pkg/errors"
)

// GetAs gets an object from the object store by key. If the object is not found, return false and a nil error.
func GetAs(ctx context.Context, store Store, key Key, as interface{}) (bool, error) {
	u, found, err := store.Get(ctx, key)
	if err != nil {
		return false, errors.Wrap(err, "get object from object store")
	}

	if !found {
		return false, nil
	}

	// NOTE: (bryanl) vendored converter can't convert from int64 to float64. Watching
	// https://github.com/kubernetes-sigs/yaml/pull/14 to see when it gets pulled into
	// a release so Octant can switch back.
	if err := tpunstructured.DefaultUnstructuredConverter.FromUnstructured(u.Object, as); err != nil {
		return false, errors.Wrap(err, "unable to convert object to unstructured")
	}

	if err := copyObjectMeta(as, u); err != nil {
		return false, errors.Wrapf(err, "copy object metadata")
	}

	return true, nil
}
