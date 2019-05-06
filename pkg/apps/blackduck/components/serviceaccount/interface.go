package serviceaccount

import "github.com/blackducksoftware/horizon/pkg/components"

type ServiceAccountInterface interface {
	GetServiceAcccount() []*components.ServiceAccount
}
