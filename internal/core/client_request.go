/*
 * Copyright (c) 2019 Kubenext, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package core

import "github.com/kubenext/kubeon/pkg/action"

// ClientRequestHandler is a client request.
type ClientRequestHandler struct {
	RequestType string
	Handler     func(state State, payload action.Payload) error
}
