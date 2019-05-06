package creaters

import (
	"fmt"
	"github.com/blackducksoftware/horizon/pkg/components"
	"github.com/blackducksoftware/synopsys-operator/pkg/api"
	blackduckapi "github.com/blackducksoftware/synopsys-operator/pkg/api/blackduck/v1"
	util2 "github.com/blackducksoftware/synopsys-operator/pkg/blackduck/util"
	"github.com/blackducksoftware/synopsys-operator/pkg/protoform"
	"github.com/blackducksoftware/synopsys-operator/pkg/util"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"strings"

	opc "github.com/blackducksoftware/synopsys-operator/pkg/apps/blackduck/components/rc"
	"github.com/blackducksoftware/synopsys-operator/pkg/apps/blackduck/components/sizes"
	blackduckclientset "github.com/blackducksoftware/synopsys-operator/pkg/blackduck/client/clientset/versioned"
)

type Creater struct {
	Config          *protoform.Config
	KubeConfig      *rest.Config
	KubeClient      *kubernetes.Clientset
	BlackduckClient *blackduckclientset.Clientset
}

func (hc *Creater) autoRegisterHub(bdspec *blackduckapi.BlackduckSpec) error {
	// Filter the registration pod to auto register the hub using the registration key from the environment variable
	//registrationPod, err := util.FilterPodByNamePrefixInNamespace(hc.KubeClient, bdspec.Namespace, "registration")
	//if err != nil {
	//	return err
	//}
	//
	//registrationKey := bdspec.LicenseKey
	//
	//if registrationPod != nil && !strings.EqualFold(registrationKey, "") {
	//	for i := 0; i < 20; i++ {
	//		registrationPod, err := util.GetPods(hc.KubeClient, bdspec.Namespace, registrationPod.Name)
	//		if err != nil {
	//			return err
	//		}
	//
	//		// Create the exec into kubernetes pod request
	//		req := util.CreateExecContainerRequest(hc.KubeClient, registrationPod)
	//		// Exec into the kubernetes pod and execute the commands
	//		err = util.ExecContainer(hc.KubeConfig, req, []string{fmt.Sprintf(`curl -k -X POST "https://127.0.0.1:8443/registration/HubRegistration?registrationid=%s&action=activate" -k --cert /opt/blackduck/hub/hub-registration/security/blackduck_system.crt --key /opt/blackduck/hub/hub-registration/security/blackduck_system.key`, registrationKey)})
	//
	//		if err == nil {
	//			return nil
	//		}
	//		time.Sleep(10 * time.Second)
	//	}
	//}
	return fmt.Errorf("unable to register the blackduck %s", bdspec.Namespace)
}

func (hc *Creater) isBinaryAnalysisEnabled(bdspec *blackduckapi.BlackduckSpec) bool {
	for _, value := range bdspec.Environs {
		if strings.Contains(value, "USE_BINARY_UPLOADS") {
			values := strings.SplitN(value, ":", 2)
			if len(values) == 2 {
				mapValue := strings.TrimSpace(values[1])
				if strings.EqualFold(mapValue, "1") {
					return true
				}
			}
			return false
		}
	}
	return false
}

func (hc *Creater) GetPVCVolumeName(namespace string, name string) (string, error) {
	pvc, err := util.GetPVC(hc.KubeClient, namespace, name)
	if err != nil {
		return "", fmt.Errorf("unable to get pvc in %s namespace because %s", namespace, err.Error())
	}

	return pvc.Spec.VolumeName, nil
}

func (hc *Creater) GetComponents(blackduck *blackduckapi.Blackduck, bd Blackduck) *api.ComponentList {
	size := bd.Size.GetSize(blackduck.Spec.Size)

	// Replication Controllers
	var rcs []*components.ReplicationController
	for k, v := range bd.Rc {
		s, ok := size[k]
		if ok {
			rc := v.Func(hc.getReplicationController(k, &blackduck.Spec, v.Tags, s), blackduck)
			if rc != nil{
				rcs = append(rcs, rc.GetRc())
			}
		}
	}

	// Services
	var services []*components.Service
	for _, v := range bd.Service {
		services = append(services, v(blackduck).GetService())
	}

	//Configmap
	var configmaps []*components.ConfigMap
	configmaps = append(configmaps, bd.Configmap(blackduck).GetCM()...)

	//PVC
	var PVCs []*components.PersistentVolumeClaim
	PVCs = append(PVCs, bd.PVC(blackduck).GetPVCs()...)

	//Secret
	var secrets []*components.Secret
	secrets = append(secrets, bd.Secret(blackduck, hc.Config, hc.KubeClient).GetSecrets()...)

	// Service accounts
	var serviceAccounts []*components.ServiceAccount
	serviceAccounts = append(serviceAccounts, bd.ServiceAccount(blackduck).GetServiceAcccount()...)

	// Cluster Role Binding
	var clusterRoleBindings []*components.ClusterRoleBinding
	clusterRoleBindings = append(clusterRoleBindings, bd.ClusterRoleBinding(blackduck).GetClusterRoleBinding()...)

	return &api.ComponentList{
		ReplicationControllers: rcs,
		Services:               services,
		ConfigMaps:             configmaps,
		PersistentVolumeClaims: PVCs,
		Secrets:                secrets,
		ServiceAccounts:        serviceAccounts,
		ClusterRoleBindings:    clusterRoleBindings,
	}
}

func (hc *Creater) getReplicationController(rcname string, bdspec *blackduckapi.BlackduckSpec, containers map[string]TagOrImage, sizes *sizes.Size) *opc.ReplicationController {

	c := map[string]opc.Container{}

	if len(containers) == 0 {

	}
	for k, v := range containers {
		tmpContainer := opc.Container{
			Image: hc.getImageTag(k, v, bdspec),
		}
		containerSize, ok := sizes.Containers[k]
		if ok {
			tmpContainer.MinMem = containerSize.MinMem
			tmpContainer.MaxMem = containerSize.MaxMem
			tmpContainer.MinCPU = containerSize.MinCPU
			tmpContainer.MaxCPU = containerSize.MaxCPU
		}
		c[k] = tmpContainer
	}

	return &opc.ReplicationController{
		Namespace:      bdspec.Namespace,
		Replicas:       sizes.Replica,
		PullSecret:     bdspec.RegistryConfiguration.PullSecrets,
		LivenessProbes: bdspec.LivenessProbes,
		Containers:     c,
	}
}

func (hc *Creater) getImageTag(name string, tagOrImage TagOrImage, bdspec *blackduckapi.BlackduckSpec) string {
	confImageTag := hc.getFullContainerNameFromImageRegistryConf(name, bdspec.ImageRegistries)
	if len(confImageTag) > 0 {
		return confImageTag
	}

	if len(bdspec.RegistryConfiguration.Registry) > 0 && len(bdspec.RegistryConfiguration.Namespace) > 0 {
		return fmt.Sprintf("%s/%s/%s:%s", bdspec.RegistryConfiguration.Registry, bdspec.RegistryConfiguration.Namespace, name, tagOrImage.Tag)
	}

	if len(tagOrImage.Image) > 1 {
		return  tagOrImage.Image
	}

	return fmt.Sprintf("docker.io/blackducksoftware/%s:%s", name, tagOrImage.Tag)
}

func (hc *Creater) getFullContainerNameFromImageRegistryConf(baseContainer string, imageRegistries []string) string {
	//blackduckVersion := hubutils.GetHubVersion(c.hubSpec.Environs)
	for _, reg := range imageRegistries {
		// normal case: we expect registries
		if strings.Contains(reg, baseContainer) {
			_, err := util2.ParseImageString(reg)
			if err != nil {
				break
			}
			return reg
		}
	}
	return ""
}
