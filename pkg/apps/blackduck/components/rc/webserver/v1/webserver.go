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
	"github.com/blackducksoftware/synopsys-operator/pkg/api/blackduck/v1"
	components2 "github.com/blackducksoftware/synopsys-operator/pkg/apps/blackduck/components"
	utils2 "github.com/blackducksoftware/synopsys-operator/pkg/apps/blackduck/components/rc/utils"
	"github.com/blackducksoftware/synopsys-operator/pkg/apps/blackduck/components/utils"
	"github.com/blackducksoftware/synopsys-operator/pkg/apps/blackduck/types"
	"github.com/blackducksoftware/synopsys-operator/pkg/util"

	opc "github.com/blackducksoftware/synopsys-operator/pkg/apps/blackduck/components/rc"

)

type deploymentVersion struct {
	*opc.ReplicationController
	blackduck *v1.Blackduck
}

func NewDeploymentVersion(component *opc.ReplicationController, blackduck *v1.Blackduck) types.ReplicationControllerInterface {
	return &deploymentVersion{ReplicationController: component, blackduck: blackduck}
}

func init() {
	components2.Register(types.RcWebserverV1, NewDeploymentVersion)
}


// GetWebserverDeployment will return the webserver deployment
func (c *deploymentVersion) GetRc() *components.ReplicationController {

	containerConfig, ok := c.Containers["blackduck-nginx"]
	if !ok {
		return nil
	}


	webServerContainerConfig := &util.Container{
		ContainerConfig: &horizonapi.ContainerConfig{Name: "webserver", Image: containerConfig.Image,
			PullPolicy: horizonapi.PullAlways, MinMem: fmt.Sprintf("%dM", containerConfig.MinMem), MaxMem: fmt.Sprintf("%dM",containerConfig.MaxMem), MinCPU: fmt.Sprintf("%d", containerConfig.MinCPU), MaxCPU:fmt.Sprintf("%d", containerConfig.MaxCPU)},
		EnvConfigs:   []*horizonapi.EnvConfig{utils2.GetHubConfigEnv()},
		VolumeMounts: c.getWebserverVolumeMounts(),
		PortConfig:   []*horizonapi.PortConfig{{ContainerPort: "8443", Protocol: horizonapi.ProtocolTCP}},
	}

	if c.LivenessProbes {
		webServerContainerConfig.LivenessProbeConfigs = []*horizonapi.ProbeConfig{{
			ActionConfig:    horizonapi.ActionConfig{Command: []string{"/usr/local/bin/docker-healthcheck.sh", "https://localhost:8443/health-checks/liveness", "/tmp/secrets/WEBSERVER_CUSTOM_CERT_FILE"}},
			Delay:           180,
			Interval:        30,
			Timeout:         10,
			MinCountFailure: 10,
		}}
	}

	//c.PostEditContainer(webServerContainerConfig)

	webserver := util.CreateReplicationControllerFromContainer(&horizonapi.ReplicationControllerConfig{
		Namespace: c.Namespace, Name: "webserver", Replicas: util.IntToInt32(1)}, "",
		[]*util.Container{webServerContainerConfig}, c.getWebserverVolumes(), []*util.Container{}, []horizonapi.AffinityConfig{}, utils.GetVersionLabel("webserver", c.blackduck.Spec.Version), utils.GetLabel("webserver"), c.PullSecret)
	// log.Infof("webserver : %v\n", webserver.GetObj())
	return webserver
}

// getWebserverVolumes will return the authentication volumes
func (c *deploymentVersion) getWebserverVolumes() []*components.Volume {
	webServerEmptyDir, _ := util.CreateEmptyDirVolumeWithoutSizeLimit("dir-webserver")
	webServerSecretVol, _ := util.CreateSecretVolume("certificate", "blackduck-certificate", 0777)

	volumes := []*components.Volume{webServerEmptyDir, webServerSecretVol}

	// Custom CA auth
	if len(c.blackduck.Spec.AuthCustomCA) > 1 {
		authCustomCaVolume, _ := util.CreateSecretVolume("auth-custom-ca", "auth-custom-ca", 0777)
		volumes = append(volumes, authCustomCaVolume)
	}
	return volumes
}

// getWebserverVolumeMounts will return the authentication volume mounts
func (c *deploymentVersion) getWebserverVolumeMounts() []*horizonapi.VolumeMountConfig {
	volumesMounts := []*horizonapi.VolumeMountConfig{
		{Name: "dir-webserver", MountPath: "/opt/blackduck/hub/webserver/security"},
		{Name: "certificate", MountPath: "/tmp/secrets/WEBSERVER_CUSTOM_CERT_FILE", SubPath: "WEBSERVER_CUSTOM_CERT_FILE"},
		{Name: "certificate", MountPath: "/tmp/secrets/WEBSERVER_CUSTOM_KEY_FILE", SubPath: "WEBSERVER_CUSTOM_KEY_FILE"},
	}

	if len(c.blackduck.Spec.AuthCustomCA) > 1 {
		volumesMounts = append(volumesMounts, &horizonapi.VolumeMountConfig{
			Name:      "auth-custom-ca",
			MountPath: "/tmp/secrets/AUTH_CUSTOM_CA",
			SubPath:   "AUTH_CUSTOM_CA",
		})
	}

	return volumesMounts
}

// GetWebServerNodePortService will return the webserver nodeport service
func (c *deploymentVersion) GetWebServerNodePortService() *components.Service {
	return util.CreateService("webserver-exposed", utils.GetLabel("webserver"), c.Namespace, "443", "8443", horizonapi.ClusterIPServiceTypeNodePort, utils.GetLabel("webserver-exposed"))
}

// GetWebServerLoadBalancerService will return the webserver loadbalancer service
func (c *deploymentVersion) GetWebServerLoadBalancerService() *components.Service {
	return util.CreateService("webserver-exposed", utils.GetLabel("webserver"), c.Namespace, "443", "8443", horizonapi.ClusterIPServiceTypeLoadBalancer, utils.GetLabel("webserver-exposed"))
}
