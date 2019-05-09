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
	opc "github.com/blackducksoftware/synopsys-operator/pkg/apps/blackduck/components/rc"
	"github.com/blackducksoftware/synopsys-operator/pkg/apps/blackduck/components/utils"
	"github.com/blackducksoftware/synopsys-operator/pkg/util"
)


type deploymentVersion struct {
	*opc.ReplicationController
	blackduck *v1.Blackduck
}

func NewDeploymentVersion(replicationController *opc.ReplicationController, blackduck *v1.Blackduck) opc.ReplicationControllerInterface {
	return &deploymentVersion{ReplicationController: replicationController, blackduck: blackduck}
}


// GetRabbitmqDeployment will return the rabbitmq deployment
func (c *deploymentVersion) GetRc() *components.ReplicationController {

	containerConfig, ok := c.Containers["blackduck-rabbitmq"]
	if !ok {
		return nil
	}

	volumeMounts := c.getRabbitmqVolumeMounts()

	rabbitmqContainerConfig := &util.Container{
		ContainerConfig: &horizonapi.ContainerConfig{Name: "rabbitmq", Image: containerConfig.Image,
			PullPolicy: horizonapi.PullAlways, MinMem: fmt.Sprintf("%dM", containerConfig.MinMem), MaxMem: fmt.Sprintf("%dM", containerConfig.MaxMem), MinCPU: fmt.Sprintf("%d", containerConfig.MinCPU), MaxCPU: fmt.Sprintf("%d", containerConfig.MinCPU)},
		EnvConfigs:   []*horizonapi.EnvConfig{c.GetHubConfigEnv()},
		VolumeMounts: volumeMounts,
		PortConfig:   []*horizonapi.PortConfig{{ContainerPort: "5671", Protocol: horizonapi.ProtocolTCP}},
	}

	var initContainers []*util.Container
	if c.blackduck.Spec.PersistentStorage {
		initContainerConfig := &util.Container{
			ContainerConfig: &horizonapi.ContainerConfig{Name: "alpine", Image: "alpine", Command: []string{"sh", "-c", "chmod -cR 777 /var/lib/rabbitmq"}},
			VolumeMounts:    volumeMounts,
		}
		initContainers = append(initContainers, initContainerConfig)
	}

	//c.PostEditContainer(rabbitmqContainerConfig)

	rabbitmq := util.CreateReplicationControllerFromContainer(&horizonapi.ReplicationControllerConfig{Namespace: c.Namespace,
		Name: "rabbitmq", Replicas: util.IntToInt32(1)}, "", []*util.Container{rabbitmqContainerConfig}, c.getRabbitmqVolumes(), initContainers,
		[]horizonapi.AffinityConfig{}, utils.GetVersionLabel("rabbitmq", c.blackduck.Spec.Version), utils.GetLabel("rabbitmq"), c.PullSecret)

	return rabbitmq
}

// getRabbitmqVolumes will return the rabbitmq volumes
func (c *deploymentVersion) getRabbitmqVolumes() []*components.Volume {
	rabbitmqSecurityEmptyDir, _ := util.CreateEmptyDirVolumeWithoutSizeLimit("dir-rabbitmq-security")
	var rabbitmqDataEmptyDir *components.Volume
	if c.blackduck.Spec.PersistentStorage {
		rabbitmqDataEmptyDir, _ = util.CreatePersistentVolumeClaimVolume("dir-rabbitmq-data", "blackduck-rabbitmq")
	} else {
		rabbitmqDataEmptyDir, _ = util.CreateEmptyDirVolumeWithoutSizeLimit("dir-rabbitmq-data")
	}
	volumes := []*components.Volume{rabbitmqSecurityEmptyDir, rabbitmqDataEmptyDir}
	return volumes
}

// getRabbitmqVolumeMounts will return the rabbitmq volume mounts
func (c *deploymentVersion) getRabbitmqVolumeMounts() []*horizonapi.VolumeMountConfig {
	volumesMounts := []*horizonapi.VolumeMountConfig{
		{Name: "dir-rabbitmq-security", MountPath: "/opt/blackduck/rabbitmq/security"},
		{Name: "dir-rabbitmq-data", MountPath: "/var/lib/rabbitmq"},
	}
	return volumesMounts
}