package app

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

func NewComputeMachineImagesCommand(ctx context.Context) *cobra.Command {
	options := newOptions()

	cmd := &cobra.Command{
		Use:   "compute-machine-images",
		Short: "Computes machine images",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := options.complete(); err != nil {
				fmt.Print(err)
				return err
			}

			return options.run(ctx)
		},
	}

	options.addFlags(cmd.Flags())

	return cmd
}
