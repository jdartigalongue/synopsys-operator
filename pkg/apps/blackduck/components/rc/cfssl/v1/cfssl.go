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

package containers

import (
	"fmt"
	horizonapi "github.com/blackducksoftware/horizon/pkg/api"
	"github.com/blackducksoftware/horizon/pkg/components"
	v1 "github.com/blackducksoftware/synopsys-operator/pkg/api/blackduck/v1"
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

// GetCfsslDeployment will return the cfssl deployment
func (c *deploymentVersion) GetRc() *components.ReplicationController {
	containerConfig, ok := c.Containers["blackduck-cfssl"]
	if !ok {
		return nil
	}


	cfsslVolumeMounts := c.getCfsslolumeMounts()
	cfsslContainerConfig := &util.Container{
		ContainerConfig: &horizonapi.ContainerConfig{Name: "cfssl", Image: containerConfig.Image,
			PullPolicy: horizonapi.PullAlways, MinMem: fmt.Sprintf("%dM", containerConfig.MinMem), MaxMem: fmt.Sprintf("%dM", containerConfig.MaxMem), MinCPU: fmt.Sprintf("%dm", containerConfig.MinCPU), MaxCPU: fmt.Sprintf("%dm", containerConfig.MinCPU)},
		EnvConfigs:   []*horizonapi.EnvConfig{c.GetHubConfigEnv()},
		VolumeMounts: cfsslVolumeMounts,
		PortConfig:   []*horizonapi.PortConfig{{ContainerPort: "8888", Protocol: horizonapi.ProtocolTCP}},
	}

	if c.blackduck.Spec.LivenessProbes {
		cfsslContainerConfig.LivenessProbeConfigs = []*horizonapi.ProbeConfig{{
			ActionConfig:    horizonapi.ActionConfig{Command: []string{"/usr/local/bin/docker-healthcheck.sh", "http://localhost:8888/api/v1/cfssl/scaninfo"}},
			Delay:           240,
			Interval:        30,
			Timeout:         10,
			MinCountFailure: 10,
		}}
	}

	var initContainers []*util.Container
	if c.blackduck.Spec.PersistentStorage {
		initContainerConfig := &util.Container{
			ContainerConfig: &horizonapi.ContainerConfig{Name: "alpine", Image: "alpine", Command: []string{"sh", "-c", "chmod -cR 777 /etc/cfssl"}},
			VolumeMounts:    cfsslVolumeMounts,
		}
		initContainers = append(initContainers, initContainerConfig)
	}

	//c.PostEditContainer(cfsslContainerConfig)

	cfssl := util.CreateReplicationControllerFromContainer(&horizonapi.ReplicationControllerConfig{Namespace: c.Namespace, Name: "cfssl", Replicas: util.IntToInt32(1)}, "",
		[]*util.Container{cfsslContainerConfig}, c.getCfsslVolumes(), initContainers,
		[]horizonapi.AffinityConfig{}, utils.GetVersionLabel("cfssl", c.blackduck.Spec.Version), utils.GetLabel("cfssl"), c.PullSecret)

	return cfssl
}

// getCfsslVolumes will return the cfssl volumes
func (c *deploymentVersion) getCfsslVolumes() []*components.Volume {
	var cfsslVolume *components.Volume
	if c.blackduck.Spec.PersistentStorage {
		cfsslVolume, _ = util.CreatePersistentVolumeClaimVolume("dir-cfssl", "blackduck-cfssl")
	} else {
		cfsslVolume, _ = util.CreateEmptyDirVolumeWithoutSizeLimit("dir-cfssl")
	}

	volumes := []*components.Volume{cfsslVolume}
	return volumes
}

// getCfsslolumeMounts will return the cfssl volume mounts
func (c *deploymentVersion) getCfsslolumeMounts() []*horizonapi.VolumeMountConfig {
	volumesMounts := []*horizonapi.VolumeMountConfig{
		{Name: "dir-cfssl", MountPath: "/etc/cfssl"},
	}
	return volumesMounts
}
