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
	"github.com/blackducksoftware/synopsys-operator/pkg/util"
)

// GetUploadCacheDeployment will return the uploadCache deployment
func (c *Creater) GetUploadCacheDeployment(imageName string) (*components.Deployment, error) {
	volumeMounts := c.getUploadCacheVolumeMounts()

	uploadCacheContainerConfig := &util.Container{
		ContainerConfig: &horizonapi.ContainerConfig{Name: "uploadcache", Image: imageName,
			PullPolicy: horizonapi.PullAlways, MinMem: c.hubContainerFlavor.UploadCacheMemoryLimit, MaxMem: c.hubContainerFlavor.UploadCacheMemoryLimit,
			MinCPU: "", MaxCPU: ""},
		EnvConfigs: []*horizonapi.EnvConfig{
			c.getHubConfigEnv(),
		},
		VolumeMounts: volumeMounts,
		PortConfig: []*horizonapi.PortConfig{{ContainerPort: uploadCachePort1, Protocol: horizonapi.ProtocolTCP},
			{ContainerPort: uploadCachePort2, Protocol: horizonapi.ProtocolTCP}},
	}

	if c.blackDuck.Spec.LivenessProbes {
		uploadCacheContainerConfig.LivenessProbeConfigs = []*horizonapi.ProbeConfig{{
			ActionConfig: horizonapi.ActionConfig{
				Type:    horizonapi.ActionTypeCommand,
				Command: []string{"curl", "--insecure", "-X", "GET", "--verbose", "http://localhost:8086/live?full=1"},
			},
			Delay:           240,
			Interval:        30,
			Timeout:         10,
			MinCountFailure: 5,
		}}
	}
	podConfig := &util.PodConfig{
		Volumes:             c.getUploadCacheVolumes(),
		Containers:          []*util.Container{uploadCacheContainerConfig},
		Labels:              c.GetVersionLabel("uploadcache"),
		NodeAffinityConfigs: c.GetNodeAffinityConfigs("uploadcache"),
	}

	if c.blackDuck.Spec.RegistryConfiguration != nil && len(c.blackDuck.Spec.RegistryConfiguration.PullSecrets) > 0 {
		podConfig.ImagePullSecrets = c.blackDuck.Spec.RegistryConfiguration.PullSecrets
	}

	if !c.config.IsOpenshift {
		podConfig.FSGID = util.IntToInt64(0)
	}

	return util.CreateDeploymentFromContainer(
		&horizonapi.DeploymentConfig{Namespace: c.blackDuck.Spec.Namespace, Name: util.GetResourceName(c.blackDuck.Name, util.BlackDuckName, "uploadcache"), Replicas: util.IntToInt32(1)},
		podConfig, c.GetLabel("uploadcache"))
}

// getUploadCacheVolumes will return the uploadCache volumes
func (c *Creater) getUploadCacheVolumes() []*components.Volume {
	uploadCacheSecurityEmptyDir, _ := util.CreateEmptyDirVolumeWithoutSizeLimit("dir-uploadcache-security")
	sealKeySecretVol, _ := util.CreateSecretVolume("dir-seal-key", util.GetResourceName(c.blackDuck.Name, util.BlackDuckName, "upload-cache"), 0444)
	var uploadCacheDataDir *components.Volume
	if c.blackDuck.Spec.PersistentStorage {
		uploadCacheDataDir, _ = util.CreatePersistentVolumeClaimVolume("dir-uploadcache-data", c.getPVCName("uploadcache-data"))
	} else {
		uploadCacheDataDir, _ = util.CreateEmptyDirVolumeWithoutSizeLimit("dir-uploadcache-data")
	}
	volumes := []*components.Volume{uploadCacheSecurityEmptyDir, uploadCacheDataDir, sealKeySecretVol}
	return volumes
}

// getUploadCacheVolumeMounts will return the uploadCache volume mounts
func (c *Creater) getUploadCacheVolumeMounts() []*horizonapi.VolumeMountConfig {
	volumesMounts := []*horizonapi.VolumeMountConfig{
		{Name: "dir-uploadcache-security", MountPath: "/opt/blackduck/hub/blackduck-upload-cache/security"},
		{Name: "dir-uploadcache-data", MountPath: "/opt/blackduck/hub/blackduck-upload-cache/uploads", SubPath: "uploads"},
		{Name: "dir-uploadcache-data", MountPath: "/opt/blackduck/hub/blackduck-upload-cache/keys", SubPath: "keys"},
		{Name: "dir-seal-key", MountPath: "/tmp/secrets"},
	}
	return volumesMounts
}

// GetUploadCacheService will return the uploadCache service
func (c *Creater) GetUploadCacheService() *components.Service {
	return util.CreateServiceWithMultiplePort(util.GetResourceName(c.blackDuck.Name, util.BlackDuckName, "uploadcache"), c.GetLabel("uploadcache"), c.blackDuck.Spec.Namespace, []int32{uploadCachePort1, uploadCachePort2},
		horizonapi.ServiceTypeServiceIP, c.GetVersionLabel("uploadcache"))
}
