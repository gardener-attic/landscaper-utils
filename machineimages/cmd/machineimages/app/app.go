// SPDX-FileCopyrightText: 2021 SAP SE or an SAP affiliate company and Gardener contributors.
//
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"context"
	"fmt"
	"os"

	"github.com/gardener/landscaper-utils/machineimages/pkg/logger"

	"github.com/spf13/cobra"
)

func NewComputeMachineImagesCommand(ctx context.Context) *cobra.Command {
	options := newOptions()

	cmd := &cobra.Command{
		Use:   "compute-machine-images",
		Short: "Computes machine images",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			log, err := logger.NewCliLogger()
			if err != nil {
				fmt.Println("unable to setup logger")
				fmt.Println(err.Error())
				os.Exit(1)
			}
			logger.SetLogger(log)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := options.complete(); err != nil {
				fmt.Print(err)
				return err
			}

			return options.run(ctx)
		},
	}

	logger.InitFlags(cmd.PersistentFlags())
	options.addFlags(cmd.Flags())

	return cmd
}
