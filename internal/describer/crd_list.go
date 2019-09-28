/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package describer

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	apiextv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/vmware/octant/internal/link"
	"github.com/vmware/octant/internal/printer"
	"github.com/vmware/octant/pkg/icon"
	"github.com/vmware/octant/pkg/store"
	"github.com/vmware/octant/pkg/view/component"
)

type crdListPrinter func(
	crdName string,
	crd *apiextv1beta1.CustomResourceDefinition,
	objects *unstructured.UnstructuredList,
	linkGenerator link.Interface,
	isLoading bool) (component.Component, error)

type crdListDescriptionOption func(*crdList)

type crdList struct {
	base

	name    string
	path    string
	printer crdListPrinter
}

var _ Describer = (*crdList)(nil)

func newCRDList(name, path string, options ...crdListDescriptionOption) *crdList {
	d := &crdList{
		name:    name,
		path:    path,
		printer: printer.CustomResourceListHandler,
	}

	for _, option := range options {
		option(d)
	}

	return d
}

func (cld *crdList) Describe(ctx context.Context, namespace string, options Options) (component.ContentResponse, error) {
	objectStore := options.ObjectStore()
	crd, err := CustomResourceDefinition(ctx, cld.name, objectStore)
	if err != nil {
		return component.EmptyContentResponse, err
	}

	objects, isLoading, err := ListCustomResources(ctx, crd, namespace, objectStore, options.LabelSet)
	if err != nil {
		return component.EmptyContentResponse, err
	}

	table, err := cld.printer(cld.name, crd, objects, options.Link, isLoading)
	if err != nil {
		return component.EmptyContentResponse, err
	}

	list := component.NewList(fmt.Sprintf("Custom Resources / %s", cld.name), []component.Component{
		table,
	})

	iconName, iconSource := loadIcon(icon.CustomResourceDefinition)
	list.SetIcon(iconName, iconSource)

	return component.ContentResponse{
		Components: []component.Component{list},
	}, nil
}

func ListCustomResources(
	ctx context.Context,
	crd *apiextv1beta1.CustomResourceDefinition,
	namespace string,
	o store.Store,
	selector *labels.Set) (*unstructured.UnstructuredList, bool, error) {
	if crd == nil {
		return nil, false, errors.New("crd is nil")
	}
	gvk := schema.GroupVersionKind{
		Group:   crd.Spec.Group,
		Version: crd.Spec.Version,
		Kind:    crd.Spec.Names.Kind,
	}

	apiVersion, kind := gvk.ToAPIVersionAndKind()

	key := store.Key{
		Namespace:  namespace,
		APIVersion: apiVersion,
		Kind:       kind,
		Selector:   selector,
	}

	objects, isLoading, err := o.List(ctx, key)
	if err != nil {
		return nil, false, errors.Wrapf(err, "listing custom resources for %q", crd.Name)
	}

	return objects, isLoading, nil
}

func (cld *crdList) PathFilters() []PathFilter {
	return []PathFilter{
		*NewPathFilter(cld.path, cld),
	}
}
