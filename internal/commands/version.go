/*
 * Copyright (c) 2019 Kubenext, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package commands

import (
	"fmt"
	"github.com/spf13/cobra"
)

func newVersionCmd(version, gitCommit, buildTime string) *cobra.Command {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Show version",
		Long:  "Version for Kubeon binary",
		Run: func(cmd *cobra.Command, args []string) {
			out := cmd.OutOrStdout()
			_, _ = fmt.Fprintln(out, "Version: ", version)
			_, _ = fmt.Fprintln(out, "Git commit: ", gitCommit)
			_, _ = fmt.Fprintln(out, "Build time: ", buildTime)
		},
	}
	return versionCmd
}
