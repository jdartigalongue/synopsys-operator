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



// GetWebappLogstashDeployment will return the webapp and logstash deployment
func (c *deploymentVersion) GetRc() *components.ReplicationController {

	containerConfigWebapp, ok := c.Containers["blackduck-webapp"]
	if !ok {
		return nil
	}

	containerConfigLogstasb, ok := c.Containers["blackduck-logstash"]
	if !ok {
		return nil
	}

	webappEnvs := []*horizonapi.EnvConfig{c.GetHubConfigEnv(), c.GetHubDBConfigEnv()}
	webappEnvs = append(webappEnvs, &horizonapi.EnvConfig{Type: horizonapi.EnvVal, NameOrPrefix: "HUB_MAX_MEMORY", KeyOrVal: fmt.Sprintf("%dM",containerConfigWebapp.MaxMem - 512)})

	webappVolumeMounts := c.getWebappVolumeMounts()

	webappContainerConfig := &util.Container{
		ContainerConfig: &horizonapi.ContainerConfig{Name: "webapp", Image: containerConfigWebapp.Image,
			PullPolicy: horizonapi.PullAlways, MinMem: fmt.Sprintf("%dM", containerConfigWebapp.MinMem), MaxMem: fmt.Sprintf("%dM", containerConfigWebapp.MaxMem), MinCPU: fmt.Sprintf("%dm", containerConfigWebapp.MinCPU), MaxCPU: fmt.Sprintf("%dm", containerConfigWebapp.MinCPU)},
		EnvConfigs:   webappEnvs,
		VolumeMounts: webappVolumeMounts,
		PortConfig:   []*horizonapi.PortConfig{{ContainerPort: "8443", Protocol: horizonapi.ProtocolTCP}},
	}

	if c.LivenessProbes {
		webappContainerConfig.LivenessProbeConfigs = []*horizonapi.ProbeConfig{{
			ActionConfig: horizonapi.ActionConfig{
				Command: []string{
					"/usr/local/bin/docker-healthcheck.sh",
					"https://127.0.0.1:8443/api/health-checks/liveness",
					"/opt/blackduck/hub/hub-webapp/security/root.crt",
					"/opt/blackduck/hub/hub-webapp/security/blackduck_system.crt",
					"/opt/blackduck/hub/hub-webapp/security/blackduck_system.key",
				},
			},
			Delay:           360,
			Interval:        30,
			Timeout:         10,
			MinCountFailure: 1000,
		}}
	}

	//c.PostEditContainer(webappContainerConfig)

	logstashVolumeMounts := c.getLogstashVolumeMounts()

	logstashContainerConfig := &util.Container{
		ContainerConfig: &horizonapi.ContainerConfig{Name: "logstash", Image: containerConfigLogstasb.Image,
			PullPolicy: horizonapi.PullAlways, MinMem: fmt.Sprintf("%dM", containerConfigLogstasb.MinMem), MaxMem: fmt.Sprintf("%dM", containerConfigLogstasb.MaxMem), MinCPU: fmt.Sprintf("%dm", containerConfigLogstasb.MinCPU), MaxCPU: fmt.Sprintf("%dm", containerConfigLogstasb.MinCPU)},
		EnvConfigs:   []*horizonapi.EnvConfig{c.GetHubConfigEnv()},
		VolumeMounts: logstashVolumeMounts,
		PortConfig:   []*horizonapi.PortConfig{{ContainerPort: "5044", Protocol: horizonapi.ProtocolTCP}},
	}

	if c.LivenessProbes {
		logstashContainerConfig.LivenessProbeConfigs = []*horizonapi.ProbeConfig{{
			ActionConfig:    horizonapi.ActionConfig{Command: []string{"/usr/local/bin/docker-healthcheck.sh", "http://localhost:9600/"}},
			Delay:           240,
			Interval:        30,
			Timeout:         10,
			MinCountFailure: 1000,
		}}
	}

	//c.PostEditContainer(logstashContainerConfig)

	var initContainers []*util.Container
	if c.blackduck.Spec.PersistentStorage {
		initContainerConfig := &util.Container{
			ContainerConfig: &horizonapi.ContainerConfig{Name: "alpine-webapp", Image: "alpine", Command: []string{"sh", "-c", "chmod -cR 777 /opt/blackduck/hub/hub-webapp/ldap"}},
			VolumeMounts:    webappVolumeMounts,
		}
		initContainers = append(initContainers, initContainerConfig)
	}
	if c.blackduck.Spec.PersistentStorage {
		initContainerConfig := &util.Container{
			ContainerConfig: &horizonapi.ContainerConfig{Name: "alpine-logstash", Image: "alpine", Command: []string{"sh", "-c", "chmod -cR 777 /var/lib/logstash/data"}},
			VolumeMounts:    logstashVolumeMounts,
		}
		initContainers = append(initContainers, initContainerConfig)
	}

	webappLogstash := util.CreateReplicationControllerFromContainer(&horizonapi.ReplicationControllerConfig{Namespace: c.Namespace, Name: "webapp-logstash", Replicas: util.IntToInt32(1)},
		"", []*util.Container{webappContainerConfig, logstashContainerConfig}, c.getWebappLogtashVolumes(),
		initContainers, []horizonapi.AffinityConfig{}, utils.GetVersionLabel("webapp-logstash", c.blackduck.Spec.Version), utils.GetLabel("webapp-logstash"), c.PullSecret)
	return webappLogstash
}

// getWebappLogtashVolumes will return the webapp and logstash volumes
func (c *deploymentVersion) getWebappLogtashVolumes() []*components.Volume {
	webappSecurityEmptyDir, _ := util.CreateEmptyDirVolumeWithoutSizeLimit("dir-webapp-security")
	var webappVolume *components.Volume
	if c.blackduck.Spec.PersistentStorage {
		webappVolume, _ = util.CreatePersistentVolumeClaimVolume("dir-webapp", "blackduck-webapp")
	} else {
		webappVolume, _ = util.CreateEmptyDirVolumeWithoutSizeLimit("dir-webapp")
	}

	var logstashVolume *components.Volume
	if c.blackduck.Spec.PersistentStorage {
		logstashVolume, _ = util.CreatePersistentVolumeClaimVolume("dir-logstash", "blackduck-logstash")
	} else {
		logstashVolume, _ = util.CreateEmptyDirVolumeWithoutSizeLimit("dir-logstash")
	}

	volumes := []*components.Volume{webappSecurityEmptyDir, webappVolume, logstashVolume, c.GetDBSecretVolume()}
	// Mount the HTTPS proxy certificate if provided
	if len(c.blackduck.Spec.ProxyCertificate) > 0 {
		volumes = append(volumes, c.GetProxyVolume())
	}

	return volumes
}

// getLogstashVolumeMounts will return the Logstash volume mounts
func (c *deploymentVersion) getLogstashVolumeMounts() []*horizonapi.VolumeMountConfig {
	volumesMounts := []*horizonapi.VolumeMountConfig{
		{Name: "dir-logstash", MountPath: "/var/lib/logstash/data"},
	}
	return volumesMounts
}

// getWebappVolumeMounts will return the Webapp volume mounts
func (c *deploymentVersion) getWebappVolumeMounts() []*horizonapi.VolumeMountConfig {
	volumesMounts := []*horizonapi.VolumeMountConfig{
		{Name: "db-passwords", MountPath: "/tmp/secrets/HUB_POSTGRES_ADMIN_PASSWORD_FILE", SubPath: "HUB_POSTGRES_ADMIN_PASSWORD_FILE"},
		{Name: "db-passwords", MountPath: "/tmp/secrets/HUB_POSTGRES_USER_PASSWORD_FILE", SubPath: "HUB_POSTGRES_USER_PASSWORD_FILE"},
		{Name: "dir-webapp", MountPath: "/opt/blackduck/hub/hub-webapp/ldap"},
		{Name: "dir-webapp-security", MountPath: "/opt/blackduck/hub/hub-webapp/security"},
		{Name: "dir-logstash", MountPath: "/opt/blackduck/hub/logs"},
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
