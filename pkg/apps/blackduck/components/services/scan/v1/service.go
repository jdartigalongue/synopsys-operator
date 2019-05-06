package v1

import (
	horizonapi "github.com/blackducksoftware/horizon/pkg/api"
	"github.com/blackducksoftware/horizon/pkg/components"
	v1 "github.com/blackducksoftware/synopsys-operator/pkg/api/blackduck/v1"
	"github.com/blackducksoftware/synopsys-operator/pkg/apps/blackduck/components/services"
	"github.com/blackducksoftware/synopsys-operator/pkg/apps/blackduck/components/utils"
	"github.com/blackducksoftware/synopsys-operator/pkg/util"
)

type service struct {
	blackduck *v1.Blackduck
}

func NewService(blackduck *v1.Blackduck) services.ServiceInterface {
	return &service{blackduck: blackduck}
}

func (s *service) GetService() *components.Service {
	return util.CreateService("scan", utils.GetLabel("scan"), s.blackduck.Spec.Namespace, "8443", "8443", horizonapi.ClusterIPServiceTypeDefault, utils.GetVersionLabel("scan",s.blackduck.Spec.Version))
}
