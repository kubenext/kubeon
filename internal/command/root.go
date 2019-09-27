/*
 * Copyright (c) 2019 Kubenext, Inc. All Rights Reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package command

import (
	"fmt"
	"os"
)

func Execute(version, gitCommit, buildTime string) {
	rootCmd := newKubeonCmd()
	versionCmd := newVersionCmd(version, gitCommit, buildTime)
	rootCmd.AddCommand(versionCmd)
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
