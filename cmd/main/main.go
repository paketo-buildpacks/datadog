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
	logger := bard.NewLogger(os.Stdout)
	libpak.Main(
		datadog.Detect{Logger: logger},
		datadog.Build{Logger: logger},
	)
}
