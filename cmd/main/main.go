/*
 * Copyright 2018-2020, VMware, Inc. All Rights Reserved.
 * Proprietary and Confidential.
 * Unauthorized use, copying or distribution of this source code via any medium is
 * strictly prohibited without the express written consent of VMware, Inc.
 */

package main

import (
	"os"

	"github.com/paketo-buildpacks/datadog/datadog"
	"github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/libpak/bard"
)

func main() {
	libpak.Main(
		datadog.Detect{},
		datadog.Build{Logger: bard.NewLogger(os.Stdout)},
	)
}
