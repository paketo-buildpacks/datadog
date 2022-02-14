/*
 * Copyright 2018-2020, VMware, Inc. All Rights Reserved.
 * Proprietary and Confidential.
 * Unauthorized use, copying or distribution of this source code via any medium is
 * strictly prohibited without the express written consent of VMware, Inc.
 */

package datadog_test

import (
	"testing"

	"github.com/buildpacks/libcnb"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"

	"github.com/paketo-buildpacks/datadog/datadog"
)

func testDetect(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		ctx    libcnb.DetectContext
		detect datadog.Detect
	)

	it("fails without service", func() {
		Expect(detect.Detect(ctx)).To(Equal(libcnb.DetectResult{}))
	})

	it("passes with service", func() {
		ctx.Platform.Bindings = libcnb.Bindings{
			{Name: "test-service", Type: "Datadog"},
		}

		Expect(detect.Detect(ctx)).To(Equal(libcnb.DetectResult{
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
		}))
	})
}
