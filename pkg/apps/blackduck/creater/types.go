package creaters

import (
	"github.com/blackducksoftware/synopsys-operator/pkg/apps/blackduck/components/clusterrolebinding"
	"github.com/blackducksoftware/synopsys-operator/pkg/apps/blackduck/components/configmaps"
	"github.com/blackducksoftware/synopsys-operator/pkg/apps/blackduck/components/pvc"
	opc "github.com/blackducksoftware/synopsys-operator/pkg/apps/blackduck/components/rc"
	"github.com/blackducksoftware/synopsys-operator/pkg/apps/blackduck/components/secrets"
	"github.com/blackducksoftware/synopsys-operator/pkg/apps/blackduck/components/serviceaccount"
	"github.com/blackducksoftware/synopsys-operator/pkg/apps/blackduck/components/services"
	"github.com/blackducksoftware/synopsys-operator/pkg/apps/blackduck/components/sizes"
	v1 "github.com/blackducksoftware/synopsys-operator/pkg/api/blackduck/v1"
	"github.com/blackducksoftware/synopsys-operator/pkg/protoform"
	"k8s.io/client-go/kubernetes"
)

type Blackduck struct {
	Rc map[string]ComponentReplicationController
	Service  []func(blackduck *v1.Blackduck) services.ServiceInterface
	Configmap func(blackduck *v1.Blackduck) configmaps.ConfigMapInterface
	PVC func(blackduck *v1.Blackduck) pvc.PVCInterface
	Secret func(blackduck *v1.Blackduck, config *protoform.Config, kubeClient *kubernetes.Clientset) secrets.SecretInterface
	Size  sizes.SizeInterface
	ServiceAccount func(blackduck *v1.Blackduck) serviceaccount.ServiceAccountInterface
	ClusterRoleBinding func(blackduck *v1.Blackduck) clusterrolebinding.ClusterRoleBindingInterface
}

//type ComponentReplicationController func(rc *opc.ReplicationController, blackduck *v1.Blackduck) opc.ReplicationControllerInterface

type TagOrImage struct {
	Tag string
	Image string
}

type ComponentReplicationController struct {
	Tags map[string]TagOrImage
	Func func(rc *opc.ReplicationController, blackduck *v1.Blackduck) opc.ReplicationControllerInterface
}