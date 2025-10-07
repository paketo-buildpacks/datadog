/*
 * Copyright 2018-2025 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package datadog

import (
	"fmt"

	"github.com/buildpacks/libcnb"
	"github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/libpak/bard"
)

type Detect struct {
	Logger bard.Logger
}

func (d Detect) Detect(context libcnb.DetectContext) (libcnb.DetectResult, error) {
	cr, err := libpak.NewConfigurationResolver(context.Buildpack, &d.Logger)
	if err != nil {
		return libcnb.DetectResult{}, fmt.Errorf("unable to create configuration resolver\n%w", err)
	}

	if !cr.ResolveBool("BP_DATADOG_ENABLED") {
		d.Logger.Info("SKIPPED: variable 'BP_DATADOG_ENABLED' not set to true")
		return libcnb.DetectResult{Pass: false}, nil
	}

	// If both BP_DATADOG_ENABLED and BP_NATIVE_IMAGE are enabled, don't require jvm-application plan and prepare for native-image only
	if cr.ResolveBool("BP_NATIVE_IMAGE") {
		d.Logger.Info("PASSED: native-image was enabled via BP_NATIVE_IMAGE")
		return libcnb.DetectResult{
			Pass: true,
			Plans: []libcnb.BuildPlan{
				{
					Provides: []libcnb.BuildPlanProvide{
						{Name: "datadog-java"},
					},
					Requires: []libcnb.BuildPlanRequire{
						{Name: "datadog-java"},
					},
				},
			},
		}, nil
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
		},
	}, nil
}
