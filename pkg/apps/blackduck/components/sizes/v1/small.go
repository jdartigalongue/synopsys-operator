package v1

import (
	"github.com/blackducksoftware/synopsys-operator/pkg/apps/blackduck/types"
)

func init() {
	RegisterSize("small", map[string]*types.Size{
		"authentication": {
			Replica: 1,
			Containers: map[string]types.ContainerSize{
				"blackduck-authentication": {
					MinCPU: 0,
					MaxCPU: 0,
					MinMem: 1024,
					MaxMem: 1024,
				},
			},
		},
		"binaryscanner": {
			Replica: 1,
			Containers: map[string]types.ContainerSize{
				"blackduck-binaryscanner": {
					MinCPU: 1,
					MaxCPU: 1,
					MinMem: 2048,
					MaxMem: 2048,
				},
			},
		},
		"cfssl": {
			Replica: 1,
			Containers: map[string]types.ContainerSize{
				"blackduck-cfssl": {
					MinCPU: 0,
					MaxCPU: 0,
					MinMem: 640,
					MaxMem: 640,
				},
			},
		},
		"documentation": {
			Replica: 1,
			Containers: map[string]types.ContainerSize{
				"blackduck-documentation": {
					MinCPU: 0,
					MaxCPU: 0,
					MinMem: 512,
					MaxMem: 512,
				},
			},
		},
		"jobrunner": {
			Replica: 1,
			Containers: map[string]types.ContainerSize{
				"blackduck-jobrunner": {
					MinCPU: 1,
					MaxCPU: 1,
					MinMem: 4608,
					MaxMem: 4608,
				},
			},
		},
		"postgres": {
			Replica: 1,
			Containers: map[string]types.ContainerSize{
				"blackduck-postgres": {
					MinCPU: 1,
					MaxCPU: 1,
					MinMem: 3072,
					MaxMem: 3072,
				},
			},
		},
		"rabbitmq": {
			Replica: 1,
			Containers: map[string]types.ContainerSize{
				"blackduck-rabbitmq": {
					MinCPU: 0,
					MaxCPU: 0,
					MinMem: 1024,
					MaxMem: 1024,
				},
			},
		},
		"registration": {
			Replica: 1,
			Containers: map[string]types.ContainerSize{
				"blackduck-registration": {
					MinCPU: 1,
					MaxCPU: 1,
					MinMem: 640,
					MaxMem: 640,
				},
			},
		},
		"scan": {
			Replica: 1,
			Containers: map[string]types.ContainerSize{
				"blackduck-scan": {
					MinCPU: 1,
					MaxCPU: 1,
					MinMem: 2560,
					MaxMem: 2560,
				},
			},
		},
		"uploadcache": {
			Replica: 1,
			Containers: map[string]types.ContainerSize{
				"blackduck-uploadcache": {
					MinCPU: 0,
					MaxCPU: 0,
					MinMem: 512,
					MaxMem: 512,
				},
			},
		},
		"webapp-logstash": {
			Replica: 1,
			Containers: map[string]types.ContainerSize{
				"blackduck-webapp": {
					MinCPU: 1,
					MaxCPU: 1,
					MinMem: 2560,
					MaxMem: 2560,
				},
				"blackduck-logstash": {
					MinCPU: 1,
					MaxCPU: 1,
					MinMem: 1024,
					MaxMem: 1024,
				},
			},
		},
		"webserver": {
			Replica: 1,
			Containers: map[string]types.ContainerSize{
				"blackduck-nginx": {
					MinCPU: 0,
					MaxCPU: 0,
					MinMem: 512,
					MaxMem: 512,
				},
			},
		},
		"zookeeper": {
			Replica: 1,
			Containers: map[string]types.ContainerSize{
				"blackduck-zookeeper": {
					MinCPU: 1,
					MaxCPU: 1,
					MinMem: 640,
					MaxMem: 640,
				},
			},
		},
	})
}
