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
	horizonapi "github.com/blackducksoftware/horizon/pkg/api"
	"github.com/blackducksoftware/horizon/pkg/components"
	apputil "github.com/blackducksoftware/synopsys-operator/pkg/apps/util"
	"github.com/blackducksoftware/synopsys-operator/pkg/util"
)

// GetCfsslDeployment will return the cfssl deployment
func (c *Creater) GetCfsslDeployment(imageName string) (*components.Deployment, error) {
	podName := "cfssl"

	cfsslVolumeMounts := c.getCfsslolumeMounts()
	cfsslContainerConfig := &util.Container{
		ContainerConfig: &horizonapi.ContainerConfig{Name: podName, Image: imageName,
			PullPolicy: horizonapi.PullAlways, MinMem: c.hubContainerFlavor.CfsslMemoryLimit, MaxMem: c.hubContainerFlavor.CfsslMemoryLimit, MinCPU: "", MaxCPU: ""},
		EnvConfigs:   []*horizonapi.EnvConfig{c.getHubConfigEnv()},
		VolumeMounts: cfsslVolumeMounts,
		PortConfig:   []*horizonapi.PortConfig{{ContainerPort: cfsslPort, Protocol: horizonapi.ProtocolTCP}},
	}

	if c.blackDuck.Spec.LivenessProbes {
		cfsslContainerConfig.LivenessProbeConfigs = []*horizonapi.ProbeConfig{{
			ActionConfig: horizonapi.ActionConfig{
				Type:    horizonapi.ActionTypeCommand,
				Command: []string{"/usr/local/bin/docker-healthcheck.sh", "http://localhost:8888/api/v1/cfssl/scaninfo"},
			},
			Delay:           240,
			Interval:        30,
			Timeout:         10,
			MinCountFailure: 10,
		}}
	}

	podConfig := &util.PodConfig{
		Volumes:             c.getCfsslVolumes(),
		Containers:          []*util.Container{cfsslContainerConfig},
		Labels:              c.GetVersionLabel(podName),
		NodeAffinityConfigs: c.GetNodeAffinityConfigs(podName),
		ServiceAccount:      util.GetResourceName(c.blackDuck.Name, util.BlackDuckName, "service-account"),
	}

	if c.blackDuck.Spec.RegistryConfiguration != nil && len(c.blackDuck.Spec.RegistryConfiguration.PullSecrets) > 0 {
		podConfig.ImagePullSecrets = c.blackDuck.Spec.RegistryConfiguration.PullSecrets
	}

	apputil.ConfigurePodConfigSecurityContext(podConfig, c.blackDuck.Spec.SecurityContexts, "blackduck-cfssl", c.config.IsOpenshift)

	return util.CreateDeploymentFromContainer(
		&horizonapi.DeploymentConfig{Namespace: c.blackDuck.Spec.Namespace, Name: util.GetResourceName(c.blackDuck.Name, util.BlackDuckName, podName), Replicas: util.IntToInt32(1)},
		podConfig, c.GetLabel(podName))
}

// getCfsslVolumes will return the cfssl volumes
func (c *Creater) getCfsslVolumes() []*components.Volume {
	var cfsslVolume *components.Volume
	if c.blackDuck.Spec.PersistentStorage {
		cfsslVolume, _ = util.CreatePersistentVolumeClaimVolume("dir-cfssl", c.getPVCName("cfssl"))
	} else {
		cfsslVolume, _ = util.CreateEmptyDirVolumeWithoutSizeLimit("dir-cfssl")
	}

	volumes := []*components.Volume{cfsslVolume}
	return volumes
}

// getCfsslolumeMounts will return the cfssl volume mounts
func (c *Creater) getCfsslolumeMounts() []*horizonapi.VolumeMountConfig {
	volumesMounts := []*horizonapi.VolumeMountConfig{
		{Name: "dir-cfssl", MountPath: "/etc/cfssl"},
	}
	return volumesMounts
}

// GetCfsslService will return the cfssl service
func (c *Creater) GetCfsslService() *components.Service {
	return util.CreateService(util.GetResourceName(c.blackDuck.Name, util.BlackDuckName, "cfssl"), c.GetLabel("cfssl"), c.blackDuck.Spec.Namespace, cfsslPort, cfsslPort, horizonapi.ServiceTypeServiceIP, c.GetVersionLabel("cfssl"))
}
