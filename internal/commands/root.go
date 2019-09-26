/*
 * Copyright (c) 2019 Kubenext, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
)

// Remove timestamp from log.
func init() {
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
}

func Execute(version, gitCommit, buildTime string) {
	rootCmd := newRoot(version, gitCommit, buildTime)
	if err := rootCmd.Execute(); err != nil {
		defer os.Exit(1)
		_, _ = fmt.Fprintln(os.Stderr, err)
	}
}

func newRoot(version, gitCommit, buildTime string) *cobra.Command {
	rootCmd := newDashboardCmd()
	rootCmd.AddCommand(newVersionCmd(version, gitCommit, buildTime))
	return rootCmd
}
