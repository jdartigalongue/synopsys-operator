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
	"github.com/blackducksoftware/horizon/pkg/components"
	v1 "github.com/blackducksoftware/synopsys-operator/pkg/api/blackduck/v1"
	opc "github.com/blackducksoftware/synopsys-operator/pkg/apps/blackduck/components/rc"
	"github.com/blackducksoftware/synopsys-operator/pkg/apps/blackduck/components/utils"
	"github.com/blackducksoftware/synopsys-operator/pkg/apps/database/postgres"
)


type deploymentVersion struct {
	*opc.ReplicationController
	blackduck *v1.Blackduck
}

func NewDeploymentVersion(replicationController *opc.ReplicationController, blackduck *v1.Blackduck) opc.ReplicationControllerInterface {
	return &deploymentVersion{ReplicationController: replicationController, blackduck: blackduck}
}


// GetPostgres will return the postgres object
func (c *deploymentVersion) GetRc() *components.ReplicationController {
	containerConfig, ok := c.Containers["blackduck-postgres"]
	if !ok {
		return nil
	}

	if len(containerConfig.Image) == 0 {
		containerConfig.Image = "registry.access.redhat.com/rhscl/postgresql-96-rhel7:1"
	}
	var pvcName string
	if c.blackduck.Spec.PersistentStorage {
		pvcName = "blackduck-postgres"
	}

	p :=  postgres.Postgres{
		Namespace:              c.Namespace,
		PVCName:                pvcName,
		Port:                   "5432",
		Image:                  containerConfig.Image,
		MinCPU:                 fmt.Sprintf("%d", containerConfig.MinCPU),
		MaxCPU:                 fmt.Sprintf("%d", containerConfig.MaxCPU),
		MinMemory:              fmt.Sprintf("%dM", containerConfig.MinMem),
		MaxMemory:              fmt.Sprintf("%dM", containerConfig.MaxMem),
		Database:               "blackduck",
		User:                   "blackduck",
		PasswordSecretName:     "db-creds",
		UserPasswordSecretKey:  "HUB_POSTGRES_USER_PASSWORD_FILE",
		AdminPasswordSecretKey: "HUB_POSTGRES_ADMIN_PASSWORD_FILE",
		MaxConnections:         300,
		SharedBufferInMB:       1024,
		EnvConfigMapRefs:       []string{"blackduck-db-config"},
		Labels:                 utils.GetLabel("postgres"),
	}
	return p.GetPostgresReplicationController()
}