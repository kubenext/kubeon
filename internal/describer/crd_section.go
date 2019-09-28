/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package describer

import (
	"context"
	"sort"
	"sync"

	"github.com/vmware/octant/internal/log"
	"github.com/vmware/octant/pkg/view/component"
)

type CRDSection struct {
	describers map[string]Describer
	path       string
	title      string

	mu sync.Mutex
}

var _ Describer = (*CRDSection)(nil)

func NewCRDSection(p, title string) *CRDSection {
	return &CRDSection{
		describers: make(map[string]Describer),
		path:       p,
		title:      title,
	}
}

func (csd *CRDSection) Add(name string, describer Describer) {
	csd.mu.Lock()
	defer csd.mu.Unlock()

	csd.describers[name] = describer
}

func (csd *CRDSection) Remove(name string) {
	csd.mu.Lock()
	defer csd.mu.Unlock()

	delete(csd.describers, name)
}

func (csd *CRDSection) Describe(ctx context.Context, namespace string, options Options) (component.ContentResponse, error) {
	csd.mu.Lock()
	defer csd.mu.Unlock()

	var names []string
	for name := range csd.describers {
		names = append(names, name)
	}

	sort.Strings(names)

	list := component.NewList("Custom Resources", nil)

	for _, name := range names {
		resp, err := csd.describers[name].Describe(ctx, namespace, options)
		if err != nil {
			return component.EmptyContentResponse, err
		}

		for i := range resp.Components {
			if nestedList, ok := resp.Components[i].(*component.List); ok {
				for i := range nestedList.Config.Items {
					item := nestedList.Config.Items[i]
					if !item.IsEmpty() {
						list.Add(item)
					}
				}
			}
		}
	}

	cr := component.ContentResponse{
		Components: []component.Component{list},
		Title:      component.TitleFromString(csd.title),
	}

	return cr, nil
}

func (csd *CRDSection) PathFilters() []PathFilter {
	return []PathFilter{
		*NewPathFilter(csd.path, csd),
	}
}

func (csd *CRDSection) Reset(ctx context.Context) error {
	csd.mu.Lock()
	defer csd.mu.Unlock()

	logger := log.From(ctx)

	for name := range csd.describers {
		logger.With("describer-name", name, "crd-section-path", csd.path).
			Debugf("removing crd from section")
		delete(csd.describers, name)
	}

	return nil
}
