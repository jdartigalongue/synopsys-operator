/*
 * Copyright (C) 2019 Synopsys, Inc.
 *
 *  Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements. See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership. The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 *  with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 *  under the License.
 */

package polaris

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/blackducksoftware/synopsys-operator/pkg/util"

	appsv1 "k8s.io/api/apps/v1"
	appsv1beta1 "k8s.io/api/apps/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
)

func fromYaml(content string, polaris Polaris) (map[string]runtime.Object, error) {
	// Basic
	content = strings.ReplaceAll(content, "${NAMESPACE}", polaris.Namespace)
	content = strings.ReplaceAll(content, "${VERSION}", polaris.Version)
	content = strings.ReplaceAll(content, "${ENVIRONMENT_NAME}", polaris.Namespace)
	content = strings.ReplaceAll(content, "${POLARIS_ROOT_DOMAIN}", polaris.EnvironmentDNS)
	content = strings.ReplaceAll(content, "${IMAGE_PULL_SECRETS}", polaris.ImagePullSecrets)
	content = strings.ReplaceAll(content, "${INGRESS_CLASS}", polaris.IngressClass)

	// PVC
	content = strings.ReplaceAll(content, "${DOWNLOAD_SERVER_PV_SIZE}", polaris.PolarisSpec.DownloadServerDetails.Storage.StorageSize)
	content = strings.ReplaceAll(content, "${EVENTSTORE_PV_SIZE}", polaris.PolarisDBSpec.EventstoreDetails.Storage.StorageSize)
	content = strings.ReplaceAll(content, "${MONGODB_PV_SIZE}", polaris.PolarisDBSpec.MongoDBDetails.Storage.StorageSize)
	content = strings.ReplaceAll(content, "${UPLOAD_SERVER_PV_SIZE}", polaris.PolarisDBSpec.UploadServerDetails.Storage.StorageSize)
	content = strings.ReplaceAll(content, "${POSTGRES_PV_SIZE}", polaris.PolarisDBSpec.PostgresDetails.Storage.StorageSize)
	content = strings.ReplaceAll(content, "${REPORT_STORAGE_PV_SIZE}", polaris.ReportingSpec.ReportStorageDetails.Storage.StorageSize)

	// SMTP
	content = strings.ReplaceAll(content, "${SMTP_SENDER_EMAIL}", polaris.PolarisDBSpec.SMTPDetails.SenderEmail)
	content = strings.ReplaceAll(content, "${SMTP_HOST}", polaris.PolarisDBSpec.SMTPDetails.Host)
	content = strings.ReplaceAll(content, "${SMTP_PORT}", fmt.Sprintf("\"%d\"", polaris.PolarisDBSpec.SMTPDetails.Port))

	content = strings.ReplaceAll(content, "${SMTP_TLS_MODE}", string(polaris.PolarisDBSpec.SMTPDetails.TLSMode))
	content = strings.ReplaceAll(content, "${SMTP_TLS_CHECK_SERVER_IDENTITY}", fmt.Sprintf("\"%s\"", strconv.FormatBool(polaris.PolarisDBSpec.SMTPDetails.TLSCheckServerIdentity)))
	content = strings.ReplaceAll(content, "${SMTP_TLS_TRUSTED_HOSTS}", fmt.Sprintf("\"%s\"", polaris.PolarisDBSpec.SMTPDetails.TLSTrustedHosts))

	if len(polaris.PolarisDBSpec.SMTPDetails.Username) != 0 {
		content = strings.ReplaceAll(content, "${SMTP_USERNAME}", util.EncodeStringToBase64(polaris.PolarisDBSpec.SMTPDetails.Username))
	} else {
		content = strings.ReplaceAll(content, "${SMTP_USERNAME}", "Cg==")
	}
	if len(polaris.PolarisDBSpec.SMTPDetails.Password) != 0 {
		content = strings.ReplaceAll(content, "${SMTP_PASSWORD}", fmt.Sprintf("\"%s\"", util.EncodeStringToBase64(polaris.PolarisDBSpec.SMTPDetails.Password)))
	} else {
		content = strings.ReplaceAll(content, "${SMTP_PASSWORD}", "Cg==")
	}

	// Postgres
	content = strings.ReplaceAll(content, "${POSTGRES_USERNAME}", util.EncodeStringToBase64(polaris.PolarisDBSpec.PostgresDetails.Username))
	content = strings.ReplaceAll(content, "${POSTGRES_PASSWORD}", util.EncodeStringToBase64(polaris.PolarisDBSpec.PostgresDetails.Password))
	content = strings.ReplaceAll(content, "${POSTGRES_SSL_MODE}", string(polaris.PolarisDBSpec.PostgresDetails.SSLMode))
	content = strings.ReplaceAll(content, "${POSTGRES_HOST}", polaris.PolarisDBSpec.PostgresDetails.Host)
	content = strings.ReplaceAll(content, "${POSTGRES_PORT}", fmt.Sprintf("\"%d\"", polaris.PolarisDBSpec.PostgresDetails.Port))

	// TODO Do we really need this?
	if polaris.PolarisDBSpec.PostgresDetails.IsInternal {
		content = strings.ReplaceAll(content, "${POSTGRES_TYPE}", "internal")
	} else {
		content = strings.ReplaceAll(content, "${POSTGRES_TYPE}", "external")
	}

	// Reporting
	if polaris.EnableReporting {
		content = strings.ReplaceAll(content, "${REPORTING_URL}", fmt.Sprintf("https://%s/reporting", polaris.EnvironmentDNS))
	} else {
		content = strings.ReplaceAll(content, "${REPORTING_URL}", "")
	}

	// Org job
	content = strings.ReplaceAll(content, "${ORG_DESCRIPTION}", polaris.OrganizationDetails.OrganizationProvisionOrganizationDescription)
	content = strings.ReplaceAll(content, "${ORG_NAME}", polaris.OrganizationDetails.OrganizationProvisionOrganizationName)
	content = strings.ReplaceAll(content, "${ADMIN_NAME}", polaris.OrganizationDetails.OrganizationProvisionAdminName)
	content = strings.ReplaceAll(content, "${ADMIN_USERNAME}", polaris.OrganizationDetails.OrganizationProvisionAdminUsername)
	content = strings.ReplaceAll(content, "${ADMIN_EMAIL}", polaris.OrganizationDetails.OrganizationProvisionAdminEmail)

	mapOfUniqueIDToBaseRuntimeObject, err := util.ConvertYamlFileToRuntimeObjects(content)
	if err != nil {
		return nil, err
	}
	mapOfUniqueIDToBaseRuntimeObject = removeTestManifests(mapOfUniqueIDToBaseRuntimeObject)

	patcher := Patcher{
		polaris:                          polaris,
		mapOfUniqueIDToBaseRuntimeObject: mapOfUniqueIDToBaseRuntimeObject,
	}

	return patcher.patch()
}

// Patcher holds the Polaris run time objects and it is having methods to patch it
type Patcher struct {
	polaris                          Polaris
	mapOfUniqueIDToBaseRuntimeObject map[string]runtime.Object
}

func (p *Patcher) patch() (map[string]runtime.Object, error) {
	patches := []func() error{
		p.patchNamespace,
		p.patchLabel,
		p.patchStorageClass,
		p.patchRegistry,
	}
	for _, f := range patches {
		err := f()
		if err != nil {
			return nil, err
		}
	}
	return p.mapOfUniqueIDToBaseRuntimeObject, nil
}

// patchNamespace will change the resource namespace
func (p *Patcher) patchNamespace() error {
	accessor := meta.NewAccessor()
	for _, runtimeObject := range p.mapOfUniqueIDToBaseRuntimeObject {
		accessor.SetNamespace(runtimeObject, p.polaris.Namespace)
	}
	return nil
}

// patchStorageClass will iterate over the runtime objects and update the storage class
func (p *Patcher) patchStorageClass() error {
	if len(p.polaris.StorageClass) > 0 {
		for k, v := range p.mapOfUniqueIDToBaseRuntimeObject {
			switch v.(type) {
			case *appsv1beta1.StatefulSet:
				for claimTemplateIndex := range p.mapOfUniqueIDToBaseRuntimeObject[k].(*appsv1beta1.StatefulSet).Spec.VolumeClaimTemplates {
					p.mapOfUniqueIDToBaseRuntimeObject[k].(*appsv1beta1.StatefulSet).Spec.VolumeClaimTemplates[claimTemplateIndex].Spec.StorageClassName = &p.polaris.StorageClass
				}
			case *appsv1.StatefulSet:
				for claimTemplateIndex := range p.mapOfUniqueIDToBaseRuntimeObject[k].(*appsv1.StatefulSet).Spec.VolumeClaimTemplates {
					p.mapOfUniqueIDToBaseRuntimeObject[k].(*appsv1.StatefulSet).Spec.VolumeClaimTemplates[claimTemplateIndex].Spec.StorageClassName = &p.polaris.StorageClass
				}
			case *corev1.PersistentVolumeClaim:
				p.mapOfUniqueIDToBaseRuntimeObject[k].(*corev1.PersistentVolumeClaim).Spec.StorageClassName = &p.polaris.StorageClass
			}
		}
	}
	return nil
}

// patchRegistry will update the image registry in the pod specs
func (p *Patcher) patchRegistry() error {
	if len(p.polaris.Registry) > 0 {
		if _, err := util.UpdateRegistry(p.mapOfUniqueIDToBaseRuntimeObject, p.polaris.Registry); err != nil {
			return err
		}
	}
	return nil
}

// patchLabel will set some useful labels
func (p *Patcher) patchLabel() error {
	accessor := meta.NewAccessor()
	for _, runtimeObject := range p.mapOfUniqueIDToBaseRuntimeObject {
		labels, err := accessor.Labels(runtimeObject)
		if err != nil {
			return err
		}
		if labels == nil {
			labels = make(map[string]string)
		}
		labels["synopsysctl.synopsys.com/polaris"] = p.polaris.Namespace
		if err := accessor.SetLabels(runtimeObject, labels); err != nil {
			return err
		}
	}
	return nil
}
