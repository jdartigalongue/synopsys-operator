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

func NewDeploymentVersion(component *opc.ReplicationController, blackduck *v1.Blackduck) opc.ReplicationControllerInterface {
	return &deploymentVersion{ReplicationController: component, blackduck: blackduck}
}


func (c *deploymentVersion) GetRc() *components.ReplicationController {
	containerConfig, ok := c.Containers["blackduck-jobrunner"]
	if !ok {
		return nil
	}

	jobRunnerEmptyDir, _ := util.CreateEmptyDirVolumeWithoutSizeLimit("dir-jobrunner")

	jobRunnerEnvs := []*horizonapi.EnvConfig{
		c.GetHubConfigEnv(),
		c.GetHubDBConfigEnv(),
	}

	jobRunnerEnvs = append(jobRunnerEnvs, &horizonapi.EnvConfig{Type: horizonapi.EnvVal, NameOrPrefix: "HUB_MAX_MEMORY", KeyOrVal: fmt.Sprintf("%dM",containerConfig.MaxMem - 512)})
	jobRunnerContainerConfig := &util.Container{
		ContainerConfig: &horizonapi.ContainerConfig{Name: "jobrunner", Image: containerConfig.Image,
			PullPolicy: horizonapi.PullAlways, MinMem: fmt.Sprintf("%dM", containerConfig.MinMem), MaxMem: fmt.Sprintf("%dM",containerConfig.MaxMem), MinCPU: fmt.Sprintf("%d", containerConfig.MinCPU), MaxCPU:fmt.Sprintf("%d", containerConfig.MaxCPU)},
		EnvConfigs: jobRunnerEnvs,
		VolumeMounts: []*horizonapi.VolumeMountConfig{
			{Name: "db-passwords", MountPath: "/tmp/secrets/HUB_POSTGRES_ADMIN_PASSWORD_FILE", SubPath: "HUB_POSTGRES_ADMIN_PASSWORD_FILE"},
			{Name: "db-passwords", MountPath: "/tmp/secrets/HUB_POSTGRES_USER_PASSWORD_FILE", SubPath: "HUB_POSTGRES_USER_PASSWORD_FILE"},
			{Name: "dir-jobrunner", MountPath: "/opt/blackduck/hub/jobrunner/security"},
		},
		PortConfig: []*horizonapi.PortConfig{{ContainerPort: "3001", Protocol: horizonapi.ProtocolTCP}},
	}

	if c.LivenessProbes {
		jobRunnerContainerConfig.LivenessProbeConfigs = []*horizonapi.ProbeConfig{{
			ActionConfig:    horizonapi.ActionConfig{Command: []string{"/usr/local/bin/docker-healthcheck.sh"}},
			Delay:           240,
			Interval:        30,
			Timeout:         10,
			MinCountFailure: 10,
		}}
	}

	//c.PostEditContainer(jobRunnerContainerConfig)
	jobRunnerVolumes := []*components.Volume{c.GetDBSecretVolume(), jobRunnerEmptyDir}

	// Mount the HTTPS proxy certificate if provided
	if len (c.blackduck.Spec.ProxyCertificate) > 0  {
		jobRunnerContainerConfig.VolumeMounts = append(jobRunnerContainerConfig.VolumeMounts, &horizonapi.VolumeMountConfig{
			Name:      "blackduck-proxy-certificate",
			MountPath: "/tmp/secrets/HUB_PROXY_CERT_FILE",
			SubPath:   "HUB_PROXY_CERT_FILE",
		})
		jobRunnerVolumes = append(jobRunnerVolumes, c.GetProxyVolume())
	}

	jobRunner := util.CreateReplicationControllerFromContainer(&horizonapi.ReplicationControllerConfig{Namespace: c.Namespace, Name: "jobrunner", Replicas: util.IntToInt32(c.Replicas)}, "",
		[]*util.Container{jobRunnerContainerConfig}, jobRunnerVolumes, []*util.Container{},
		[]horizonapi.AffinityConfig{}, utils.GetVersionLabel("jobrunner", c.blackduck.Spec.Version), utils.GetLabel("jobrunner"), c.PullSecret)
	return jobRunner
}