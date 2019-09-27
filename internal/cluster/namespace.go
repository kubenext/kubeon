/*
 * Copyright (c) 2019 Kubenext, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package cluster

type NamespaceInterface interface {
	Names() ([]string, error)
	InitialNamespace() string
}
