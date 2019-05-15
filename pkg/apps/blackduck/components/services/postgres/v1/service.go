package v1

import (
	"github.com/blackducksoftware/horizon/pkg/components"
	v1 "github.com/blackducksoftware/synopsys-operator/pkg/api/blackduck/v1"
	"github.com/blackducksoftware/synopsys-operator/pkg/apps/blackduck/components/utils"
	"github.com/blackducksoftware/synopsys-operator/pkg/apps/blackduck/types"
	"github.com/blackducksoftware/synopsys-operator/pkg/apps/database/postgres"
)

type service struct {
	blackduck *v1.Blackduck
}

func NewService(blackduck *v1.Blackduck) types.ServiceInterface {
	return &service{blackduck: blackduck}
}

func (s *service) GetService() *components.Service {
	p := postgres.Postgres{
		Namespace:              s.blackduck.Spec.Namespace,
		Port:                   "5432",
		Labels:                 utils.GetVersionLabel("postgres", s.blackduck.Spec.Version),
	}
	return p.GetPostgresService()
}
