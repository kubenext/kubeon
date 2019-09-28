/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package localcontent

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/vmware/octant/internal/module"
	"github.com/vmware/octant/internal/octant"
	"github.com/vmware/octant/pkg/navigation"
	"github.com/vmware/octant/pkg/view/component"
)

type LocalContent struct {
	root string
}

var _ module.Module = (*LocalContent)(nil)

func New(root string) *LocalContent {
	return &LocalContent{
		root: root,
	}
}

func (l *LocalContent) Root() string {
	return l.root
}

func (l *LocalContent) Name() string {
	return "local"
}

func (l *LocalContent) ClientRequestHandlers() []octant.ClientRequestHandler {
	return nil
}

func (l *LocalContent) Content(ctx context.Context, contentPath string, opts module.ContentOptions) (component.ContentResponse, error) {
	if contentPath == "/" || contentPath == "" {
		return l.list()
	}

	fileName := fmt.Sprintf("%s.json", contentPath)
	return l.content(fileName)
}

func (l *LocalContent) list() (component.ContentResponse, error) {
	var out component.ContentResponse

	cols := component.NewTableCols("Title", "File")
	table := component.NewTable("Local Components", "We couldn't find any local content!", cols)

	err := l.walk(func(name, base string, content component.ContentResponse) error {
		title, err := l.titleToText(content.Title)
		if err != nil {
			return errors.Wrap(err, "convert title to text")
		}

		table.Add(component.TableRow{
			"Title": component.NewLink("", title, path.Join("/local", base)),
			"File":  component.NewText(name),
		})

		return nil
	})

	if err != nil {
		return out, nil
	}

	out.Title = component.Title(component.NewText("Local Contents"))

	out.Components = []component.Component{
		table,
	}

	return out, nil
}

func (l *LocalContent) ContentPath() string {
	return fmt.Sprintf("%s", l.Name())
}

func (l *LocalContent) content(name string) (component.ContentResponse, error) {
	var out component.ContentResponse

	f, err := os.Open(filepath.Join(l.root, name))
	if err != nil {
		return out, errors.Wrap(err, "open local content")
	}
	defer f.Close()

	if err = json.NewDecoder(f).Decode(&out); err != nil {
		return out, errors.Wrap(err, "read JSON")
	}

	return out, nil
}

type walkFn func(name, base string, content component.ContentResponse) error

func (l *LocalContent) walk(fn walkFn) error {
	if fn == nil {
		return errors.New("walkFn is nil")
	}

	fis, err := ioutil.ReadDir(l.root)
	if err != nil {
		return errors.Wrapf(err, "read %s", l.root)
	}

	for _, fi := range fis {
		if fi.IsDir() {
			continue
		}

		ext := filepath.Ext(fi.Name())
		if ext != ".json" {
			continue
		}

		content, err := l.content(fi.Name())
		if err != nil {
			return err
		}

		base := strings.TrimSuffix(fi.Name(), ext)
		if err = fn(fi.Name(), base, content); err != nil {
			return err
		}
	}

	return nil
}

func (l *LocalContent) Navigation(ctx context.Context, namespace, root string) ([]navigation.Navigation, error) {
	if !strings.HasSuffix(root, "/") {
		root = fmt.Sprintf("%s/", root)
	}
	nav := navigation.Navigation{
		Title:    "Local Content",
		Path:     root,
		Children: []navigation.Navigation{},
	}

	err := l.walk(func(name, base string, content component.ContentResponse) error {
		title, err := l.titleToText(content.Title)
		if err != nil {
			return errors.Wrap(err, "convert title to text")
		}

		nav.Children = append(nav.Children, navigation.Navigation{
			Title: title,
			Path:  path.Join(root, base),
		})

		return nil
	})

	if err != nil {
		return nil, err
	}

	return []navigation.Navigation{nav}, nil
}

func (l *LocalContent) titleToText(title []component.TitleComponent) (string, error) {
	if len(title) < 1 {
		return "", errors.New("title is empty")
	}

	var parts []string
	for _, titlePart := range title {
		text, ok := titlePart.(*component.Text)
		if !ok {
			return "", errors.New("title has a not text component")
		}
		parts = append(parts, text.Config.Text)
	}

	return strings.Join(parts, " > "), nil
}

func (l *LocalContent) SetNamespace(namespace string) error {
	return nil
}

func (l *LocalContent) Start() error {
	return nil
}

func (l *LocalContent) Stop() {
}

func (l *LocalContent) SupportedGroupVersionKind() []schema.GroupVersionKind {
	return []schema.GroupVersionKind{}
}

func (l *LocalContent) GroupVersionKindPath(namespace, apiVersion, kind, name string) (string, error) {
	return "", errors.Errorf("local content can't create paths for %s %s", apiVersion, kind)
}

func (l *LocalContent) AddCRD(ctx context.Context, crd *unstructured.Unstructured) error {
	return nil
}

func (l *LocalContent) RemoveCRD(ctx context.Context, crd *unstructured.Unstructured) error {
	return nil
}

func (l *LocalContent) ResetCRDs(ctx context.Context) error {
	return nil
}

// Generators allow modules to send events to the frontend.
func (l *LocalContent) Generators() []octant.Generator {
	return []octant.Generator{}
}

func (l *LocalContent) SetContext(ctx context.Context, contextName string) error {
	return nil
}
