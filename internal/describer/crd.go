/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package describer

import (
	"context"

	"github.com/pkg/errors"
	apiextv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/vmware/octant/internal/config"
	"github.com/vmware/octant/internal/link"
	"github.com/vmware/octant/internal/modules/overview/yamlviewer"
	"github.com/vmware/octant/internal/printer"
	"github.com/vmware/octant/internal/queryer"
	"github.com/vmware/octant/internal/resourceviewer"
	"github.com/vmware/octant/pkg/icon"
	"github.com/vmware/octant/pkg/store"
	"github.com/vmware/octant/pkg/view/component"
)

type crdPrinter func(ctx context.Context, crd *apiextv1beta1.CustomResourceDefinition, object *unstructured.Unstructured, options printer.Options) (component.Component, error)
type resourceViewerPrinter func(ctx context.Context, object *unstructured.Unstructured, dashConfig config.Dash, q queryer.Queryer) (component.Component, error)
type yamlPrinter func(runtime.Object) (*component.YAML, error)

type crdOption func(*crd)

type crd struct {
	base

	path                  string
	name                  string
	summaryPrinter        crdPrinter
	resourceViewerPrinter resourceViewerPrinter
	yamlPrinter           yamlPrinter
}

var _ Describer = (*crd)(nil)

func newCRD(name, path string, options ...crdOption) *crd {
	d := &crd{
		path:                  path,
		name:                  name,
		summaryPrinter:        printer.CustomResourceHandler,
		resourceViewerPrinter: createCRDResourceViewer,
		yamlPrinter:           yamlviewer.ToComponent,
	}

	for _, option := range options {
		option(d)
	}

	return d
}

func (c *crd) Describe(ctx context.Context, namespace string, options Options) (component.ContentResponse, error) {
	objectStore := options.ObjectStore()
	crd, err := CustomResourceDefinition(ctx, c.name, objectStore)
	if err != nil {
		return component.EmptyContentResponse, err
	}

	// TODO: crd.Spec.Version is incorrect. Use crd.Spec.Version instead.
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
		Name:       options.Fields["name"],
	}

	object, found, err := objectStore.Get(ctx, key)
	if err != nil {
		return component.EmptyContentResponse, err
	}

	if !found {
		return component.EmptyContentResponse, err
	}

	title := component.Title(
		component.NewText("Custom Resources"),
		component.NewText(crd.Name),
		component.NewText(object.GetName()))

	iconName, iconSource := loadIcon(icon.CustomResourceDefinition)
	cr := component.NewContentResponse(title)
	cr.IconName = iconName
	cr.IconSource = iconSource

	linkGenerator, err := link.NewFromDashConfig(options)
	if err != nil {
		return component.EmptyContentResponse, err
	}

	printOptions := printer.Options{
		DashConfig: options,
		Link:       linkGenerator,
	}

	summary, err := c.summaryPrinter(ctx, crd, object, printOptions)
	if err != nil {
		return component.EmptyContentResponse, err
	}
	summary.SetAccessor("summary")

	cr.Add(summary)

	resourceViewerComponent, err := c.resourceViewerPrinter(ctx, object, options, options.Queryer)
	if err != nil {
		return component.EmptyContentResponse, err
	}

	resourceViewerComponent.SetAccessor("resourceViewer")
	cr.Add(resourceViewerComponent)

	yvComponent, err := c.yamlPrinter(object)
	if err != nil {
		return component.EmptyContentResponse, err
	}

	yvComponent.SetAccessor("yaml")
	cr.Add(yvComponent)

	pluginPrinter := options.PluginManager()
	tabs, err := pluginPrinter.Tabs(ctx, object)
	if err != nil {
		return component.EmptyContentResponse, errors.Wrap(err, "getting tabs from plugins")
	}

	for _, tab := range tabs {
		tab.Contents.SetAccessor(tab.Name)
		cr.Add(&tab.Contents)
	}

	return *cr, nil
}

func (c *crd) PathFilters() []PathFilter {
	return []PathFilter{
		*NewPathFilter(c.path, c),
	}
}

func createCRDResourceViewer(ctx context.Context, object *unstructured.Unstructured, dashConfig config.Dash, q queryer.Queryer) (component.Component, error) {
	rv, err := resourceviewer.New(dashConfig, resourceviewer.WithDefaultQueryer(dashConfig, q))
	if err != nil {
		return nil, err
	}

	handler, err := resourceviewer.NewHandler(dashConfig)
	if err != nil {
		return nil, err
	}

	if err := rv.Visit(ctx, object, handler); err != nil {
		return nil, err
	}

	return resourceviewer.GenerateComponent(ctx, handler, "")
}
