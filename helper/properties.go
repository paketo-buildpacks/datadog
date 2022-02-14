/*
 * Copyright 2018-2020, VMware, Inc. All Rights Reserved.
 * Proprietary and Confidential.
 * Unauthorized use, copying or distribution of this source code via any medium is
 * strictly prohibited without the express written consent of VMware, Inc.
 */

package helper

import (
	"fmt"
	"strings"

	"github.com/buildpacks/libcnb"
	"github.com/paketo-buildpacks/libpak/bard"
	"github.com/paketo-buildpacks/libpak/bindings"
)

type Properties struct {
	Bindings libcnb.Bindings
	Logger   bard.Logger
}

func (p Properties) Execute() (map[string]string, error) {
	b, ok, err := bindings.ResolveOne(p.Bindings, bindings.OfType("Datadog"))
	if err != nil {
		return nil, fmt.Errorf("unable to resolve binding Datadog\n%w", err)
	} else if !ok {
		return nil, nil
	}

	p.Logger.Info("Configuring Datadog Agent properties")

	e := make(map[string]string, len(b.Secret))
	for k, v := range b.Secret {
		s := strings.ToUpper(k)
		s = strings.ReplaceAll(s, "-", "_")
		s = strings.ReplaceAll(s, ".", "_")

		e[fmt.Sprintf("DD_%s", s)] = v
	}

	return e, nil
}
