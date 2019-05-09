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

// GetSolrDeployment will return the solr deployment
func (c *deploymentVersion) GetRc() *components.ReplicationController {


	containerConfig, ok := c.Containers["blackduck-solr"]
	if !ok {
		return nil
	}

	solrVolumeMount := c.getSolrVolumeMounts()
	solrContainerConfig := &util.Container{
		ContainerConfig: &horizonapi.ContainerConfig{Name: "solr", Image: containerConfig.Image,
			PullPolicy: horizonapi.PullAlways, MinMem: fmt.Sprintf("%dM", containerConfig.MinMem), MaxMem: fmt.Sprintf("%dM", containerConfig.MaxMem), MinCPU: fmt.Sprintf("%d", containerConfig.MinCPU), MaxCPU: fmt.Sprintf("%d", containerConfig.MinCPU)},
		EnvConfigs:   []*horizonapi.EnvConfig{c.GetHubConfigEnv()},
		VolumeMounts: solrVolumeMount,
		PortConfig:   []*horizonapi.PortConfig{{ContainerPort: "8983", Protocol: horizonapi.ProtocolTCP}},
	}

	if c.blackduck.Spec.LivenessProbes {
		solrContainerConfig.LivenessProbeConfigs = []*horizonapi.ProbeConfig{{
			ActionConfig:    horizonapi.ActionConfig{Command: []string{"/usr/local/bin/docker-healthcheck.sh", "http://localhost:8983/solr/project/admin/ping?wt=json"}},
			Delay:           240,
			Interval:        30,
			Timeout:         10,
			MinCountFailure: 10,
		}}
	}

	var initContainers []*util.Container
	if c.blackduck.Spec.PersistentStorage {
		initContainerConfig := &util.Container{
			ContainerConfig: &horizonapi.ContainerConfig{Name: "alpine", Image: "alpine", Command: []string{"sh", "-c", "chmod -cR 777 /opt/blackduck/hub/solr/cores.data"}},
			VolumeMounts:    solrVolumeMount,
		}
		initContainers = append(initContainers, initContainerConfig)
	}

	//c.PostEditContainer(solrContainerConfig)

	solr := util.CreateReplicationControllerFromContainer(&horizonapi.ReplicationControllerConfig{Namespace: c.Namespace, Name: "solr", Replicas: util.IntToInt32(1)}, "",
		[]*util.Container{solrContainerConfig}, c.getSolrVolumes(), initContainers,
		[]horizonapi.AffinityConfig{}, utils.GetVersionLabel("solr", c.blackduck.Spec.Version), utils.GetLabel("solr"), c.PullSecret)

	return solr
}

// getSolrVolumes will return the solr volumes
func (c *deploymentVersion) getSolrVolumes() []*components.Volume {
	var solrVolume *components.Volume
	if c.blackduck.Spec.PersistentStorage {
		solrVolume, _ = util.CreatePersistentVolumeClaimVolume("dir-solr", "blackduck-solr")
	} else {
		solrVolume, _ = util.CreateEmptyDirVolumeWithoutSizeLimit("dir-solr")
	}

	volumes := []*components.Volume{solrVolume}
	return volumes
}

// getSolrVolumeMounts will return the solr volume mounts
func (c *deploymentVersion) getSolrVolumeMounts() []*horizonapi.VolumeMountConfig {
	volumesMounts := []*horizonapi.VolumeMountConfig{
		{Name: "dir-solr", MountPath: "/opt/blackduck/hub/solr/cores.data"},
	}
	return volumesMounts
}

