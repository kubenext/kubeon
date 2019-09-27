/*
 * Copyright (c) 2019 Kubenext, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package objectstore

import (
	"github.com/kubenext/kubeon/pkg/store"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sync"
)

type informerSynced struct {
	status map[string]bool
	mu     sync.RWMutex
}

func initInformerSynced() *informerSynced {
	return &informerSynced{
		status: make(map[string]bool),
	}
}

func (c *informerSynced) setSynced(key store.Key, value bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.status[key.String()] = value
}

func (c *informerSynced) hasSynced(key store.Key) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	v, ok := c.status[key.String()]
	if !ok {
		return true
	}

	return v
}

func (c *informerSynced) hasSeen(key store.Key) bool {
	_, ok := c.status[key.String()]
	return ok
}

func (c *informerSynced) reset() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for key := range c.status {
		delete(c.status, key)
	}
}

type factoriesCache struct {
	factories map[string]InformerFactory
	mu        sync.RWMutex
}

func initFactoriesCache() *factoriesCache {
	return &factoriesCache{
		factories: make(map[string]InformerFactory),
	}
}

func (c *factoriesCache) set(key string, value InformerFactory) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.factories[key] = value
}

func (c *factoriesCache) keys() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var list []string
	for k := range c.factories {
		list = append(list, k)
	}
	return list
}

func (c *factoriesCache) get(key string) (InformerFactory, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	v, ok := c.factories[key]
	return v, ok
}

func (c *factoriesCache) delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.factories, key)
}

func (c *factoriesCache) reset() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for key := range c.factories {
		delete(c.factories, key)
	}
}

type seenGvksCache struct {
	seenGvks map[string]map[schema.GroupVersionKind]bool
	mu       sync.RWMutex
}

func initSeenGvksCache() *seenGvksCache {
	return &seenGvksCache{
		seenGvks: make(map[string]map[schema.GroupVersionKind]bool),
	}
}

func (c *seenGvksCache) setSeen(key string, gvk schema.GroupVersionKind, value bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	cur, ok := c.seenGvks[key]
	if !ok {
		cur = make(map[schema.GroupVersionKind]bool)
	}

	cur[gvk] = value
	c.seenGvks[key] = cur
}

func (c *seenGvksCache) hasSeen(key string, gvk schema.GroupVersionKind) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	v, ok := c.seenGvks[key]
	if !ok {
		return false
	}

	seen, ok := v[gvk]
	if !ok {
		return false
	}
	return seen
}

func (c *seenGvksCache) reset() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for k := range c.seenGvks {
		delete(c.seenGvks, k)
	}
}

type informerContextCache struct {
	cache map[schema.GroupVersionResource]chan struct{}
	mu    sync.Mutex
}

func initInformerContextCache() *informerContextCache {
	return &informerContextCache{
		cache: make(map[schema.GroupVersionResource]chan struct{}),
	}
}

func (c *informerContextCache) addChild(gvr schema.GroupVersionResource) <-chan struct{} {
	c.mu.Lock()
	defer c.mu.Unlock()

	ch := make(chan struct{}, 1)
	c.cache[gvr] = ch
	return ch
}

func (c *informerContextCache) delete(gvr schema.GroupVersionResource) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if stopCh, ok := c.cache[gvr]; ok {
		close(stopCh)
		delete(c.cache, gvr)
	}
}

func (c *informerContextCache) reset() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for k, stopCh := range c.cache {
		close(stopCh)
		delete(c.cache, k)
	}
}
