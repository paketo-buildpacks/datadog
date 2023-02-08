/*
 * Copyright 2018-2020, VMware, Inc. All Rights Reserved.
 * Proprietary and Confidential.
 * Unauthorized use, copying or distribution of this source code via any medium is
 * strictly prohibited without the express written consent of VMware, Inc.
 */

package datadog_test

import (
	"os"
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

	context("BP_DATADOG_ENABLE is not set", func() {
		it.Before(func() {
			Expect(os.Unsetenv("BP_DATADOG_ENABLED")).To(Succeed())
		})

		it("fails without BP_DATADOG_ENABLED", func() {
			Expect(detect.Detect(ctx)).To(Equal(libcnb.DetectResult{}))
		})
	})

	context("BP_DATADOG_ENABLE is set", func() {
		it.Before(func() {
			Expect(os.Setenv("BP_DATADOG_ENABLED", "true")).To(Succeed())
		})

		it.After(func() {
			Expect(os.Unsetenv("BP_DATADOG_ENABLED")).To(Succeed())
		})

		it("passes with BP_DATADOG_ENABLED set to true", func() {
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
							{Name: "node_modules"},
							{Name: "node", Metadata: map[string]interface{}{"build": true}},
						},
					},
				},
			}))
		})
	})

	context("BP_DATADOG_ENABLED and BP_NATIVE_IMAGE are both enabled", func() {
		it.Before(func() {
			Expect(os.Setenv("BP_DATADOG_ENABLED", "true")).To(Succeed())
			Expect(os.Setenv("BP_NATIVE_IMAGE", "true")).To(Succeed())
		})

		it.After(func() {
			Expect(os.Unsetenv("BP_DATADOG_ENABLED")).To(Succeed())
			Expect(os.Unsetenv("BP_NATIVE_IMAGE")).To(Succeed())
		})

		it("passes with native-image plan", func() {
			Expect(detect.Detect(ctx)).To(Equal(libcnb.DetectResult{
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
			}))
		})
		
	})
}
