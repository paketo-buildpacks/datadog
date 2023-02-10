package helper_test

import (
	"os"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"

	"github.com/paketo-buildpacks/datadog/helper"
)

func testToggle(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		toggle = helper.Toggle{}
	)

	it.Before(func() {
		Expect(os.Setenv("BPL_DATADOG_DISABLED", "true")).To(Succeed())
	})

	it.After(func() {
		Expect(os.Unsetenv("BPL_DATADOG_DISABLED")).To(Succeed())
	})

	it("returns if $BPL_DATADOG_DISABLED is not set", func() {
		Expect(toggle.Execute()).To(BeNil())
	})

	context("$BPL_DATADOG_DISABLED", func() {
		it.Before(func() {
			Expect(os.Unsetenv("BPL_DATADOG_DISABLED")).To(Succeed())
		})

		it.After(func() {
			Expect(os.Unsetenv("BPL_DATADOG_DISABLED")).To(Succeed())
		})

		context("$JAVA_TOOL_OPTIONS", func() {
			it.Before(func() {
				Expect(os.Setenv("JAVA_TOOL_OPTIONS", "test-java-tool-options")).To(Succeed())
			})

			it.After(func() {
				Expect(os.Unsetenv("JAVA_TOOL_OPTIONS")).To(Succeed())
			})

			it("returns error if $BPI_DATADOG_AGENT_PATH is not set", func() {
				_, err := toggle.Execute()
				Expect(err).To(MatchError("$BPI_DATADOG_AGENT_PATH must be set"))
			})

			context("$BPI_DATADOG_AGENT_PATH", func() {
				it.Before(func() {
					Expect(os.Setenv("BPI_DATADOG_AGENT_PATH", "/mock/path/to/agent.jar")).To(Succeed())
				})

				it.After(func() {
					Expect(os.Unsetenv("BPI_DATADOG_AGENT_PATH")).To(Succeed())
				})

				it("contributes configuration", func() {
					Expect(toggle.Execute()).To(Equal(map[string]string{
						"JAVA_TOOL_OPTIONS": "test-java-tool-options -javaagent:/mock/path/to/agent.jar",
					}))
				})
			})
		})
	})
}
