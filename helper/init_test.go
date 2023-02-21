package helper_test

import (
	"testing"

	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestUnit(t *testing.T) {
	suite := spec.New("helper", spec.Report(report.Terminal{}))
	suite("Properties", testToggle)
	suite.Run(t)
}
