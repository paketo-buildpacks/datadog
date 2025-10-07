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
	"testing"

	"github.com/buildpacks/libcnb"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"

	"github.com/paketo-buildpacks/datadog/datadog"
)

func testBuild(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		ctx libcnb.BuildContext
	)

	it.Before(func() {
		t.Setenv("BP_ARCH", "amd64")
	})

	it("contributes Java agent API <= 0.6 and toggle", func() {
		ctx.Plan.Entries = append(ctx.Plan.Entries, libcnb.BuildpackPlanEntry{Name: "datadog-java"})
		ctx.Buildpack.Metadata = map[string]interface{}{
			"dependencies": []map[string]interface{}{
				{
					"id":      "datadog-agent-java",
					"version": "1.1.1",
					"stacks":  []interface{}{"test-stack-id"},
				},
			},
		}
		ctx.Buildpack.API = "0.6"
		ctx.StackID = "test-stack-id"

		result, err := datadog.Build{}.Build(ctx)
		Expect(err).NotTo(HaveOccurred())

		Expect(result.Layers).To(HaveLen(2))
		Expect(result.Layers[0].Name()).To(Equal("datadog-agent-java"))
		Expect(result.Layers[1].Name()).To(Equal("helper"))
	})

	it("contributes Java agent API >= 0.7 and toggle", func() {
		ctx.Plan.Entries = append(ctx.Plan.Entries, libcnb.BuildpackPlanEntry{Name: "datadog-java"})
		ctx.Buildpack.Metadata = map[string]interface{}{
			"dependencies": []map[string]interface{}{
				{
					"id":      "datadog-agent-java",
					"version": "1.1.1",
					"stacks":  []interface{}{"test-stack-id"},
					"cpes":    []interface{}{"cpe:2.3:a:datadog-agent:java-agent:1.1.1:*:*:*:*:*:*:*"},
					"purl":    "pkg:generic/datadog-agent-java-agent@1.1.1?arch=amd64",
				},
			},
		}
		ctx.Buildpack.API = "0.7"
		ctx.StackID = "test-stack-id"

		result, err := datadog.Build{}.Build(ctx)
		Expect(err).NotTo(HaveOccurred())

		Expect(result.Layers).To(HaveLen(2))
		Expect(result.Layers[0].Name()).To(Equal("datadog-agent-java"))
		Expect(result.Layers[1].Name()).To(Equal("helper"))
	})

	context("using Native Image", func() {
		it.Before(func() {
			t.Setenv("BP_NATIVE_IMAGE", "true")
		})

		it("contributes Java agent API >= 0.7 only", func() {
			ctx.Plan.Entries = append(ctx.Plan.Entries, libcnb.BuildpackPlanEntry{Name: "datadog-java"})
			ctx.Buildpack.Metadata = map[string]interface{}{
				"dependencies": []map[string]interface{}{
					{
						"id":      "datadog-agent-java",
						"version": "1.1.1",
						"stacks":  []interface{}{"test-stack-id"},
						"cpes":    []interface{}{"cpe:2.3:a:datadog-agent:java-agent:1.1.1:*:*:*:*:*:*:*"},
						"purl":    "pkg:generic/datadog-agent-java-agent@1.1.1?arch=amd64",
					},
				},
			}
			ctx.Buildpack.API = "0.7"
			ctx.StackID = "test-stack-id"

			result, err := datadog.Build{}.Build(ctx)
			Expect(err).NotTo(HaveOccurred())

			Expect(result.Layers).To(HaveLen(1))
			Expect(result.Layers[0].Name()).To(Equal("datadog-agent-java"))
		})
	})
}
