package v1

import (
	horizonapi "github.com/blackducksoftware/horizon/pkg/api"
	"github.com/blackducksoftware/horizon/pkg/components"
	v1 "github.com/blackducksoftware/synopsys-operator/pkg/api/blackduck/v1"
	"github.com/blackducksoftware/synopsys-operator/pkg/apps/blackduck/components/utils"
	"github.com/blackducksoftware/synopsys-operator/pkg/apps/blackduck/types"
)

type serviceaccount struct {
	blackduck *v1.Blackduck
}

func (s *serviceaccount) GetServiceAcccount() []*components.ServiceAccount {
	svc := components.NewServiceAccount(horizonapi.ServiceAccountConfig{
		Name:      s.blackduck.Spec.Namespace,
		Namespace: s.blackduck.Spec.Namespace,
	})

	svc.AddLabels(utils.GetVersionLabel("serviceAccount", s.blackduck.Spec.Namespace))

	return []*components.ServiceAccount{svc}
}

func NewServiceaccount(blackduck *v1.Blackduck) types.ServiceAccountInterface {
	return &serviceaccount{blackduck: blackduck}
}