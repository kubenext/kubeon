/*
 * Copyright (c) 2019 Kubenext, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package cluster

import (
	"github.com/magiconair/properties/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/client-go/tools/clientcmd"
	"testing"
)

func TestClusterInfo(t *testing.T) {
	tests := []struct {
		name          string
		kubeConfig    []byte
		expectContext string
		expectCluster string
		expectServer  string
		expectUser    string
	}{
		{
			name: "general",
			kubeConfig: []byte(
				`
apiVersion: v1
clusters:
- cluster:
    server: https://other-localhost:6443
  name: other-cluster
- cluster:
    server: https://localhost:6443
  name: docker-for-desktop
contexts:
- context:
    cluster: docker-for-desktop
    user: docker-user
  name: main-context
- context:
    cluster: other-cluster
  name: other-context
current-context: main-context
`,
			),
			expectContext: "main-context",
			expectCluster: "docker-for-desktop",
			expectServer:  "https://localhost:6443",
			expectUser:    "docker-user",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			config, err := clientcmd.NewClientConfigFromBytes(tc.kubeConfig)
			require.NoError(t, err)

			ci := newClusterInfo(config)
			assert.Equal(t, tc.expectContext, ci.Context(), "unexpected context")
			assert.Equal(t, tc.expectCluster, ci.Cluster(), "unexpected cluster")
			assert.Equal(t, tc.expectServer, ci.Server(), "unexpected server")
			assert.Equal(t, tc.expectUser, ci.User(), "unexpected user")
		})
	}

}
