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

// GetBinaryScannerDeployment will return the binary scanner deployment
func (c *deploymentVersion) GetRc() *components.ReplicationController {
	binaryScannerContainerConfig := &util.Container{
		ContainerConfig: &horizonapi.ContainerConfig{Name: "binaryscanner", Image: c.Image,
			PullPolicy: horizonapi.PullAlways, MinMem: fmt.Sprintf("%dM", c.MinMem), MaxMem: fmt.Sprintf("%dM", c.MaxMem), MinCPU: fmt.Sprintf("%dm", c.MinCPU), MaxCPU: fmt.Sprintf("%dm", c.MinCPU),
			Command: []string{"/docker-entrypoint.sh"}},
		EnvConfigs: []*horizonapi.EnvConfig{c.GetHubConfigEnv()},
		PortConfig: []*horizonapi.PortConfig{{ContainerPort: "3001", Protocol: horizonapi.ProtocolTCP}},
	}

	//c.PostEditContainer(binaryScannerContainerConfig)

	binaryScanner := util.CreateReplicationControllerFromContainer(&horizonapi.ReplicationControllerConfig{Namespace: c.Namespace,
		Name: "binaryscanner", Replicas: util.IntToInt32(1)}, c.Namespace, []*util.Container{binaryScannerContainerConfig},
		[]*components.Volume{}, []*util.Container{}, []horizonapi.AffinityConfig{}, c.GetVersionLabel("binaryscanner", c.blackduck.Spec.Version), c.GetLabel("binaryscanner"), c.PullSecret)
	// log.Infof("binaryScanner : %v\n", binaryScanner.GetObj())
	return binaryScanner
}
