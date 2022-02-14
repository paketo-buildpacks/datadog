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
	"github.com/paketo-buildpacks/libpak/bindings"
)

type Detect struct{}

func (d Detect) Detect(context libcnb.DetectContext) (libcnb.DetectResult, error) {
	if _, ok, err := bindings.ResolveOne(context.Platform.Bindings, bindings.OfType("Datadog")); err != nil {
		return libcnb.DetectResult{}, fmt.Errorf("unable to resolve binding Datadog\n%w", err)
	} else if !ok {
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
