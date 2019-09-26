/*
 * Copyright (c) 2019 Kubenext, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package cluster

import "k8s.io/client-go/tools/clientcmd"

//go:generate mockgen -source=info.go -destination=./fake/mock_info_interface.go -package=fake github.com/kubenext/kubeon/internal/cluster InfoInterface

// InfoInterface provides connection details for a cluster
type InfoInterface interface {
	Context() string
	Cluster() string
	Server() string
	User() string
}

type clusterInfo struct {
	clientConfig clientcmd.ClientConfig
}

func (ci clusterInfo) Context() string {
	raw, err := ci.clientConfig.RawConfig()
	if err != nil {
		return ""
	}
	return raw.CurrentContext
}

func (ci clusterInfo) Cluster() string {
	raw, err := ci.clientConfig.RawConfig()
	if err != nil {
		return ""
	}
	ktx, ok := raw.Contexts[raw.CurrentContext]
	if !ok || ktx == nil {
		return ""
	}
	return ktx.Cluster
}

func (ci clusterInfo) Server() string {
	raw, err := ci.clientConfig.RawConfig()
	if err != nil {
		return ""
	}
	ktx, ok := raw.Contexts[raw.CurrentContext]
	if !ok || ktx == nil {
		return ""
	}
	c, ok := raw.Clusters[ktx.Cluster]
	if !ok || c == nil {
		return ""
	}
	return c.Server
}

func (ci clusterInfo) User() string {
	raw, err := ci.clientConfig.RawConfig()
	if err != nil {
		return ""
	}
	ktx, ok := raw.Contexts[raw.CurrentContext]
	if !ok || ktx == nil {
		return ""
	}
	return ktx.AuthInfo
}

func newClusterInfo(clientConfig clientcmd.ClientConfig) clusterInfo {
	return clusterInfo{clientConfig: clientConfig}
}
