/*
Copyright (C) 2019 Synopsys, Inc.

Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements. See the NOTICE file
distributed with this work for additional information
regarding copyright ownership. The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License. You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied. See the License for the
specific language governing permissions and limitations
under the License.
*/

package v1

import (
	"fmt"
	horizonapi "github.com/blackducksoftware/horizon/pkg/api"
	"github.com/blackducksoftware/horizon/pkg/components"
	v1 "github.com/blackducksoftware/synopsys-operator/pkg/api/blackduck/v1"
	utils2 "github.com/blackducksoftware/synopsys-operator/pkg/apps/blackduck/components/rc/utils"
	"github.com/blackducksoftware/synopsys-operator/pkg/apps/blackduck/components/utils"
	"github.com/blackducksoftware/synopsys-operator/pkg/util"
	opc "github.com/blackducksoftware/synopsys-operator/pkg/apps/blackduck/components/rc"
)


type deploymentVersion struct {
	*opc.ReplicationController
	blackduck *v1.Blackduck
}

func NewDeploymentVersion(replicationController *opc.ReplicationController, blackduck *v1.Blackduck) opc.ReplicationControllerInterface {
	return &deploymentVersion{ReplicationController: replicationController, blackduck: blackduck}
}

// GetDocumentationDeployment will return the documentation deployment
func (c *deploymentVersion) GetRc() *components.ReplicationController {
	containerConfig, ok := c.Containers["blackduck-documentation"]
	if !ok {
		return nil
	}

	documentationEmptyDir, _ := util.CreateEmptyDirVolumeWithoutSizeLimit("dir-documentation")
	documentationContainerConfig := &util.Container{
		ContainerConfig: &horizonapi.ContainerConfig{Name: "documentation", Image: containerConfig.Image,
			PullPolicy: horizonapi.PullAlways, MinMem: fmt.Sprintf("%dM", containerConfig.MinMem), MaxMem: fmt.Sprintf("%dM", containerConfig.MaxMem), MinCPU: fmt.Sprintf("%d", containerConfig.MinCPU), MaxCPU: fmt.Sprintf("%d", containerConfig.MinCPU)},
		EnvConfigs: []*horizonapi.EnvConfig{utils2.GetHubConfigEnv()},
		VolumeMounts: []*horizonapi.VolumeMountConfig{
			{Name: "dir-documentation", MountPath: "/opt/blackduck/hub/hub-documentation/security"},
		},
		PortConfig: []*horizonapi.PortConfig{{ContainerPort: "8443", Protocol: horizonapi.ProtocolTCP}},
	}

	if c.blackduck.Spec.LivenessProbes {
		documentationContainerConfig.LivenessProbeConfigs = []*horizonapi.ProbeConfig{{
			ActionConfig:    horizonapi.ActionConfig{Command: []string{"/usr/local/bin/docker-healthcheck.sh", "https://127.0.0.1:8443/hubdoc/health-checks/liveness", "/opt/blackduck/hub/hub-documentation/security/root.crt"}},
			Delay:           240,
			Interval:        30,
			Timeout:         10,
			MinCountFailure: 10,
		}}
	}
	//c.PostEditContainer(documentationContainerConfig)

	documentation := util.CreateReplicationControllerFromContainer(&horizonapi.ReplicationControllerConfig{Namespace: c.Namespace, Name: "documentation", Replicas: util.IntToInt32(1)}, "",
		[]*util.Container{documentationContainerConfig}, []*components.Volume{documentationEmptyDir}, []*util.Container{}, []horizonapi.AffinityConfig{}, utils.GetVersionLabel("documentation",c.blackduck.Spec.Version), utils.GetLabel("documentation"), c.PullSecret)

	return documentation
}
