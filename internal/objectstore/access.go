/*
 * Copyright (c) 2019 Kubenext, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package objectstore

import (
	"sync"
)

// AccessKey is used at a key in an access map, It is made up of a namespace, Group, Resource and Verb.
type AccessKey struct {
	Namespace string
	Group     string
	Resource  string
	Verb      string
}

type accessMap map[AccessKey]bool

type accessCache struct {
	access accessMap
	mu     sync.RWMutex
}

func newAccessCache() *accessCache {
	return &accessCache{
		access: accessMap{},
	}
}

func (c *accessCache) reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.access = accessMap{}
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
