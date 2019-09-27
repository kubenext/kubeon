/*
 * Copyright (c) 2019 Kubenext, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package objectstore

import "fmt"

type AccessError struct {
	Key AccessKey
}

func (ae *AccessError) Error() string {
	return fmt.Sprintf("access denied: no %s access in %s to %s/%s", ae.Key.Verb, ae.Key.Namespace, ae.Key.Group, ae.Key.Resource)
}
