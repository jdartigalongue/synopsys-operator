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

// GetRegistrationDeployment will return the registration deployment
func (c *deploymentVersion) GetRc() *components.ReplicationController {

	containerConfig, ok := c.Containers["blackduck-registration"]
	if !ok {
		return nil
	}


	volumeMounts := c.getRegistrationVolumeMounts()

	registrationContainerConfig := &util.Container{
		ContainerConfig: &horizonapi.ContainerConfig{Name: "registration", Image: containerConfig.Image,
			PullPolicy: horizonapi.PullAlways, MinMem: fmt.Sprintf("%dM", containerConfig.MinMem), MaxMem: fmt.Sprintf("%dM", containerConfig.MaxMem), MinCPU: fmt.Sprintf("%d", containerConfig.MinCPU), MaxCPU: fmt.Sprintf("%d", containerConfig.MinCPU)},
		EnvConfigs:   []*horizonapi.EnvConfig{utils2.GetHubConfigEnv()},
		VolumeMounts: c.getRegistrationVolumeMounts(),
		PortConfig:   []*horizonapi.PortConfig{{ContainerPort: "8443", Protocol: horizonapi.ProtocolTCP}},
	}

	if c.blackduck.Spec.LivenessProbes {
		registrationContainerConfig.LivenessProbeConfigs = []*horizonapi.ProbeConfig{{
			ActionConfig: horizonapi.ActionConfig{
				Command: []string{
					"/usr/local/bin/docker-healthcheck.sh",
					"https://localhost:8443/registration/health-checks/liveness",
					"/opt/blackduck/hub/hub-registration/security/root.crt",
					"/opt/blackduck/hub/hub-registration/security/blackduck_system.crt",
					"/opt/blackduck/hub/hub-registration/security/blackduck_system.key",
				},
			},
			Delay:           240,
			Interval:        30,
			Timeout:         10,
			MinCountFailure: 10,
		}}
	}

	var initContainers []*util.Container
	if c.blackduck.Spec.PersistentStorage {
		initContainerConfig := &util.Container{
			ContainerConfig: &horizonapi.ContainerConfig{Name: "alpine", Image: "alpine", Command: []string{"sh", "-c", "chmod -cR 777 /opt/blackduck/hub/hub-registration/config"}},
			VolumeMounts:    volumeMounts,
		}
		initContainers = append(initContainers, initContainerConfig)
	}

	//c.PostEditContainer(registrationContainerConfig)

	registration := util.CreateReplicationControllerFromContainer(&horizonapi.ReplicationControllerConfig{Namespace: c.Namespace, Name: "registration", Replicas: util.IntToInt32(1)}, "",
		[]*util.Container{registrationContainerConfig}, c.getRegistrationVolumes(), initContainers,
		[]horizonapi.AffinityConfig{}, utils.GetVersionLabel("registration", c.blackduck.Spec.Version), utils.GetLabel("registration"), c.PullSecret)

	return registration
}

// getRegistrationVolumes will return the registration volumes
func (c *deploymentVersion) getRegistrationVolumes() []*components.Volume {
	var registrationVolume *components.Volume
	registrationSecurityEmptyDir, _ := util.CreateEmptyDirVolumeWithoutSizeLimit("dir-registration-security")

	if c.blackduck.Spec.PersistentStorage {
		registrationVolume, _ = util.CreatePersistentVolumeClaimVolume("dir-registration", "blackduck-registration")
	} else {
		registrationVolume, _ = util.CreateEmptyDirVolumeWithoutSizeLimit("dir-registration")
	}

	volumes := []*components.Volume{registrationVolume, registrationSecurityEmptyDir}

	// Mount the HTTPS proxy certificate if provided
	if len(c.blackduck.Spec.ProxyCertificate) > 0 {
		volumes = append(volumes, utils2.GetProxyVolume())
	}
	return volumes
}

// getRegistrationVolumeMounts will return the registration volume mounts
func (c *deploymentVersion) getRegistrationVolumeMounts() []*horizonapi.VolumeMountConfig {
	volumesMounts := []*horizonapi.VolumeMountConfig{
		{Name: "dir-registration", MountPath: "/opt/blackduck/hub/hub-registration/config"},
		{Name: "dir-registration-security", MountPath: "/opt/blackduck/hub/hub-registration/security"},
	}

	// Mount the HTTPS proxy certificate if provided
	if len(c.blackduck.Spec.ProxyCertificate) > 0 {
		volumesMounts = append(volumesMounts, &horizonapi.VolumeMountConfig{
			Name:      "blackduck-proxy-certificate",
			MountPath: "/tmp/secrets/HUB_PROXY_CERT_FILE",
			SubPath:   "HUB_PROXY_CERT_FILE",
		})
	}

	return volumesMounts
}