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
