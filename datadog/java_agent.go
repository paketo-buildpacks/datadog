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

package datadog

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/buildpacks/libcnb"
	"github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/libpak/bard"
	"github.com/paketo-buildpacks/libpak/sherpa"
)

type JavaAgent struct {
	LayerContributor libpak.DependencyLayerContributor
	Logger           bard.Logger
	NativeImage      bool
}

func NewJavaAgent(dependency libpak.BuildpackDependency, cache libpak.DependencyCache, logger bard.Logger, nativeImage bool) JavaAgent {
	contrib, _ := libpak.NewDependencyLayer(dependency, cache, libcnb.LayerTypes{
		Build:  nativeImage,
		Launch: true,
	})
	return JavaAgent{LayerContributor: contrib, Logger: logger, NativeImage: nativeImage}
}

func (j JavaAgent) Contribute(layer libcnb.Layer) (libcnb.Layer, error) {
	j.LayerContributor.Logger = j.Logger

	return j.LayerContributor.Contribute(layer, func(artifact *os.File) (libcnb.Layer, error) {
		j.Logger.Bodyf("Copying to %s", layer.Path)

		file := filepath.Join(layer.Path, filepath.Base(j.LayerContributor.Dependency.URI))
		if err := sherpa.CopyFile(artifact, file); err != nil {
			return libcnb.Layer{}, fmt.Errorf("unable to copy artifact to %s\n%w", file, err)
		}
		if err := os.Symlink(file, filepath.Join(layer.Path, "dd-java-agent.jar")); err != nil {
			return libcnb.Layer{}, fmt.Errorf("unable to create symlink to %s\n%w", file, err)
		}

		if j.NativeImage {
			layer.BuildEnvironment.Appendf("BP_NATIVE_IMAGE_BUILD_ARGUMENTS", " ", "-J-javaagent:%s", file)
		}
		layer.LaunchEnvironment.Default("BPI_DATADOG_AGENT_PATH", file)

		return layer, nil
	})
}

func (j JavaAgent) Name() string {
	return j.LayerContributor.LayerName()
}
