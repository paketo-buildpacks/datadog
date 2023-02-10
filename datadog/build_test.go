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

func testBuild(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		ctx libcnb.BuildContext
	)

	it("contributes Java agent API <= 0.6", func() {
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

	it("contributes Java agent API >= 0.7", func() {
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

	it("contributes NodeJS agent API <= 0.6", func() {
		ctx.Plan.Entries = append(ctx.Plan.Entries, libcnb.BuildpackPlanEntry{Name: "datadog-nodejs"})
		ctx.Buildpack.Metadata = map[string]interface{}{
			"dependencies": []map[string]interface{}{
				{
					"id":      "datadog-agent-nodejs",
					"version": "1.1.1",
					"stacks":  []interface{}{"test-stack-id"},
				},
			},
		}
		ctx.Buildpack.API = "0.6"
		ctx.StackID = "test-stack-id"

		result, err := datadog.Build{}.Build(ctx)
		Expect(err).NotTo(HaveOccurred())

		Expect(result.Layers).To(HaveLen(1))
		Expect(result.Layers[0].Name()).To(Equal("datadog-agent-nodejs"))
	})

	it("contributes NodeJS agent API >= 0.7", func() {
		ctx.Plan.Entries = append(ctx.Plan.Entries, libcnb.BuildpackPlanEntry{Name: "datadog-nodejs"})
		ctx.Buildpack.Metadata = map[string]interface{}{
			"dependencies": []map[string]interface{}{
				{
					"id":      "datadog-agent-nodejs",
					"version": "1.1.1",
					"stacks":  []interface{}{"test-stack-id"},
					"cpes":    []interface{}{"cpe:2.3:a:datadog-agent:nodejs-agent:1.1.1:*:*:*:*:*:*:*"},
					"purl":    "pkg:generic/datadog-agent-nodejs-agent@1.1.1?arch=amd64",
				},
			},
		}
		ctx.Buildpack.API = "0.7"
		ctx.StackID = "test-stack-id"

		result, err := datadog.Build{}.Build(ctx)
		Expect(err).NotTo(HaveOccurred())

		Expect(result.Layers).To(HaveLen(1))
		Expect(result.Layers[0].Name()).To(Equal("datadog-agent-nodejs"))
	})
}
