package clusterrolebinding

import "github.com/blackducksoftware/horizon/pkg/components"

type ClusterRoleBindingInterface interface {
	GetClusterRoleBinding() []*components.ClusterRoleBinding
}
