package v1

import (
	horizonapi "github.com/blackducksoftware/horizon/pkg/api"
	"github.com/blackducksoftware/horizon/pkg/components"
	v1 "github.com/blackducksoftware/synopsys-operator/pkg/api/blackduck/v1"
	"github.com/blackducksoftware/synopsys-operator/pkg/apps/blackduck/components/utils"
	"github.com/blackducksoftware/synopsys-operator/pkg/apps/blackduck/types"
)

type clusterrolebinding struct {
	blackduck *v1.Blackduck
}

func (c *clusterrolebinding) GetClusterRoleBinding() []*components.ClusterRoleBinding {
	clusterRoleBinding := components.NewClusterRoleBinding(horizonapi.ClusterRoleBindingConfig{
		Name:       "blackduck",
		APIVersion: "rbac.authorization.k8s.io/v1",
	})

	clusterRoleBinding.AddSubject(horizonapi.SubjectConfig{
		Kind:      "ServiceAccount",
		Name:      c.blackduck.Spec.Namespace,
		Namespace: c.blackduck.Spec.Namespace,
	})
	clusterRoleBinding.AddRoleRef(horizonapi.RoleRefConfig{
		APIGroup: "",
		Kind:     "ClusterRole",
		Name:     "synopsys-operator-admin",
	})

	clusterRoleBinding.AddLabels(utils.GetVersionLabel("clusterRoleBinding", c.blackduck.Spec.Version))

	return []*components.ClusterRoleBinding{clusterRoleBinding}
}

func NewClusterrolebinding(blackduck *v1.Blackduck) types.ClusterRoleBindingInterface {
	return &clusterrolebinding{blackduck: blackduck}
}