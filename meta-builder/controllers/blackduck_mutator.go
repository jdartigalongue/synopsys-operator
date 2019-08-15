/*
 * Copyright (C) $year Synopsys, Inc.
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

/*

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	synopsysv1 "github.com/blackducksoftware/synopsys-operator/meta-builder/api/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
)

func patchBlackduck(blackduck *synopsysv1.Blackduck, objects map[string]runtime.Object) map[string]runtime.Object {
	patcher := BlackduckPatcher{
		blackduck: blackduck,
		objects:   objects,
	}
	return patcher.patch()
}

type BlackduckPatcher struct {
	blackduck *synopsysv1.Blackduck
	objects   map[string]runtime.Object
}

func (p *BlackduckPatcher) patch() map[string]runtime.Object {
	p.patchNamespace()
	p.patchStorage()
	p.patchLiveness()
	return p.objects
}

func (p *BlackduckPatcher) patchNamespace() error {
	accessor := meta.NewAccessor()
	for _, runtimeObject := range p.objects {
		accessor.SetNamespace(runtimeObject, p.blackduck.Spec.Namespace)
	}
	return nil
}

func (p *BlackduckPatcher) patchLiveness() error {
	// Removes liveness probes if Spec.LivenessProbes is set to false
	for _, v := range p.objects {
		switch v.(type) {
		case *v1.ReplicationController:
			if !p.blackduck.Spec.LivenessProbes {
				for i := range v.(*v1.ReplicationController).Spec.Template.Spec.Containers {
					v.(*v1.ReplicationController).Spec.Template.Spec.Containers[i].LivenessProbe = nil
				}
			}
		}
	}
	return nil
}

func (p *BlackduckPatcher) patchStorage() error {
	for _, v := range p.objects {
		switch v.(type) {
		case *v1.ReplicationController:
			if !p.blackduck.Spec.PersistentStorage {
				for i := range v.(*v1.ReplicationController).Spec.Template.Spec.Volumes {
					v.(*v1.ReplicationController).Spec.Template.Spec.Volumes[i].VolumeSource = v1.VolumeSource{
						EmptyDir: &v1.EmptyDirVolumeSource{
							Medium:    v1.StorageMediumDefault,
							SizeLimit: nil,
						},
					}
				}
			}
		}
	}
	return nil
}

// TODO: Create functions to patch the remaining spec fields
