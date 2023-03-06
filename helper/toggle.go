package helper

import (
	"fmt"
	"os"
	"strings"

	"github.com/heroku/color"
	"github.com/paketo-buildpacks/libpak/bard"
	"github.com/paketo-buildpacks/libpak/sherpa"
)

type Toggle struct {
	Logger bard.Logger
}

func (t Toggle) Execute() (map[string]string, error) {
	t.Logger.Infof(color.CyanString("Datadog toggle process start ..."))
	if sherpa.ResolveBool("BPL_DATADOG_DISABLED") {
		t.Logger.Infof(color.CyanString("Datadog agent disabled by property BPL_DATADOG_DISABLED"))
		return nil, nil
	}

	var values []string
	s, ok := os.LookupEnv("JAVA_TOOL_OPTIONS")
	if s == "" {
		return nil, fmt.Errorf("disabling Datadog at launch time is unsupported for Node")
	}
	values = append(values, s)

	p, ok := os.LookupEnv("BPI_DATADOG_AGENT_PATH")
	if !ok {
		t.Logger.Infof(color.RedString("$BPI_DATADOG_AGENT_PATH is not set during build."))
		return nil, fmt.Errorf("$BPI_DATADOG_AGENT_PATH must be set")
	}
	t.Logger.Infof(color.GreenString("Datadog agent path: %s", p))

	values = append(values, fmt.Sprintf("-javaagent:%s", p))
	java_tool_options := strings.Join(values, " ")

	return map[string]string{"JAVA_TOOL_OPTIONS": java_tool_options}, nil
}
