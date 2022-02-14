/*
 * Copyright 2018-2020, VMware, Inc. All Rights Reserved.
 * Proprietary and Confidential.
 * Unauthorized use, copying or distribution of this source code via any medium is
 * strictly prohibited without the express written consent of VMware, Inc.
 */

package datadog_test

import (
	"testing"

	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestUnit(t *testing.T) {
	suite := spec.New("Datadog", spec.Report(report.Terminal{}))
	suite("Build", testBuild)
	suite("Detect", testDetect)
	suite("JavaAgent", testJavaAgent)
	suite("NodeJSAgent", testNodeJSAgent)
	suite.Run(t)
}
