/*
Copyright (C) 2019Synopsys, Inc.

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


// GetZookeeperDeployment will return the zookeeper deployment
func (c *deploymentVersion) GetRc() *components.ReplicationController {

	containerConfig, ok := c.Containers["blackduck-zookeeper"]
	if !ok {
		return nil
	}

	volumeMounts := c.getZookeeperVolumeMounts()

	zookeeperContainerConfig := &util.Container{
		ContainerConfig: &horizonapi.ContainerConfig{Name: "zookeeper", Image: containerConfig.Image,
			PullPolicy: horizonapi.PullAlways, MinMem: fmt.Sprintf("%dM", containerConfig.MinMem), MaxMem: fmt.Sprintf("%dM",containerConfig.MaxMem), MinCPU: fmt.Sprintf("%d", containerConfig.MinCPU), MaxCPU:fmt.Sprintf("%d", containerConfig.MaxCPU)},
		EnvConfigs:   []*horizonapi.EnvConfig{c.GetHubConfigEnv()},
		VolumeMounts: volumeMounts,
		PortConfig:   []*horizonapi.PortConfig{{ContainerPort: "2181", Protocol: horizonapi.ProtocolTCP}},
	}

	if c.LivenessProbes {
		zookeeperContainerConfig.LivenessProbeConfigs = []*horizonapi.ProbeConfig{{
			ActionConfig:    horizonapi.ActionConfig{Command: []string{"zkServer.sh", "status", "/opt/blackduck/zookeeper/conf/zoo.cfg"}},
			Delay:           240,
			Interval:        30,
			Timeout:         10,
			MinCountFailure: 10,
		}}
	}

	var initContainers []*util.Container
	if c.blackduck.Spec.PersistentStorage {
		initContainerConfig := &util.Container{
			ContainerConfig: &horizonapi.ContainerConfig{Name: "alpine", Image: "alpine", Command: []string{"sh", "-c", "chmod -cR 777 /opt/blackduck/zookeeper/data && chmod -cR 777 /opt/blackduck/zookeeper/datalog"}},
			VolumeMounts:    volumeMounts,
		}
		initContainers = append(initContainers, initContainerConfig)
	}

	//c.PostEditContainer(zookeeperContainerConfig)

	zookeeper := util.CreateReplicationControllerFromContainer(&horizonapi.ReplicationControllerConfig{Namespace: c.Namespace, Name: "zookeeper", Replicas: util.IntToInt32(1)}, "",
		[]*util.Container{zookeeperContainerConfig}, c.getZookeeperVolumes(), initContainers, []horizonapi.AffinityConfig{}, utils.GetVersionLabel("zookeeper", c.blackduck.Spec.Version), utils.GetLabel("zookeeper"), c.PullSecret)

	return zookeeper
}

// getZookeeperVolumes will return the zookeeper volumes
func (c *deploymentVersion) getZookeeperVolumes() []*components.Volume {
	var zookeeperDataVolume *components.Volume
	var zookeeperDatalogVolume *components.Volume

	if c.blackduck.Spec.PersistentStorage {
		zookeeperDataVolume, _ = util.CreatePersistentVolumeClaimVolume("dir-zookeeper-data", "blackduck-zookeeper-data")
	} else {
		zookeeperDataVolume, _ = util.CreateEmptyDirVolumeWithoutSizeLimit("dir-zookeeper-data")
	}

	if c.blackduck.Spec.PersistentStorage {
		zookeeperDatalogVolume, _ = util.CreatePersistentVolumeClaimVolume("dir-zookeeper-datalog", "blackduck-zookeeper-datalog")
	} else {
		zookeeperDatalogVolume, _ = util.CreateEmptyDirVolumeWithoutSizeLimit("dir-zookeeper-datalog")
	}

	volumes := []*components.Volume{zookeeperDataVolume, zookeeperDatalogVolume}
	return volumes
}

// getZookeeperVolumeMounts will return the zookeeper volume mounts
func (c *deploymentVersion) getZookeeperVolumeMounts() []*horizonapi.VolumeMountConfig {
	volumesMounts := []*horizonapi.VolumeMountConfig{
		{Name: "dir-zookeeper-data", MountPath: "/opt/blackduck/zookeeper/data"},
		{Name: "dir-zookeeper-datalog", MountPath: "/opt/blackduck/zookeeper/datalog"},
	}
	return volumesMounts
}

