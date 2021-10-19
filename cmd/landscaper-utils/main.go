// SPDX-FileCopyrightText: 2021 SAP SE or an SAP affiliate company and Gardener contributors.
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"
	"github.com/gardener/landscaper-utils/cmd/landscaper-utils/app"
	"os"
)

func main() {
	ctx := context.Background()
	defer ctx.Done()

	landscaperUtilsCmd := app.NewLandscaperUtilsCommand(ctx)

	if err := landscaperUtilsCmd.Execute(); err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
}
