/*
 * Copyright (c) 2019 Kubenext, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package cluster

import (
	"github.com/magiconair/properties/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic/fake"
	"testing"
)

func TestNamespaceClient_Names(t *testing.T) {
	scheme := runtime.NewScheme()
	dc := fake.NewSimpleDynamicClient(scheme,
		newUnstructured("v1", "Namespace", "", "default"),
		newUnstructured("v1", "Namespace", "", "app-1"),
	)

	nc := newNamespaceClient(dc, "default")

	got, err := nc.Names()
	require.NoError(t, err)

	expected := []string{"default", "app-1"}
	assert.Equal(t, expected, got)
}

func newUnstructured(apiVersion, kind, namespace, name string) *unstructured.Unstructured {
	return &unstructured.Unstructured{Object: map[string]interface{}{
		"apiVersion": apiVersion,
		"kind":       kind,
		"metadata": map[string]interface{}{
			"namespace": namespace,
			"name":      name,
		},
	}}
}
