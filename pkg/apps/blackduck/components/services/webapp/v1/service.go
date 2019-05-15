package v1

import (
	horizonapi "github.com/blackducksoftware/horizon/pkg/api"
	"github.com/blackducksoftware/horizon/pkg/components"
	v1 "github.com/blackducksoftware/synopsys-operator/pkg/api/blackduck/v1"
	"github.com/blackducksoftware/synopsys-operator/pkg/apps/blackduck/components/utils"
	"github.com/blackducksoftware/synopsys-operator/pkg/apps/blackduck/types"
	"github.com/blackducksoftware/synopsys-operator/pkg/util"
)

type service struct {
	blackduck *v1.Blackduck
}

func NewService(blackduck *v1.Blackduck) types.ServiceInterface {
	return &service{blackduck: blackduck}
}

func (s *service) GetService() *components.Service {
	return util.CreateService("webapp", utils.GetLabel("webapp-logstash"), s.blackduck.Spec.Namespace, "8443", "8443", horizonapi.ClusterIPServiceTypeDefault, utils.GetVersionLabel("webapp-logstash", s.blackduck.Spec.Version))
}
