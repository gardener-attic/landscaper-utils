// SPDX-FileCopyrightText: 2020 SAP SE or an SAP affiliate company and Gardener contributors.
//
// SPDX-License-Identifier: Apache-2.0

//go:generate go run -mod=vendor ../hack/generate-docs ../docs/reference

package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/gardener/landscaper-utils/pkg/shared/logger"
)

func NewLandscaperCliCommand(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "landscaper-cli",
		Short: "landscaper cli",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			log, err := logger.NewCliLogger()
			if err != nil {
				fmt.Println("unable to setup logger")
				fmt.Println(err.Error())
				os.Exit(1)
			}
			logger.SetLogger(log)
		},
	}

	logger.InitFlags(cmd.PersistentFlags())

	return cmd
}
