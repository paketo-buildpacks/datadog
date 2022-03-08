/*
 * Copyright 2018-2020, VMware, Inc. All Rights Reserved.
 * Proprietary and Confidential.
 * Unauthorized use, copying or distribution of this source code via any medium is
 * strictly prohibited without the express written consent of VMware, Inc.
 */

package datadog

import (
	"fmt"

	"github.com/buildpacks/libcnb"
	"github.com/paketo-buildpacks/libpak"
)

type Detect struct{}

func (d Detect) Detect(context libcnb.DetectContext) (libcnb.DetectResult, error) {
	cr, err := libpak.NewConfigurationResolver(context.Buildpack, nil)
	if err != nil {
		return libcnb.DetectResult{}, fmt.Errorf("unable to create configuration resolver\n%w", err)
	}

	if !cr.ResolveBool("BP_DATADOG_ENABLED") {
		return libcnb.DetectResult{Pass: false}, nil
	}

	return libcnb.DetectResult{
		Pass: true,
		Plans: []libcnb.BuildPlan{
			{
				Provides: []libcnb.BuildPlanProvide{
					{Name: "datadog-java"},
				},
				Requires: []libcnb.BuildPlanRequire{
					{Name: "datadog-java"},
					{Name: "jvm-application"},
				},
			},
			{
				Provides: []libcnb.BuildPlanProvide{
					{Name: "datadog-nodejs"},
				},
				Requires: []libcnb.BuildPlanRequire{
					{Name: "datadog-nodejs"},
					{Name: "node", Metadata: map[string]interface{}{"build": true}},
				},
			},
		},
	}, nil
}
