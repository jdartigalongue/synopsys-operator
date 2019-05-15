package v1

import (
	horizonapi "github.com/blackducksoftware/horizon/pkg/api"
	"github.com/blackducksoftware/horizon/pkg/components"
	v1 "github.com/blackducksoftware/synopsys-operator/pkg/api/blackduck/v1"
	"github.com/blackducksoftware/synopsys-operator/pkg/apps/blackduck/components/utils"
	"github.com/blackducksoftware/synopsys-operator/pkg/apps/blackduck/types"
	"github.com/blackducksoftware/synopsys-operator/pkg/util"
	"strings"
)

type service struct {
	blackduck *v1.Blackduck
}

func NewService(blackduck *v1.Blackduck) types.ServiceInterface {
	return &service{blackduck: blackduck}
}

func (s *service) GetService() *components.Service {
	switch strings.ToUpper(s.blackduck.Spec.ExposeService){
	case "NODEPORT":
		return s.GetWebServerNodePortService()
	case "LOADBALANCER":
		return s.GetWebServerLoadBalancerService()
	default:
		return nil
	}
}


//GetWebServerNodePortService will return the webserver nodeport service
func (s *service) GetWebServerNodePortService() *components.Service {
	return util.CreateService("webserver-exposed", utils.GetLabel("webserver"), s.blackduck.Spec.Namespace, "443", "8443", horizonapi.ClusterIPServiceTypeNodePort, utils.GetLabel("webserver-exposed"))
}

// GetWebServerLoadBalancerService will return the webserver loadbalancer service
func (s *service) GetWebServerLoadBalancerService() *components.Service {
	return util.CreateService("webserver-exposed", utils.GetLabel("webserver"), s.blackduck.Spec.Namespace, "443", "8443", horizonapi.ClusterIPServiceTypeLoadBalancer, utils.GetLabel("webserver-exposed"))
}
