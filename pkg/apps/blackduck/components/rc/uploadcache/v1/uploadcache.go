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
	utils2 "github.com/blackducksoftware/synopsys-operator/pkg/apps/blackduck/components/rc/utils"
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


// GetUploadCacheDeployment will return the uploadCache deployment
func (c *deploymentVersion) GetRc() *components.ReplicationController {

	containerConfig, ok := c.Containers["blackduck-uploadcache"]
	if !ok {
		return nil
	}


	volumeMounts := c.getUploadCacheVolumeMounts()

	uploadCacheContainerConfig := &util.Container{
		ContainerConfig: &horizonapi.ContainerConfig{Name: "uploadcache", Image: containerConfig.Image,
			PullPolicy: horizonapi.PullAlways, MinMem: fmt.Sprintf("%dM", containerConfig.MinMem), MaxMem: fmt.Sprintf("%dM", containerConfig.MaxMem), MinCPU: fmt.Sprintf("%d", containerConfig.MinCPU), MaxCPU: fmt.Sprintf("%d", containerConfig.MinCPU)},
		EnvConfigs: []*horizonapi.EnvConfig{
			utils2.GetHubConfigEnv(),
			// {NameOrPrefix: "SEAL_KEY", Type: horizonapi.EnvFromSecret, KeyOrVal: "SEAL_KEY", FromName: "upload-cache"},
		},
		VolumeMounts: volumeMounts,
		PortConfig: []*horizonapi.PortConfig{{ContainerPort: "9443", Protocol: horizonapi.ProtocolTCP},
			{ContainerPort: "9444", Protocol: horizonapi.ProtocolTCP}},
	}

	if c.LivenessProbes {
		uploadCacheContainerConfig.LivenessProbeConfigs = []*horizonapi.ProbeConfig{{
			ActionConfig:    horizonapi.ActionConfig{Command: []string{"curl", "--insecure", "-X", "GET", "--verbose", "http://localhost:8086/live?full=1"}},
			Delay:           240,
			Interval:        30,
			Timeout:         10,
			MinCountFailure: 5,
		}}
	}

	var initContainers []*util.Container
	if c.blackduck.Spec.PersistentStorage {
		initContainerConfig := &util.Container{
			ContainerConfig: &horizonapi.ContainerConfig{Name: "alpine", Image: "alpine", Command: []string{"sh", "-c", "chmod -cR 777 /opt/blackduck/hub/blackduck-upload-cache/security /opt/blackduck/hub/blackduck-upload-cache/keys /opt/blackduck/hub/blackduck-upload-cache/uploads"}},
			VolumeMounts:    volumeMounts,
		}
		initContainers = append(initContainers, initContainerConfig)
	}

	//c.PostEditContainer(uploadCacheContainerConfig)

	uploadCache := util.CreateReplicationControllerFromContainer(&horizonapi.ReplicationControllerConfig{Namespace: c.Namespace,
		Name: "uploadcache", Replicas: util.IntToInt32(1)}, "", []*util.Container{uploadCacheContainerConfig}, c.getUploadCacheVolumes(),
		initContainers, []horizonapi.AffinityConfig{}, utils.GetVersionLabel("uploadcache", c.blackduck.Spec.Version), utils.GetLabel("uploadcache"), c.PullSecret)

	return uploadCache
}

// getUploadCacheVolumes will return the uploadCache volumes
func (c *deploymentVersion) getUploadCacheVolumes() []*components.Volume {
	uploadCacheSecurityEmptyDir, _ := util.CreateEmptyDirVolumeWithoutSizeLimit("dir-uploadcache-security")
	sealKeySecretVol, _ := util.CreateSecretVolume("dir-seal-key", "upload-cache", 0777)
	var uploadCacheDataDir *components.Volume
	var uploadCacheDataKey *components.Volume
	if c.blackduck.Spec.PersistentStorage {
		uploadCacheDataDir, _ = util.CreatePersistentVolumeClaimVolume("dir-uploadcache-data", "blackduck-uploadcache-data")
		uploadCacheDataKey, _ = util.CreatePersistentVolumeClaimVolume("dir-uploadcache-key", "blackduck-uploadcache-key")
	} else {
		uploadCacheDataDir, _ = util.CreateEmptyDirVolumeWithoutSizeLimit("dir-uploadcache-data")
		uploadCacheDataKey, _ = util.CreateEmptyDirVolumeWithoutSizeLimit("dir-uploadcache-key")
	}
	volumes := []*components.Volume{uploadCacheSecurityEmptyDir, uploadCacheDataDir, uploadCacheDataKey, sealKeySecretVol}
	return volumes
}

// getUploadCacheVolumeMounts will return the uploadCache volume mounts
func (c *deploymentVersion) getUploadCacheVolumeMounts() []*horizonapi.VolumeMountConfig {
	volumesMounts := []*horizonapi.VolumeMountConfig{
		{Name: "dir-uploadcache-security", MountPath: "/opt/blackduck/hub/blackduck-upload-cache/security"},
		{Name: "dir-uploadcache-data", MountPath: "/opt/blackduck/hub/blackduck-upload-cache/uploads"},
		{Name: "dir-uploadcache-key", MountPath: "/opt/blackduck/hub/blackduck-upload-cache/keys"},
		{Name: "dir-seal-key", MountPath: "/tmp/secrets"},
	}
	return volumesMounts
}
