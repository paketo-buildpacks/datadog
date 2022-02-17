/*
 * Copyright 2018-2020, VMware, Inc. All Rights Reserved.
 * Proprietary and Confidential.
 * Unauthorized use, copying or distribution of this source code via any medium is
 * strictly prohibited without the express written consent of VMware, Inc.
 */

package main

import (
	"fmt"
	"os"

	"github.com/buildpacks/libcnb"
	"github.com/paketo-buildpacks/datadog/helper"
	"github.com/paketo-buildpacks/libpak/bard"
	"github.com/paketo-buildpacks/libpak/sherpa"
)

func main() {
	sherpa.Execute(func() error {
		var (
			err error
			p   = helper.Properties{Logger: bard.NewLogger(os.Stdout)}
		)

		p.Bindings, err = libcnb.NewBindingsFromEnvironment()
		if err != nil {
			return fmt.Errorf("unable to read bindings from environment\n%w", err)
		}

		return sherpa.Helpers(map[string]sherpa.ExecD{
			"properties": p,
		})
	})
}
