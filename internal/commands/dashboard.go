/*
 * Copyright (c) 2019 Kubenext, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package commands

import (
	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
	"os"
)

func newDashboardCmd() *cobra.Command {
	var (
		namespace        string
		url              string
		kubeConfig       string
		verboseLevel     int
		enableOpenCensus bool
		initialContext   string
		klogVerbosity    int
		clientQps        float32
		clientBurst      int
	)

	dashboardCmd := &cobra.Command{
		Use:   "kubethan",
		Short: "kubethan kubernetes dashboard",
		Long:  "kubethan is a dashboard for high bandwidth cluster analysis operations",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}

	kubeConfig = os.Getenv("KUBECONFIG")
	if kubeConfig == "" {
		kubeConfig = clientcmd.NewDefaultClientConfigLoadingRules().GetDefaultFilename()
	}

	dashboardCmd.Flags().StringVar(&kubeConfig, "kubeconfig", kubeConfig, "absolute path to kubectl config file")
	dashboardCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "initial namespace")
	dashboardCmd.Flags().StringVar(&url, "url", "", "dashboard url")
	dashboardCmd.Flags().CountVarP(&verboseLevel, "verbosity", "v", "verbosity level")
	dashboardCmd.Flags().BoolVarP(&enableOpenCensus, "enable-opencensus", "c", false, "enable open census")
	dashboardCmd.Flags().StringVarP(&initialContext, "context", "", "", "initial context")
	dashboardCmd.Flags().IntVarP(&klogVerbosity, "klog-verbosity", "", 0, "klog verbosity level")
	dashboardCmd.Flags().Float32VarP(&clientQps, "client-qps", "", 200, "maximum QPS for client")
	dashboardCmd.Flags().IntVarP(&clientBurst, "client-burst", "", 400, "maximum burst for client throttle")

	return dashboardCmd
}
