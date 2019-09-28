/*
 * Copyright (c) 2019 Kubenext, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package cluster

import (
	"context"
	"fmt"
	"github.com/kubenext/kubeon/internal/log"
	"github.com/kubenext/kubeon/internal/utils/strings"
	"github.com/pkg/errors"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/discovery/cached/disk"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
	"time"
)

type ClientInterface interface {
	DefaultNamespace() string
	ResourceExists(schema.GroupVersionResource) bool
	Resource(kind schema.GroupKind) (schema.GroupVersionResource, error)
	KubernetesClient() (kubernetes.Interface, error)
	DynamicClient() (dynamic.Interface, error)
	DiscoveryClient() (discovery.DiscoveryInterface, error)
	NamespaceClient() (NamespaceInterface, error)
	InfoClient() (InfoInterface, error)
	Close()
	RestInterface
}

// Cluster is a client for cluster operations
type Cluster struct {
	clientConfig     clientcmd.ClientConfig
	restConfig       *rest.Config
	logger           log.Logger
	kubernetesClient kubernetes.Interface
	dynamicClient    dynamic.Interface
	discoveryClient  discovery.DiscoveryInterface
	restMapper       *restmapper.DeferredDiscoveryRESTMapper
	closeFn          context.CancelFunc
	defaultNamespace string
}

var _ ClientInterface = (*Cluster)(nil)

func newCluster(ctx context.Context, clientConfig clientcmd.ClientConfig, restConfig *rest.Config, defaultNamespace string) (*Cluster, error) {
	logger := log.From(ctx).With("component", "cluster client")

	kubernetesClient, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, errors.Wrap(err, "create kubernetes client")
	}

	dynamicClient, err := dynamic.NewForConfig(restConfig)
	if err != nil {
		return nil, errors.Wrap(err, "create dynamic client")
	}

	discoveryClient, err := discovery.NewDiscoveryClientForConfig(restConfig)
	if err != nil {
		return nil, errors.Wrap(err, "create discovery client")
	}

	dir, err := ioutil.TempDir("", "kubeon")
	if err != nil {
		return nil, errors.Wrap(err, "create temp directory")
	}

	logger.With("dir", dir).Debugf("created temp directory")

	cacheDiscoverClient, err := disk.NewCachedDiscoveryClientForConfig(
		restConfig,
		dir,
		dir,
		180*time.Second,
	)

	if err != nil {
		return nil, errors.Wrap(err, "create cached discovery client")
	}

	restMapper := restmapper.NewDeferredDiscoveryRESTMapper(cacheDiscoverClient)

	c := &Cluster{
		clientConfig:     clientConfig,
		restConfig:       restConfig,
		logger:           logger,
		kubernetesClient: kubernetesClient,
		dynamicClient:    dynamicClient,
		discoveryClient:  discoveryClient,
		restMapper:       restMapper,
		defaultNamespace: defaultNamespace,
	}

	ctx, cancel := context.WithCancel(ctx)
	c.closeFn = cancel

	go func() {
		<-ctx.Done()
		logger.Infof("removeing cluster client temporary directory")

		if err := os.RemoveAll(dir); err != nil {
			logger.WithErr(err).Errorf("closing temporary directory")
		}
	}()

	return c, nil
}

func (c *Cluster) DefaultNamespace() string {
	return c.defaultNamespace
}

func (c *Cluster) ResourceExists(gvr schema.GroupVersionResource) bool {
	_, err := c.restMapper.KindFor(gvr)
	return err == nil
}

func (c *Cluster) Resource(gk schema.GroupKind) (schema.GroupVersionResource, error) {
	restMapping, err := c.restMapper.RESTMapping(gk)
	if err != nil {
		return schema.GroupVersionResource{}, err
	}

	return restMapping.Resource, nil
}

func (c *Cluster) KubernetesClient() (kubernetes.Interface, error) {
	return c.kubernetesClient, nil
}

func (c *Cluster) DynamicClient() (dynamic.Interface, error) {
	return c.dynamicClient, nil
}

func (c *Cluster) DiscoveryClient() (discovery.DiscoveryInterface, error) {
	return c.discoveryClient, nil
}

func (c *Cluster) NamespaceClient() (NamespaceInterface, error) {
	dc, err := c.DynamicClient()
	if err != nil {
		return nil, err
	}

	ns, _, err := c.clientConfig.Namespace()
	if err != nil {
		return nil, errors.Wrap(err, "resolving initial namespace")
	}
	return newNamespaceClient(dc, ns), nil
}

func (c *Cluster) InfoClient() (InfoInterface, error) {
	return newClusterInfo(c.clientConfig), nil
}

func (c *Cluster) Close() {
	if c.closeFn != nil {
		c.closeFn()
	}
}

func (c *Cluster) RestClient() (rest.Interface, error) {
	return rest.RESTClientFor(c.restConfig)
}

func (c *Cluster) RestConfig() *rest.Config {
	return c.restConfig
}

// Version returns a ServerVersion for the cluster.
func (c *Cluster) Version() (string, error) {
	dc, err := c.DiscoveryClient()
	if err != nil {
		return "", err
	}
	serverVersion, err := dc.ServerVersion()
	if err != nil {
		return "", err
	}
	return fmt.Sprint(serverVersion), nil
}

// FromKubeConfig creates a Cluster from a kubeConfig.
func FromKubeConfig(ctx context.Context, kubeConfig, contextName string, options RestConfigOptions) (*Cluster, error) {
	chain := strings.Deduplicate(filepath.SplitList(kubeConfig))

	rules := &clientcmd.ClientConfigLoadingRules{
		Precedence: chain,
	}

	overrides := &clientcmd.ConfigOverrides{}
	if contextName != "" {
		overrides.CurrentContext = contextName
	}
	cc := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, overrides)
	config, err := cc.ClientConfig()
	if err != nil {
		return nil, err
	}

	defaultNamespace, _, err := cc.Namespace()
	if err != nil {
		return nil, err
	}

	logger := log.From(ctx)
	logger.With("client-qps", options.Qps, "client-burst", options.Burst).
		Debugf("initializing REST client configuration")

	config = withConfigDefaults(config, options)

	return newCluster(ctx, cc, config, defaultNamespace)
}

// withConfigDefaults returns an extended rest.Config object with additional defaults applied
// See core_client.go#setConfigDefaults
func withConfigDefaults(inConfig *rest.Config, options RestConfigOptions) *rest.Config {
	config := rest.CopyConfig(inConfig)
	config.QPS = options.Qps
	config.Burst = options.Burst
	config.APIPath = "/api"
	if config.GroupVersion == nil || config.GroupVersion.Group != scheme.Scheme.PrioritizedVersionsForGroup("")[0].Group {
		gv := scheme.Scheme.PrioritizedVersionsForGroup("")[0]
		config.GroupVersion = &gv
	}
	codec := runtime.NoopEncoder{Decoder: scheme.Codecs.UniversalDecoder()}
	config.NegotiatedSerializer = serializer.NegotiatedSerializerWrapper(runtime.SerializerInfo{Serializer: codec})

	return config
}
