package types

import (
	"github.com/blackducksoftware/horizon/pkg/components"
	v1 "github.com/blackducksoftware/synopsys-operator/pkg/api/blackduck/v1"
	"github.com/blackducksoftware/synopsys-operator/pkg/protoform"
	"k8s.io/client-go/kubernetes"
)


type ReplicationControllerCreater func(replicationController *ReplicationController, blackduck *v1.Blackduck) ReplicationControllerInterface
type ServiceCreater func(blackduck *v1.Blackduck) ServiceInterface
type ConfigmapCreater func(blackduck *v1.Blackduck) ConfigMapInterface
type PvcCreater func(blackduck *v1.Blackduck) PVCInterface
type SecretCreater func(blackduck *v1.Blackduck, config *protoform.Config, kubeClient *kubernetes.Clientset) SecretInterface
type ServiceAccountCreater func(blackduck *v1.Blackduck) ServiceAccountInterface
type ClusterRoleBindingCreater func(blackduck *v1.Blackduck) ClusterRoleBindingInterface

type TagOrImage struct {
	Tag string
	Image string
}

type ClusterRoleBindingInterface interface {
	GetClusterRoleBinding() []*components.ClusterRoleBinding
}

type ConfigMapInterface interface {
	GetCM() []*components.ConfigMap
}

type PVCInterface interface {
	GetPVCs() []*components.PersistentVolumeClaim
	// TODO add deployment, rc
}


type ReplicationController struct {
	Namespace        string
	Replicas         int
	PullSecret       []string
	LivenessProbes   bool
	Containers map[string]Container
}

type Container struct {
	Image            string
	MinCPU           int
	MaxCPU           int
	MinMem           int
	MaxMem           int
}

type ReplicationControllerInterface interface {
	GetRc() *components.ReplicationController
	// TODO add deployment, rc
}

type SecretInterface interface {
	GetSecrets() []*components.Secret
}


type ServiceAccountInterface interface {
	GetServiceAcccount() []*components.ServiceAccount
}

type ServiceInterface interface {
	GetService() *components.Service
	// TODO add deployment, rc
}

type SizeInterface interface {
	GetSize(name string) map[string]*Size
}

type ContainerSize struct {
	MinCPU   int
	MaxCPU   int
	MinMem   int
	MaxMem   int
}

type Size struct {
	Replica int
	Containers map[string]ContainerSize
}