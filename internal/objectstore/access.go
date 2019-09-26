/*
 * Copyright (c) 2019 Kubenext, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package objectstore

import (
	"fmt"
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

func newAccessCache() *accessCache {
	return &accessCache{
		access: accessMap{},
	}
}
