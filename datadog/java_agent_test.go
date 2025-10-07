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
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/buildpacks/libcnb"
	. "github.com/onsi/gomega"
	"github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/libpak/bard"
	"github.com/paketo-buildpacks/libpak/sherpa"
	"github.com/sclevine/spec"

	"github.com/paketo-buildpacks/datadog/datadog"
)

func testJavaAgent(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		ctx libcnb.BuildContext
	)

	it.Before(func() {
		var err error

		ctx.Buildpack.Path, err = ioutil.TempDir("", "java-agent-buildpack")
		Expect(err).NotTo(HaveOccurred())

		ctx.Layers.Path, err = ioutil.TempDir("", "java-agent-layers")
		Expect(err).NotTo(HaveOccurred())
	})

	it.After(func() {
		Expect(os.RemoveAll(ctx.Buildpack.Path)).To(Succeed())
		Expect(os.RemoveAll(ctx.Layers.Path)).To(Succeed())
	})

	it("contributes Java agent", func() {
		dep := libpak.BuildpackDependency{
			ID:     "datadog-agent-java",
			URI:    "https://localhost/stub-datadog-agent.jar",
			SHA256: "ee23306ce5f7086219c1876652ed323970ebc249f21d1c79b737ac1120284bbf",
		}
		dc := libpak.DependencyCache{CachePath: "testdata"}

		j := datadog.NewJavaAgent(dep, dc, bard.NewLogger(io.Discard), false)

		layer, err := ctx.Layers.Layer("test-layer")
		Expect(err).NotTo(HaveOccurred())

		layer, err = j.Contribute(layer)
		Expect(err).NotTo(HaveOccurred())

		Expect(layer.Launch).To(BeTrue())
		Expect(filepath.Join(layer.Path, "stub-datadog-agent.jar")).To(BeARegularFile())
		Expect(sherpa.SymlinkExists(filepath.Join(layer.Path, "dd-java-agent.jar"))).To(BeTrue())
		Expect(layer.LaunchEnvironment["BPI_DATADOG_AGENT_PATH.default"]).To(Equal(filepath.Join(layer.Path, "stub-datadog-agent.jar")))
	})
}
