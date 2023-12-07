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
	"github.com/paketo-buildpacks/libpak/bard"
)

type Build struct {
	Logger bard.Logger
}

func (b Build) Build(context libcnb.BuildContext) (libcnb.BuildResult, error) {
	b.Logger.Title(context.Buildpack)
	result := libcnb.NewBuildResult()

	pr := libpak.PlanEntryResolver{Plan: context.Plan}

	dr, err := libpak.NewDependencyResolver(context)
	if err != nil {
		return libcnb.BuildResult{}, fmt.Errorf("unable to create dependency resolver\n%w", err)
	}

	dc, err := libpak.NewDependencyCache(context)
	if err != nil {
		return libcnb.BuildResult{}, fmt.Errorf("unable to create dependency cache\n%w", err)
	}
	dc.Logger = b.Logger

	cr, err := libpak.NewConfigurationResolver(context.Buildpack, &b.Logger)
	if err != nil {
		return libcnb.BuildResult{}, fmt.Errorf("unable to create configuration resolver\n%w", err)
	}

	if _, ok, err := pr.Resolve("datadog-java"); err != nil {
		return libcnb.BuildResult{}, fmt.Errorf("unable to resolve datadog-java plan entry\n%w", err)
	} else if ok {
		agentDependency, err := dr.Resolve("datadog-agent-java", "")
		if err != nil {
			return libcnb.BuildResult{}, fmt.Errorf("unable to find dependency\n%w", err)
		}

		nativeImage := cr.ResolveBool("BP_NATIVE_IMAGE")
		result.Layers = append(result.Layers, NewJavaAgent(agentDependency, dc, b.Logger, nativeImage))

		if !nativeImage {
			h, be := libpak.NewHelperLayer(context.Buildpack, "toggle")
			h.Logger = b.Logger
			result.Layers = append(result.Layers, h)
			result.BOM.Entries = append(result.BOM.Entries, be)
		}
	}

	if _, ok, err := pr.Resolve("datadog-nodejs"); err != nil {
		return libcnb.BuildResult{}, fmt.Errorf("unable to resolve datadog-nodejs plan entry\n%w", err)
	} else if ok {
		dep, err := dr.Resolve("datadog-agent-nodejs", "")
		if err != nil {
			return libcnb.BuildResult{}, fmt.Errorf("unable to find dependency\n%w", err)
		}

		result.Layers = append(result.Layers,
			NewNodeJSAgent(context.Application.Path, context.Buildpack.Path, dep, dc, b.Logger))
	}

	return result, nil
}
