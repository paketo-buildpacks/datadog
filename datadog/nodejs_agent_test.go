/*
 * Copyright 2018-2020, VMware, Inc. All Rights Reserved.
 * Proprietary and Confidential.
 * Unauthorized use, copying or distribution of this source code via any medium is
 * strictly prohibited without the express written consent of VMware, Inc.
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
	"github.com/paketo-buildpacks/libpak/effect"
	"github.com/paketo-buildpacks/libpak/effect/mocks"
	"github.com/sclevine/spec"
	"github.com/stretchr/testify/mock"

	"github.com/paketo-buildpacks/datadog/datadog"
)

func testNodeJSAgent(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		ctx      libcnb.BuildContext
		executor *mocks.Executor
	)

	it.Before(func() {
		var err error

		ctx.Application.Path, err = ioutil.TempDir("", "nodejs-agent-application")
		Expect(err).NotTo(HaveOccurred())

		ctx.Layers.Path, err = ioutil.TempDir("", "nodejs-agent-layers")
		Expect(err).NotTo(HaveOccurred())

		executor = &mocks.Executor{}
		executor.On("Execute", mock.Anything).Return(nil)
	})

	it.After(func() {
		Expect(os.RemoveAll(ctx.Application.Path)).To(Succeed())
		Expect(os.RemoveAll(ctx.Layers.Path)).To(Succeed())
	})

	it("contributes NodeJS agent", func() {
		Expect(ioutil.WriteFile(filepath.Join(ctx.Application.Path, "package.json"), []byte(`{ "main": "main.js" }`),
			0644)).To(Succeed())
		Expect(ioutil.WriteFile(filepath.Join(ctx.Application.Path, "main.js"), []byte{}, 0644)).To(Succeed())

		dep := libpak.BuildpackDependency{
			ID:     "datadog-agent-nodejs",
			URI:    "https://localhost/stub-datadog-agent.tgz",
			SHA256: "e6417c651cc4d3fbc0ece8c715f8098106cda1a19036805fa4746db9f05b2e9a",
		}
		dc := libpak.DependencyCache{CachePath: "testdata"}

		n := datadog.NewNodeJSAgent(ctx.Application.Path, ctx.Buildpack.Path, dep, dc, bard.NewLogger(io.Discard))
		n.Executor = executor
		layer, err := ctx.Layers.Layer("test-layer")
		Expect(err).NotTo(HaveOccurred())

		layer, err = n.Contribute(layer)
		Expect(err).NotTo(HaveOccurred())

		Expect(layer.Launch).To(BeTrue())

		execution := executor.Calls[0].Arguments[0].(effect.Execution)
		Expect(execution.Command).To(Equal("npm"))
		Expect(execution.Args).To(Equal([]string{"install", "--no-save",
			filepath.Join("testdata",
				"e6417c651cc4d3fbc0ece8c715f8098106cda1a19036805fa4746db9f05b2e9a",
				"stub-datadog-agent.tgz"),
		}))

		Expect(layer.LaunchEnvironment["NODE_PATH.delim"]).To(Equal(string(os.PathListSeparator)))
		Expect(layer.LaunchEnvironment["NODE_PATH.prepend"]).To(Equal(filepath.Join(layer.Path, "node_modules")))
	})

	it("requires datadog module", func() {
		Expect(ioutil.WriteFile(filepath.Join(ctx.Application.Path, "package.json"), []byte(`{ "main": "main.js" }`),
			0644)).To(Succeed())
		Expect(ioutil.WriteFile(filepath.Join(ctx.Application.Path, "main.js"), []byte("test"), 0644)).To(Succeed())

		dep := libpak.BuildpackDependency{
			ID:     "datadog-agent-nodejs",
			URI:    "https://localhost/stub-datadog-agent.tgz",
			SHA256: "e6417c651cc4d3fbc0ece8c715f8098106cda1a19036805fa4746db9f05b2e9a",
		}
		dc := libpak.DependencyCache{CachePath: "testdata"}

		n := datadog.NewNodeJSAgent(ctx.Application.Path, ctx.Buildpack.Path, dep, dc, bard.NewLogger(io.Discard))
		n.Executor = executor
		layer, err := ctx.Layers.Layer("test-layer")
		Expect(err).NotTo(HaveOccurred())

		layer, err = n.Contribute(layer)
		Expect(err).NotTo(HaveOccurred())

		Expect(ioutil.ReadFile(filepath.Join(ctx.Application.Path, "main.js"))).To(Equal(
			[]byte("require('dd-trace').init();\ntest")))
	})

	it("does not require datadog module", func() {
		Expect(ioutil.WriteFile(filepath.Join(ctx.Application.Path, "package.json"), []byte(`{ "main": "main.js" }`),
			0644)).To(Succeed())
		Expect(ioutil.WriteFile(filepath.Join(ctx.Application.Path, "main.js"),
			[]byte("test\nrequire('dd-trace').init()\ntest"), 0644)).To(Succeed())

		dep := libpak.BuildpackDependency{
			ID:     "datadog-agent-nodejs",
			URI:    "https://localhost/stub-datadog-agent.tgz",
			SHA256: "e6417c651cc4d3fbc0ece8c715f8098106cda1a19036805fa4746db9f05b2e9a",
		}
		dc := libpak.DependencyCache{CachePath: "testdata"}

		n := datadog.NewNodeJSAgent(ctx.Application.Path, ctx.Buildpack.Path, dep, dc, bard.NewLogger(io.Discard))
		n.Executor = executor
		layer, err := ctx.Layers.Layer("test-layer")
		Expect(err).NotTo(HaveOccurred())

		layer, err = n.Contribute(layer)
		Expect(err).NotTo(HaveOccurred())

		b, err := ioutil.ReadFile(filepath.Join(ctx.Application.Path, "main.js"))
		Expect(err).ToNot(HaveOccurred())
		Expect(string(b)).To(Equal("test\nrequire('dd-trace').init()\ntest"))
	})
}
