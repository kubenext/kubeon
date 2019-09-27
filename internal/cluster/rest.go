/*
 * Copyright (c) 2019 Kubenext, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package cluster

import "k8s.io/client-go/rest"

type RestInterface interface {
	RestClient() (rest.Interface, error)
	RestConfig() *rest.Config
}
