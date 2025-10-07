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
