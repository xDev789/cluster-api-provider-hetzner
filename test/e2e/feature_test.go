/*
Copyright 2022 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package e2e

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("[Feature] Testing Cluster Flavor with 3x control-planes 1x worker", func() {
	ctx := context.TODO()

	Context("Testing Extra Services on Loadbalancer", func() {
		CaphClusterDeploymentSpec(ctx, func() CaphClusterDeploymentSpecInput {
			return CaphClusterDeploymentSpecInput{
				E2EConfig:                e2eConfig,
				ClusterctlConfigPath:     clusterctlConfigPath,
				BootstrapClusterProxy:    bootstrapClusterProxy,
				ArtifactFolder:           artifactFolder,
				SkipCleanup:              skipCleanup,
				ControlPlaneMachineCount: 3,
				WorkerMachineCount:       1,
				Flavor:                   "hcloud-feature-load-balancer-extra-services",
			}
		})
	})

	Context("Testing Placement Groups", func() {
		CaphClusterDeploymentSpec(ctx, func() CaphClusterDeploymentSpecInput {
			return CaphClusterDeploymentSpecInput{
				E2EConfig:                e2eConfig,
				ClusterctlConfigPath:     clusterctlConfigPath,
				BootstrapClusterProxy:    bootstrapClusterProxy,
				ArtifactFolder:           artifactFolder,
				SkipCleanup:              skipCleanup,
				ControlPlaneMachineCount: 3,
				WorkerMachineCount:       1,
				Flavor:                   "hcloud-feature-placement-groups",
			}
		})
	})

	// TODO: If deactivated it's necessary to set a domain name. Currently this is not supported on the CI
	// Context("Testing deactivated loadbalancer", func() {
	// 	CaphClusterDeploymentSpec(ctx, func() CaphClusterDeploymentSpecInput {
	// 		return CaphClusterDeploymentSpecInput{
	// 			E2EConfig:                e2eConfig,
	// 			ClusterctlConfigPath:     clusterctlConfigPath,
	// 			BootstrapClusterProxy:    bootstrapClusterProxy,
	// 			ArtifactFolder:           artifactFolder,
	// 			SkipCleanup:              skipCleanup,
	// 			ControlPlaneMachineCount: 3,
	// 			WorkerMachineCount:       1,
	// 			Flavor:                   "hcloud-feature-loadbalancer-off",
	// 		}
	// 	})
	// })
})
