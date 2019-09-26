/*
 * Copyright (c) 2019 Kubenext, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package cluster

import (
	"context"
	"github.com/stretchr/testify/require"
	"path/filepath"
	"testing"
)

func TestFromKubeConfig(t *testing.T) {
	kubeConfig := filepath.Join("testdata", "kubeconfig.yaml")
	config := RestConfigOptions{}

	_, err := FromKubeConfig(context.TODO(), kubeConfig, "", config)
	require.NoError(t, err)
}
