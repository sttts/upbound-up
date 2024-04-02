// Copyright 2024 Upbound Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ctx

import (
	"context"

	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

	"github.com/upbound/up/internal/profile"
	"github.com/upbound/up/internal/upbound"
)

func DeriveState(ctx context.Context, upCtx *upbound.Context, conf *clientcmdapi.Config) (NavigationState, error) {
	// find profile and derive controlplane from kubeconfig
	profiles, err := upCtx.Cfg.GetUpboundProfiles()
	if err != nil {
		return nil, err
	}
	name, p, ctp, err := profile.FromKubeconfig(ctx, profiles, conf)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return &Profiles{}, nil
	}

	// derive navigation state
	spaceKubeconfig, err := p.GetSpaceKubeConfig()
	if err != nil {
		return nil, err
	}
	switch {
	case ctp.Namespace != "" && ctp.Name != "":
		return &ControlPlane{
			space: Space{
				profile:    name,
				kubeconfig: spaceKubeconfig,
			},
			NamespacedName: ctp,
		}, nil
	case ctp.Namespace != "":
		return &Group{
			space: Space{
				profile:    name,
				kubeconfig: spaceKubeconfig,
			},
			name: ctp.Namespace,
		}, nil
	default:
		return &Space{
			profile:    name,
			kubeconfig: spaceKubeconfig,
		}, nil
	}
}
