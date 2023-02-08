package helper

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/heroku/color"
	"github.com/paketo-buildpacks/libpak/bard"
	"github.com/paketo-buildpacks/libpak/sherpa"
)

type Toggle struct {
	Logger bard.Logger
}

func (t Toggle) Execute() (map[string]string, error) {
	t.Logger.Infof(color.CyanString("Datadog toggle process start..."))
	if datadogDisabled(t) {
		t.Logger.Infof(color.CyanString("Datadog agent disabled by property t_ENABLED"))
		return nil, nil
	}

	p, ok := os.LookupEnv("BPI_DATADOG_AGENT_PATH")
	if !ok {
		t.Logger.Infof(color.RedString("$BPI_DATADOG_AGENT_PATH is not set during build."))
		return nil, fmt.Errorf("$BPI_DATADOG_AGENT_PATH must be set")
	}
	t.Logger.Infof(color.GreenString("Datadog agent path: %s", p))

	var values []string
	if s, ok := os.LookupEnv("JAVA_TOOL_OPTIONS"); ok {
		values = append(values, s)
	}
	values = append(values, fmt.Sprintf("-javaagent:%s", p))
	java_tool_options := strings.Join(values, " ")
	//t.Logger.Infof(color.GreenString("[Datadog toggle] JAVA_TOOL_OPTIONS: %s", java_tool_options))

	return map[string]string{"JAVA_TOOL_OPTIONS": java_tool_options}, nil
}

func datadogDisabled(t Toggle) bool {
	val := sherpa.GetEnvWithDefault("BPL_DATADOG_ENABLED", "true")
	enabled, err := strconv.ParseBool(val)
	if err != nil {
		// enable by default, but warn if we couldn't understand something
		t.Logger.Infof("defaulting to enabling Datadog as '%s' could not be parsed as either true or false", val)
		return false
	}
	return !enabled
}
