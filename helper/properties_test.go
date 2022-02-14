/*
 * Copyright 2018-2020, VMware, Inc. All Rights Reserved.
 * Proprietary and Confidential.
 * Unauthorized use, copying or distribution of this source code via any medium is
 * strictly prohibited without the express written consent of VMware, Inc.
 */

package helper_test

import (
	"testing"

	"github.com/buildpacks/libcnb"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"

	"github.com/paketo-buildpacks/datadog/helper"
)

func testProperties(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		p helper.Properties
	)

	it("does not contribute properties if no binding exists", func() {
		Expect(p.Execute()).To(BeNil())
	})

	it("contributes properties if binding exists", func() {
		p.Bindings = libcnb.Bindings{
			{
				Name:   "test-binding",
				Type:   "Datadog",
				Secret: map[string]string{"test-key": "test-value"},
			},
		}

		Expect(p.Execute()).To(Equal(map[string]string{
			"DD_TEST_KEY": "test-value",
		}))
	})
}
